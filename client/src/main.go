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
		var reply = Reply{
			Domain: domain,
			JobID:  reqjson.JobID,
			Err:    err,
		}

		resp, err := send(reqjson)
		if err != nil {
			log.Println("send request error: ", err)
			reply.Err = err
			client.Emit("response", reply)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("send request error: ", err)
			reply.Err = err
			client.Emit("response", reply)
			return
		}
		defer resp.Body.Close()

		reply.Body = body
		reply.Header = resp.Header
		reply.StatusCode = resp.StatusCode

		log.Printf("reply %+v\n", reply)
		client.Emit("response", reply)
	})

	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
		log.Fatal("disconnection")
	})

	client.On("join", func(room string) {
		log.Printf("get domain name:%s\n", room)
		domain = room
	})

	client.On("echo", func(echo string) {
		log.Printf("server echo :%s\n", echo)
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		msg := string(data)
		client.Emit("echo", msg)
		log.Printf("send message:%v\n", msg)
	}
}

func send(r request) (resp *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", os.Getenv("PROXY"), r.Path), bytes.NewReader(r.Body))
	if err != nil {
		return nil, err
	}
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return client.Do(req)
}
