package job

import (
	"log"
	"net/http"
	"sync"
)

// 儲存channel,用來做每次request時溝通的管道
var jobID = &sync.Map{}

type Reply struct {
	JobID      string
	Domain     string
	Body       []byte
	Header     http.Header
	StatusCode int
	Err        error
}

// api會產生一個新的jobID 產生一個channel給他
func GetChannel(jid string) chan Reply {
	tunnel := make(chan Reply)
	jobID.Store(jid, tunnel)
	return tunnel
}

// 處理response 根據jobID找到跟api溝通的channel
func PassReplyByChannel(reply Reply) {
	fchan, ok := jobID.Load(reply.JobID)
	if !ok {
		log.Printf("Can't not find chan for domain: %s\n", reply.Domain)
	} else {
		tempChan := fchan.(chan Reply)
		tempChan <- reply //通過channel將socket回傳的東西傳給api
		jobID.Delete(reply.JobID)
	}
}
