package socket

import (
	"log"
	"math/rand"
	"net/url"
	"os"

	"itgserver/src/job"

	socketio "github.com/googollee/go-socket.io"
)

var Server = socketio.NewServer(nil)

var char = [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func SocketOn() {
	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

		domain, err := getDomain(Server, s.URL().RawQuery)
		if err != nil {
			s.Emit("regitfail", err.Error())
			return nil
		}

		// 如果有重複的則回傳錯誤
		for i := 0; i < len(Server.Rooms("/")); i++ {
			if Server.Rooms("/")[i] == domain {
				s.Emit("regitfail", "domain already use")
				return nil
			}
		}

		log.Println("connected:", s.ID())
		s.Join(domain) //將這個連線加入一個房間,房間名為domain
		s.Emit("join", domain)

		return nil
	})

	Server.OnEvent("/", "response", func(s socketio.Conn, reply job.Reply) {
		job.PassReplyByChannel(reply)
	})

	Server.OnEvent("/", "echo", func(s socketio.Conn, msg string) {
		s.Emit("echo", msg)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Printf("closed socketID: %s, by: %s", s.ID(), reason)
	})

	if err := Server.Serve(); err != nil {
		log.Fatalf("socketio listen error: %s\n", err)
	}
	defer Server.Close()
}

func getDomain(Server *socketio.Server, urlInput string) (string, error) {
	params, err := url.ParseQuery(urlInput)
	if err != nil {
		return "", err
	}
	for key, value := range params {
		// 如果有指定domain則給他這個domain
		if key == "domain" {
			return value[0] + "." + os.Getenv("ROOT_DOMAIN"), nil
		}
	}
	// 沒有指定domain則隨機給他一個
	var domain = ""
	for i := 0; i < 6; i++ {
		domain = domain + char[rand.Intn(26)]
	}
	return domain + "." + os.Getenv("ROOT_DOMAIN"), nil
}
