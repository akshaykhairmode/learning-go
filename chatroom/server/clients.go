package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

var clients = AllClients{Clients: map[ClientID]*Client{}, RWMutex: &sync.RWMutex{}}

type AllClients struct {
	Clients map[ClientID]*Client
	*sync.RWMutex
}

type Client struct {
	Conn *websocket.Conn
	Name string
	Id   ClientID
}

type ClientID string

func (ac *AllClients) add(cl *Client) {
	ac.Lock()
	defer ac.Unlock()
	ac.Clients[cl.Id] = cl
}
