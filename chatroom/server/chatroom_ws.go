package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var rooms = Rooms{Data: map[ChatroomID]*Chatroom{}, RWMutex: new(sync.RWMutex), Wg: new(sync.WaitGroup)}

type Broadcast struct {
	Message     []byte   //The actual message
	MessageType int      //as defined by the gorilla websocket package.
	Cid         ClientID //the UUID of the client.
	Cname       string   //name of the client
}

type Chatroom struct {
	Clients          map[ClientID]*Client
	Name             string
	ID               ChatroomID
	BroadcastChannel chan Broadcast
	*sync.RWMutex
}

type ChatroomID string

type PartnerID string

type Rooms struct {
	Data map[ChatroomID]*Chatroom
	Wg   *sync.WaitGroup
	*sync.RWMutex
}

type ChatData struct {
	Type    string   `json:"type"` //control or text
	Name    string   `json:"sender_name,omitempty"`
	ID      ClientID `json:"sender_id,omitempty"`
	Message string   `json:"message,omitempty"`
}

func (cd ChatData) Marshal() []byte {
	data, _ := json.Marshal(cd)
	return data
}

func (cd ChatData) String() string {
	data, _ := json.Marshal(cd)
	return string(data)
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
	cr.ID = crID

	r.Data[crID] = cr

	rooms.Wg.Add(1)
	go cr.broadcaster(rooms.Wg)
	log.Printf("Chatroom Created - URL : ws://localhost:%v/chatroom/connect", port)
	return crID
}

//delete a chatroom
func (r *Rooms) delete(crID ChatroomID) {

	r.Lock()
	defer r.Unlock()
	delete(r.Data, crID)

}

//addClient to chatroom
func (c *Chatroom) addClient(cl *Client) *Chatroom {

	c.Lock()
	defer c.Unlock()
	c.Clients[cl.Id] = cl
	return c

}

//delClient from chatroom
func (c *Chatroom) delClient(cid ClientID) *Chatroom {

	c.Lock()
	defer c.Unlock()
	delete(c.Clients, cid)
	return c

}

//close checks if any client is there else deletes the chatroom
func (c *Chatroom) close() *Chatroom {

	if len(c.Clients) > 0 { //Dont close if clients still present.
		return c
	}

	rooms.delete(c.ID) //Delete room from the global map.

	c.Lock()
	defer c.Unlock()
	close(c.BroadcastChannel) //Close the broadcast channel.
	log.Printf("Closed Chatroom ID : %v , Name : %v", c.ID, c.Name)
	return c

}

//pushToBroadcast will push the message to the broadcast channel
func (c *Chatroom) pushToBroadcast(msg []byte, mt int, cid ClientID, cname string) {
	c.BroadcastChannel <- Broadcast{Message: msg, MessageType: mt, Cid: cid, Cname: cname}
}

//broadcaster will get the data which needs to be broadcasted and broadcast it.
func (c *Chatroom) broadcaster(wg *sync.WaitGroup) {

	defer wg.Done()

	for b := range c.BroadcastChannel {

		for id, cl := range c.Clients {
			if id == b.Cid { //If its the sender then skip
				continue
			}

			chat := ChatData{ //Create chat data and write to connections
				Name:    b.Cname,
				ID:      b.Cid,
				Message: string(b.Message),
				Type:    "text",
			}

			if err := cl.Conn.WriteJSON(chat); err != nil {
				log.Printf("Error occured while sending message : %v", err)
				continue
			}
		}
	}
}

func chatroomExists(id ChatroomID) bool {

	_, ok := rooms.Data[id]

	return ok

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBuffSize,
	WriteBufferSize: writeBuffSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func chatroomConnect(cl *Client, crid string) {

	rooms.RLock()                        //Read lock the rooms
	room := rooms.Data[ChatroomID(crid)] //Get the room
	rooms.RUnlock()
	defer room.close() //Close deletes the chatroom if no client is connected.

	room.addClient(cl)          //Add client to the chatroom
	defer room.delClient(cl.Id) //Delete client from the chatroom when this function ends.

	for {
		mt, message, err := cl.Conn.ReadMessage()
		if err != nil {
			log.Println("Reading Message Error :", err)
			break
		}
		room.pushToBroadcast(message, mt, cl.Id, cl.Name) //Publish the message to everyone.
	}

}

func chatroomWSHandler(cl *Client, rw http.ResponseWriter, r *http.Request) {

	var crid string

	defer cl.Conn.Close() //Close the conn after conn ends.

	qv := r.URL.Query()
	crid = qv.Get("chatroom_id") //Get the chatroom id

	if crid == "" { //If chatroom id is blank then write a control message
		log.Printf("Got Chatroom ID : %v", crid)
		data := ChatData{Type: "control", Message: "name or chatroom_id empty"}
		cl.Conn.WriteJSON(data)
		return
	}

	if !chatroomExists(ChatroomID(crid)) { //If chatroom does not exist then write a control message
		log.Printf("Chatroom ID : %v does not exist", crid)
		data := ChatData{Type: "control", Message: "given chatroom does not exist"}
		cl.Conn.WriteJSON(data)
		return
	}

	chatroomConnect(cl, crid) //Connect to the chatroom
}
