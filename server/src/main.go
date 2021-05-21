package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"itgserver/src/socket"

	"github.com/joho/godotenv"

	"github.com/kataras/iris/v12"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	socket.SetDomain()

	go socket.SocketOn()

	app := iris.New()

	log.Println(fmt.Sprintf("Serving at localhost:%s...", os.Getenv("SERVER_LISTEN")))

	app.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(socket.Server))
	app.HandleMany("", "/{any:path}", handle)
	app.HandleDir("/", "../asset")

	if err := app.Run(
		iris.Addr(":"+os.Getenv("SERVER_LISTEN")),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	); err != nil {
		log.Fatal("failed run app: ", err)
	}
}

type request struct {
	Path   string
	Method string
	Body   []byte
	Header http.Header
}

func handle(ctx iris.Context) {
	w := ctx.Request()

	host := "localhost"
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		ctx.StatusCode(404)
		return
	}
	log.Println("body: ", string(body))

	var req = request{
		Path:   fmt.Sprintf("%s?%s", w.URL.Path, w.URL.RawQuery),
		Method: w.Method,
		Body:   body,
		Header: w.Header,
	}

	log.Printf("req: %+v\n", req)
	reqbyte, _ := json.Marshal(req)
	socket.Server.BroadcastToRoom("/", host, "request", reqbyte)

	fchan, ok := socket.CurrentDomain.Load(host)
	if !ok {
		log.Println("can't not fine chan from doamin: ", host)
		ctx.StatusCode(409)
		return
	}
	tempChan := fchan.(chan socket.Reply)
	reply := <-tempChan
	var temp struct {
		Body string
	}
	temp.Body = string(reply.Body)
	ctx.ContentType(reply.ContentType)
	for key, item := range reply.Header {
		for _, value := range item {
			ctx.Header(key, value)
		}

	}
	ctx.JSON(temp)
}
