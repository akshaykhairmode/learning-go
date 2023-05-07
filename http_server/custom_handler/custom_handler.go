package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var userData = users{data: make(map[int]details), RWMutex: &sync.RWMutex{}}

var strToIntResponse = response{statusCode: 400, Status: "failed", Message: "Invalid ID passed"}

type response struct {
	Status     string   `json:"status"`
	Message    string   `json:"message"`
	Data       *details `json:"data,omitempty"`
	statusCode int
}

type users struct {
	data   map[int]details
	currID int
	*sync.RWMutex
}

type details struct {
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
}

type jsonHandler struct{}

func (mh jsonHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	rw.Header().Add("X-AbilityRush-Domain", "abilityrush.com")
	rw.Header().Add("Content-Type", "application/json")

	var resp response

	base := r.URL.Path

	fmt.Printf("base: %v\n", base)

	switch base {
	case "/get_user":
		id := r.URL.Query().Get("id")
		resp = getUser(id)
	case "/add_user":
		resp = addUser(r.Body)
	case "/delete_user":
		id := r.URL.Query().Get("id")
		resp = delUser(id)
	default:
		resp.statusCode = 404
		resp.Message = "404 Not Found"
		resp.Status = "failed"
	}

	respBytes, _ := json.Marshal(resp)

	rw.WriteHeader(resp.statusCode)
	rw.Write(respBytes)
	fmt.Println("------------------------------------------")
}

func main() {

	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Testing123"))
	})

	fmt.Println("Started")
	http.ListenAndServe(":7000", new(jsonHandler))

}

func delUser(idStr string) response {

	fmt.Println("Got ID : ", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return strToIntResponse
	}

	var resp response

	userData.Lock()
	delete(userData.data, id)
	userData.Unlock()

	resp.statusCode = 200
	resp.Message = "Success"
	resp.Status = "User Deleted"
	return resp

}

func getUser(idStr string) response {

	fmt.Println("Got ID : ", idStr)

	var resp response

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return strToIntResponse
	}

	resp.statusCode = 200
	resp.Status = "success"

	userData.RLock()
	ud, ok := userData.data[id]
	userData.RUnlock()

	if !ok {
		resp.Message = "User does not exist"
		return resp
	}

	resp.Message = "User Found"
	resp.Data = &ud
	return resp
}

func addUser(r io.ReadCloser) response {

	var resp response

	data, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		resp.statusCode = 500
		resp.Status = "failed"
		resp.Message = "Some Error Occured"
		return resp
	}

	ud := details{}

	if err := json.Unmarshal(data, &ud); err != nil {
		fmt.Println(err)
		resp.statusCode = 400
		resp.Status = "failed"
		resp.Message = "Invalid JSON passed"
		return resp
	}

	userData.Lock()
	userData.currID += 1
	userData.data[userData.currID] = ud
	userData.Unlock()

	return response{statusCode: 200, Status: "success", Message: "User Added Successfully, UID : " + strconv.Itoa(userData.currID)}

}
