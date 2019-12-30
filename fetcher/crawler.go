package fetcher

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Spec struct {
	Id int64 `json:"id"`
	Url string `json:"url"`
	Interval int `json:"interval"`
}

type Result struct {
	Response  string `json:"response"`
	Duration  float64 `json:"duration"`
	CreatedAt float64 `json:"created_at"`
}

type entry struct {
	spec    Spec
	done    chan struct{}
	results []*Result
	mutex   *sync.RWMutex
}

func newEntry(s Spec) *entry {
	ent := &entry{spec: s}
	ent.done = make(chan struct{})
	ent.results = make([]*Result, 0, 16)
	ent.mutex = &sync.RWMutex{}
	return ent
}

type crawler struct {
	// Http client is safe for concurrent use
	client *http.Client
	// current Id
	id int64
	// Map Spec Id to entry
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

func (c *crawler) crawl(ent *entry) {
	log.Printf("Crawling id: %d, url: %s", ent.spec.Id, ent.spec.Url)
	t0 := time.Now()
	var response string
	resp, err := c.client.Get(ent.spec.Url)
	if err != nil {
		log.Printf("Http request failed id: %d, url: %s, err: %s",
			ent.spec.Id, ent.spec.Url, err.Error())
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed reading body id: %d, url: %s, err: %s",
				ent.spec.Id, ent.spec.Url, err.Error())
		} else {
			response = string(bodyBytes)
		}
	}
	r := &Result{
		Response:  response,
		Duration:  float64(time.Since(t0).Nanoseconds()) / 1e9,
		CreatedAt: float64(t0.UnixNano()) / 1e9,
	}

	ent.mutex.Lock()
	defer ent.mutex.Unlock()

	ent.results = append(ent.results, r)
}

func (c *crawler) task(e *entry) {
	log.Printf("Starting task id: %d, url: %s, interval: %d",
		e.spec.Id, e.spec.Url, e.spec.Interval)
	for {
		c.crawl(e)
		select {
			case <- e.done:
				log.Printf("Stopping task id: %d, url: %s, interval: %d",
					e.spec.Id, e.spec.Url, e.spec.Interval)
				return
			case <- time.After(time.Duration(e.spec.Interval) * time.Second):
		}
	}
}

func (c *crawler) put(s Spec) Spec {
	if s.Id == 0 {
		s.Id = c.nextId()
	}
	ent := newEntry(s)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	entPrev, ok := c.entries[s.Id]
	if ok {
		close(entPrev.done)
	}

	go c.task(ent)
	c.entries[s.Id] = ent

	return s
}

func (c *crawler) del(id int64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entPrev, ok := c.entries[id]
	if ok {
		close(entPrev.done)
		delete(c.entries, id)
	} else {
		return errors.New("entry not found")
	}

	return nil
}

func (c *crawler) getResults(id int64) ([]*Result, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ent, ok := c.entries[id]
	if !ok {
		return nil, errors.New("entry not found")
	}

	ent.mutex.RLock()
	defer ent.mutex.RUnlock()

	results := make([]*Result, len(ent.results))
	copy(results, ent.results)

	return results, nil
}

func (c *crawler) getSpecs() []*Spec {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	specs := make([]*Spec, 0, len(c.entries))
	for _, ent := range c.entries {
		specs = append(specs, &ent.spec)
	}

	return specs
}
