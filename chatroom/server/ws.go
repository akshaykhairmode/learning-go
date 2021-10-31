package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type middlewareData struct {
	c    *websocket.Conn
	vars map[string]string
	id   ClientID
}

type WSMiddleware func(middlewareData)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBuffSize,
	WriteBufferSize: writeBuffSize,
}

func handleChatroomWS(md middlewareData) {

	defer md.c.Close()

	rooms.RLock()
	room, roomExists := rooms.Data[ChatroomID(md.vars["id"])]
	rooms.RUnlock()

	if !roomExists {
		md.c.WriteMessage(websocket.CloseTryAgainLater, []byte("Room does not exist"))
		return
	}

	cl := new(Client)
	cl.Conn = md.c
	room.addClient(cl)

	for {
		mt, message, err := md.c.ReadMessage()
		if err != nil {
			log.Println("Reading Message Error :", err)
			break
		}
		log.Printf("recv: %s", message)
		room.pushToBroadcast(message, mt, md.id)
	}
}

func WS(next WSMiddleware) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {

		id := uuid.New().String()
		h := http.Header{}
		h.Add(idHeader, id)

		c, err := upgrader.Upgrade(rw, r, h)
		if err != nil {
			log.Printf("Error while upgrading connection: %v", err)
			return
		}

		vars := mux.Vars(r)

		md := middlewareData{c: c, vars: vars, id: ClientID(id)}

		next(md)
	}
}
