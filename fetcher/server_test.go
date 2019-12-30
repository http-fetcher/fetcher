package fetcher

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

type crawlerMock struct{
	id int64
}

func (c *crawlerMock) Put(s Spec) Spec {
	return Spec{Id: c.id}
}

func (c *crawlerMock) Del(id int64) error {
	return nil
}

func (c *crawlerMock) GetResults(id int64) ([]*Result, error) {
	return nil, nil
}

func (c *crawlerMock) GetSpecs() []*Spec {
	return nil
}

func TestHandlePut(t *testing.T) {
	maxBodySize := int64(1024 * 1024)

	crawlerMock := &crawlerMock{id: 123}
	srv := NewServer(maxBodySize, chi.NewRouter(), crawlerMock)

	s := Spec{Url: "https://httpbin.org/range/15", Interval: 60}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(s)
	if err != nil {
		t.Error("failed encoding input json")
	}

	r := httptest.NewRequest("POST", "/api/fetcher", &body)
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("http status code is not OK")
	}

	output := &Spec{}
	err = json.NewDecoder(w.Body).Decode(&output)
	if err != nil {
		t.Error("failed decoding output json")
	}

	if output.Id != crawlerMock.id {
		t.Error("returned invalid id")
	}
}