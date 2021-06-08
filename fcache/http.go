package fcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/__fcache__/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s]: %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.URL.Path, p.basePath) {
		p.Log("format error request url should contains basePath: %s", p.basePath)
	}
	// <basePath>/<groupName>/<path>
	p.Log(r.URL.String())
	// TODO make a summary /a/tom  and  a/tom
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	p.Log("parts: %v", parts)
	if len(parts) != 2 {
		//p.Log("request format should be <basePath>/<groupName>/<path>")
		http.Error(w, "base request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	p.Log("groupName: %s, key: %s", groupName, key)

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group ", http.StatusNotFound)
		return
	}

	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/octet-stream")
	_, _ = w.Write(value.ByteSlice())
}
