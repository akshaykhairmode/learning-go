package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var rooms = Rooms{Data: map[ChatroomID]*Chatroom{}, RWMutex: new(sync.RWMutex), Wg: new(sync.WaitGroup)}

type Broadcast struct {
	Message     []byte
	MessageType int
	Cid         ClientID
}

type Chatroom struct {
	Clients          map[ClientID]*Client
	Name             string
	BroadcastChannel chan Broadcast
	*sync.RWMutex
}

type ChatroomID string

type Rooms struct {
	Data map[ChatroomID]*Chatroom
	Wg   *sync.WaitGroup
	*sync.RWMutex
}

type CRMiddleware func(CRMiddlewareData, http.ResponseWriter)

type CRMiddlewareData struct {
	vars map[string]string
}

type CRResponse struct {
	Status string     `json:"status"`
	CRID   ChatroomID `json:"chatroom_id,omitempty"`
}

//create creates a chatroom and allocates memory.
func (r *Rooms) create(name string) ChatroomID {

	crID := ChatroomID(uuid.New().String())

	r.Lock()
	defer r.Unlock()

	cr := new(Chatroom)
	cr.Clients = map[ClientID]*Client{}
	cr.Name = name
	cr.RWMutex = new(sync.RWMutex)
	cr.BroadcastChannel = make(chan Broadcast, 50)

	r.Data[crID] = cr

	rooms.Wg.Add(1)
	go cr.broadcaster(rooms.Wg)
	return crID
}

//delete a chatroom
func (r *Rooms) delete(crID ChatroomID) {

	r.Lock()
	defer r.Unlock()
	delete(r.Data, crID)

}

//addClient to chatroom
func (c *Chatroom) addClient(cl *Client) {

	c.Lock()
	defer c.Unlock()
	c.Clients[cl.Id] = cl

}

//delClient from chatroom
func (c *Chatroom) delClient(cid ClientID) {

	c.Lock()
	defer c.Unlock()
	delete(c.Clients, cid)

}

//delClient from chatroom
func (c *Chatroom) close() {

	c.Lock()
	defer c.Unlock()
	for _, cl := range c.Clients {
		cl.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseMessage, "Closing Chatroom"))
	}

	close(c.BroadcastChannel)

	for ci := range c.Clients {
		delete(c.Clients, ci)
	}

}

//pushToBroadcast will push the message to the broadcast channel
func (c *Chatroom) pushToBroadcast(msg []byte, mt int, cid ClientID) {
	c.BroadcastChannel <- Broadcast{Message: msg, MessageType: mt, Cid: cid}
}

//broadcaster will get the data which needs to be broadcasted and broadcast it.
func (c *Chatroom) broadcaster(wg *sync.WaitGroup) {

	defer wg.Done()

	for b := range c.BroadcastChannel {

		for id, cl := range c.Clients {
			if id == b.Cid { //If its the sender do not send
				continue
			}

			cl.Conn.WriteMessage(b.MessageType, b.Message)
		}
	}
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

func DeleteChatroom(crmd CRMiddlewareData, rw http.ResponseWriter) {

	crID := crmd.vars["id"]
	r := CRResponse{}

	room, exists := rooms.Data[ChatroomID(crID)]

	if !exists {
		r.Status = "Chatroom Does not exist"
		data, _ := json.Marshal(r)
		rw.Write(data)
	}

	room.close()

	r.Status = "Success"
	data, _ := json.Marshal(r)
	rw.Write(data)
}

func generateRandomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return nameGenerator.Generate()
}
