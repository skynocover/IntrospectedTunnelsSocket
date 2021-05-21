package domain

import (
	"os"
	"strings"
	"sync"
)

// 儲存所有domain被socketid註冊的狀態
var domainID = &sync.Map{}

// 紀錄所以domain
func SetDomain() {
	var domains = strings.Split(os.Getenv("DOMAINS"), ",")
	for _, domain := range domains {
		domainID.Store(domain, "")
	}
}

// 取得一個未註冊的domain
func GetDomain() string {
	var domain = ""
	domainID.Range(func(key interface{}, value interface{}) bool {
		domain = key.(string)
		id := value.(string)
		if id != "" { //找下一個
			return true
		}
		domainID.Store(domain, id)
		return false
	})
	return domain
}

// 斷線之後釋放domain
func LetDomainFree(sid string) {
	domainID.Range(func(key interface{}, value interface{}) bool {
		domain := key.(string)
		id := value.(string)
		if id != sid { //找下一個
			return true
		}
		domainID.Store(domain, "")
		return false
	})
}
