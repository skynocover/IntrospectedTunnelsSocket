package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"fmt"
	"itgserver/src/job"
	"itgserver/src/socket"
)

type request struct {
	JobID  string
	Path   string
	Method string
	Body   []byte
	Header http.Header
}

func Handler(ctx iris.Context) {
	w := ctx.Request()

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		ctx.StatusCode(404)
		return
	}

	// 生產這次request的jobid,並取得一個channel
	jid := uuid.New().String()
	tunnel := job.GetChannel(jid)

	// 對domain發送request
	socket.Server.BroadcastToRoom("/", ctx.Host(), "request", request{
		JobID:  jid,
		Path:   fmt.Sprintf("%s?%s", w.URL.Path, w.URL.RawQuery),
		Method: w.Method,
		Body:   body,
		Header: w.Header,
	})

	// 等待channel回覆
	reply := <-tunnel
	if reply.Err != nil {
		var resp struct {
			err error
		}
		resp.err = reply.Err
		ctx.StatusCode(500)
		ctx.JSON(resp)
		return
	}
	// response
	for key, item := range reply.Header {
		for _, value := range item {
			ctx.Header(key, value)
		}
	}
	ctx.StatusCode(reply.StatusCode)
	ctx.Write(reply.Body)
}
