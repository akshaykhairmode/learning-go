package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var GC = globalClients{data: make(map[ClientID]*Client), RWMutex: &sync.RWMutex{}}

type globalClients struct {
	data map[ClientID]*Client
	*sync.RWMutex
}

type Client struct {
	Conn *websocket.Conn
	Name string
	Id   ClientID
}

type ClientID string

func (gc *globalClients) add(cl *Client) {

	gc.Lock()
	defer gc.Unlock()
	gc.data[cl.Id] = cl
}

func (gc *globalClients) del(cid ClientID) {

	gc.Lock()
	defer gc.Unlock()
	delete(gc.data, cid)
}

func clientMW(next func(cl *Client, rw http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	var cid string
	var conn *websocket.Conn
	var err error

	return func(rw http.ResponseWriter, r *http.Request) {

		qv := r.URL.Query()
		name := qv.Get("name")

		if name == "" {
			name = generateRandomName()
		}

		cid = uuid.New().String()
		h := http.Header{}
		h.Add(idHeader, cid)

		conn, err = upgrader.Upgrade(rw, r, h)
		if err != nil {
			log.Printf("Error while upgrading connection: %v", err)
			return
		}

		cl := new(Client)
		cl.Conn = conn
		cl.Id = ClientID(cid)
		cl.Name = name

		GC.add(cl)
		defer GC.del(cl.Id)

		next(cl, rw, r)
	}

}
