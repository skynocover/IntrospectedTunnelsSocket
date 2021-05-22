package main

import (
	"fmt"

	"log"

	"os"

	"itgserver/src/domain"
	"itgserver/src/handler"
	"itgserver/src/socket"

	"github.com/joho/godotenv"

	"github.com/kataras/iris/v12"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain.SetDomain()

	go socket.SocketOn()

	app := iris.New()
	log.Println(fmt.Sprintf("Serving at localhost:%s...", os.Getenv("SERVER_LISTEN")))

	app.HandleMany("GET POST", "/socket.io/{any:path}", iris.FromStd(socket.Server))
	app.HandleMany("GET POST PATCH PUT DELETE HEAD OPTIONS CONNECT", "/{any:path}", handler.Handler)

	if err := app.Run(
		iris.Addr(":"+os.Getenv("SERVER_LISTEN")),
		iris.WithoutPathCorrection,
		iris.WithoutServerError(iris.ErrServerClosed),
	); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
