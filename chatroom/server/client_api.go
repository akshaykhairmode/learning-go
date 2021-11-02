package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ClientList struct {
	Name string   `json:"client_name"`
	ID   ClientID `json:"client_id"`
}

func (gc *globalClients) list() []ClientList {

	gc.RLock()
	defer gc.RUnlock()

	out := []ClientList{}

	for _, cl := range gc.data {
		out = append(out, ClientList{Name: cl.Name, ID: cl.Id})
	}

	return out
}

func ListAllClients(rw http.ResponseWriter, r *http.Request) {

	data := GC.list()

	enc, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error occured while listing clients : %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Some error occured"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(enc)
}
