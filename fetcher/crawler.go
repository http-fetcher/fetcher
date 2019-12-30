package fetcher

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type spec struct {
	Id int64
	Url string
	Interval int
}

type result struct {
	response string
	duration time.Duration
	createdAt time.Time
}

type entry struct {
	spec spec
	done chan struct{}
	results []*result
}

func newEntry(s spec) *entry {
	e := &entry{spec: s}
	e.done = make(chan struct{})
	e.results = make([]*result, 16)
	return e
}

type crawler struct {
	// Http client is safe for concurrent use
	client *http.Client
	// current Id
	id int64
	// Map spec Id to entry
	entries map[int64]*entry
	mutex *sync.RWMutex
}

func NewCrawler(client *http.Client) *crawler {
	c := &crawler{client: client}
	c.entries = make(map[int64]*entry)
	c.mutex = &sync.RWMutex{}
	return c
}

func (c *crawler) nextId() int64 {
	return atomic.AddInt64(&c.id, 1)
}

func (c *crawler) crawl(e *entry) {
	log.Printf("Crawling Id: %d, Url: %s", e.spec.Id, e.spec.Url)
}

func (c *crawler) task(e *entry) {
	log.Printf("Starting task Id: %d, Url: %s", e.spec.Id, e.spec.Url)
	for {
		c.crawl(e)
		select {
			case <- e.done:
				log.Printf("Finising task Id: %d, Url: %s", e.spec.Id, e.spec.Url)
				return
			case <- time.After(time.Duration(e.spec.Interval) * time.Second):
		}
	}
}

func (c *crawler) put(s spec) spec {
	if s.Id == 0 {
		s.Id = c.nextId()
	}
	ent := newEntry(s)

	c.mutex.Lock()

	entPrev, ok := c.entries[s.Id]
	if ok {
		close(entPrev.done)
	}

	go c.task(ent)
	c.entries[s.Id] = ent

	c.mutex.Unlock()
	return s
}

func (c *crawler) del(id int64) {
	c.mutex.Lock()
	entPrev, ok := c.entries[id]
	if ok {
		close(entPrev.done)
		delete(c.entries, id)
	}
	c.mutex.Unlock()
}
