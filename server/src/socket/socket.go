package socket

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

var Server *socketio.Server

var freeDomain = &sync.Map{}
var CurrentDomain = &sync.Map{}
var IDDomain = &sync.Map{}

// var Chanmap = &sync.Map{} //全域宣告
// var Chanmap = make(map[string]chan Reply)

type Reply struct {
	Domain      string
	Body        []byte
	ContentType string
	Header      http.Header
}

func SetDomain() {
	var domains = strings.Split(os.Getenv("DOMAINS"), ",")
	for _, domain := range domains {
		freeDomain.Store(domain, "")
	}
}

func SocketOn() {
	Server = socketio.NewServer(nil)

	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())

		var domain = ""
		freeDomain.Range(func(key interface{}, value interface{}) bool { //遍歷需要使用func
			CurrentDomain.Store(key.(string), make(chan Reply))
			IDDomain.Store(s.ID(), key.(string))
			domain = key.(string)
			freeDomain.Delete(key.(string))
			return false //回傳true會繼續下一輪
		})
		if domain == "" {
			s.Emit("regitfail", "No free Domain")
			// return fmt.Errorf("%s regist fail", s.RemoteAddr())
		} else {
			s.Join(domain)
			s.Emit("join", domain)
		}
		return nil
	})

	Server.OnEvent("/", "response", func(s socketio.Conn, reply Reply) {
		fchan, ok := CurrentDomain.Load(reply.Domain)
		if !ok {
			log.Println("can't not fine chan from doamin: ", reply.Domain)
		} else {
			tempChan := fchan.(chan Reply)
			tempChan <- reply
		}
	})

	Server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	Server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	Server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	Server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	if err := Server.Serve(); err != nil {
		log.Fatalf("socketio listen error: %s\n", err)
	}
	defer Server.Close()
}
