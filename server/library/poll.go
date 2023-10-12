package library

import (
	"strings"
	"sync"
)

type Poll struct {
	rw sync.RWMutex

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
	p.rw.Lock()
	p.fd = 0
	p.data = strings.Split(url, ";")
	p.rw.Unlock()
}

func (p *Poll) Get() string {
	if len(p.data) == 0 {
		return ""
	}

	p.rw.Lock()
	url := p.data[p.fd]
	p.fd++
	if p.fd == len(p.data) {
		p.fd = 0
	}
	p.rw.Unlock()

	return url
}
