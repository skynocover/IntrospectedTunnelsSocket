package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

type request struct {
	Path   string
	Method string
	Body   []byte
}

type Reply struct {
	Domain string
	Body   []byte
}

func main() {

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://localhost:8080/socket.io/"

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
	client.On("reply", func(msg string) {
		log.Printf("on message:%v\n", msg)
	})
	client.On("regitfail", func(err string) {
		log.Fatal(err)
	})

	client.On("request", func(req []byte) {
		var reqjson request
		json.Unmarshal(req, &reqjson)
		log.Printf("on request: %+v\n", reqjson)

		var reply = Reply{
			Domain: "localhost",
			Body:   []byte{},
		}
		reply.Body = []byte("response from client")
		client.Emit("response", reply)
	})

	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
		log.Fatal("disconnection")
	})

	client.On("join", func(room string) {
		log.Printf("room name:%s\n", room)
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
