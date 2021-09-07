package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Request struct {
	*http.Client
	Debug bool
}

type ReqHeader struct {
	Key   string
	Value string
}

type QueryParam struct {
	Key   string
	Value string
}

func (r *Request) Dial() *Request {
	if r.Client == nil {
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		r.Client = &http.Client{Transport: tr}
	}

	return r
}

func populateWithHeaders(req *http.Request, headers []ReqHeader) {
	if len(headers) == 0 {
		return
	}

	for _, h := range headers {
		if len(h.Key) == 0 || len(h.Value) == 0 {
			continue
		}
		req.Header.Add(h.Key, h.Value)
	}
}

func (r *Request) Get(path string, params []QueryParam, headers ...ReqHeader) ([]byte, int, error) {
	var err error

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, 0, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for _, p := range params {
			q.Add(p.Key, p.Value)
		}
		req.URL.RawQuery = q.Encode()
	}

	if r.Debug {
		b, err2 := httputil.DumpRequestOut(req, true)
		if err2 != nil {
			return nil, 0, err2
		}
		log.Printf("%s\n", string(b))
	}

	populateWithHeaders(req, headers)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	reader := resp.Body
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, err
}