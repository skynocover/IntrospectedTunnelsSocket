package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

var domain = ""

type request struct {
	JobID  string
	Path   string
	Method string
	Body   []byte
	Header http.Header
}

type Reply struct {
	JobID      string
	Domain     string
	Body       []byte
	Header     http.Header
	StatusCode int
	Err        error
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	uri := fmt.Sprintf("%s/socket.io/", os.Getenv("DOMAIN"))

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func(err string) {
		log.Printf("error: %s\n", err)
		log.Fatal(err)
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})
	client.On("regitfail", func(err string) {
		log.Fatal(err)
	})

	client.On("request", func(reqjson request) {
		// var reqjson request
		// json.Unmarshal(req, &reqjson)
		// log.Printf("on request: %+v\n", reqjson)

		resp, err := send(reqjson)
		if err != nil {
			var reply = Reply{
				Domain: domain,
				JobID:  reqjson.JobID,
				Err:    err,
			}
			client.Emit("response", reply)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			var reply = Reply{
				Domain: domain,
				JobID:  reqjson.JobID,
				Err:    err,
			}
			client.Emit("response", reply)
		}
		defer resp.Body.Close()

		var reply = Reply{
			Domain:     domain,
			JobID:      reqjson.JobID,
			Body:       body,
			Header:     resp.Header,
			StatusCode: resp.StatusCode,
		}

		log.Printf("reply %+v\n", reply)
		client.Emit("response", reply)
	})

	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
		log.Fatal("disconnection")
	})

	client.On("join", func(room string) {
		log.Printf("room name:%s\n", room)
		domain = room
	})

	// client.Emit("regist")

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("request", command)
		log.Printf("send message:%v\n", command)
	}
}

func send(r request) (resp *http.Response, err error) {
	client := &http.Client{}

	//這邊可以任意變換 http method  GET、POST、PUT、DELETE
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", os.Getenv("PROXY"), r.Path), bytes.NewReader(r.Body))
	if err != nil {
		log.Println(err)
		return
	}
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return client.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()
	// return ioutil.ReadAll(resp.Body)
}
