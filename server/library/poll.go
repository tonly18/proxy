package library

import (
	"strings"
	"sync"
)

type Poll struct {
	mu sync.RWMutex

	fd   int
	data []string
}

var poll *Poll
var once sync.Once

func NewPoll() *Poll {
	once.Do(func() {
		poll = &Poll{}
	})

	return poll
}

func (p *Poll) Set(url string) {
	p.mu.Lock()
	p.fd = 0
	p.data = strings.Split(url, ";")
	p.mu.Unlock()
}

// 轮询host
func (p *Poll) Get() string {
	if len(p.data) == 0 {
		return ""
	}

	p.mu.RLock()
	url := p.data[p.fd]
	p.fd++
	if p.fd == len(p.data) {
		p.fd = 0
	}
	p.mu.RUnlock()

	return url
}

// 按ID分配host
func (p *Poll) GetByID(id int) string {
	if len(p.data) == 0 {
		return ""
	}

	p.mu.RLock()
	url := p.data[id%len(p.data)]
	p.mu.RUnlock()

	return url
}
