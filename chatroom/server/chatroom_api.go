package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/gorilla/mux"
)

type CRResponse struct {
	Status  string     `json:"status"`
	CRID    ChatroomID `json:"chatroom_id,omitempty"`
	Message string     `json:"message,omitempty"`
}

type CRMiddleware func(CRMiddlewareData, http.ResponseWriter)

type CRMiddlewareData struct {
	vars map[string]string
}

func CR(next CRMiddleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		crmd := CRMiddlewareData{}
		crmd.vars = vars

		next(crmd, rw)
	}
}

func CreateChatroom(crmd CRMiddlewareData, rw http.ResponseWriter) {

	crName := crmd.vars["name"]
	if crName == "" {
		crName = generateRandomName()
	}

	//Create the chatroom.
	crid := rooms.create(crName)

	log.Printf("Created Chatroom : %v", crName)

	r := CRResponse{}
	r.CRID = crid
	r.Status = "Success"

	out, err := json.Marshal(r)
	if err != nil {
		log.Printf("Error Occured while marhshal : %v", err)
		rw.Write([]byte(`{"status":"failed"}`))
		return
	}

	log.Printf("Response is : %s", out)

	rw.WriteHeader(200)
	rw.Write(out)
}

func ListChatroom(crmd CRMiddlewareData, rw http.ResponseWriter) {

	output := []map[string]string{}

	rooms.RLock()
	for _, cr := range rooms.Data {
		chatroomInfo := map[string]string{}
		chatroomInfo["id"] = string(cr.ID)
		chatroomInfo["name"] = cr.Name
		output = append(output, chatroomInfo)
	}
	rooms.RUnlock()

	data, _ := json.Marshal(output)

	rw.WriteHeader(200)
	rw.Write(data)
}

func generateRandomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return nameGenerator.Generate()
}
