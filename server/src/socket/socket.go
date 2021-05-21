package socket

import (
	"log"

	"itgserver/src/domain"
	"itgserver/src/job"

	socketio "github.com/googollee/go-socket.io"
)

var Server *socketio.Server

func SocketOn() {
	Server = socketio.NewServer(nil)

	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		// 發一個domain給新的連線
		domain := domain.GetDomain()
		if domain == "" {
			s.Emit("regitfail", "No free Domain")
		} else {
			s.Join(domain)
			s.Emit("join", domain)
		}
		return nil
	})

	Server.OnEvent("/", "response", func(s socketio.Conn, reply job.Reply) {
		job.PassReplyByChannel(reply)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// 斷線的時候釋放domain
		domain.LetDomainFree(s.ID())
		log.Println("closed", reason)
	})

	if err := Server.Serve(); err != nil {
		log.Fatalf("socketio listen error: %s\n", err)
	}
	defer Server.Close()
}
