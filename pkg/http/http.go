package http

import (
	"net/url"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var ins *fasthttp.Client
var once sync.Once

func client() *fasthttp.Client {
	once.Do(func() {
		ins = &fasthttp.Client{
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
		}
	})
	return ins
}

func Get(url string) (body []byte, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	req.SetRequestURI(url)
	err = client().Do(req, res)
	if err != nil {
		return
	}
	if res != nil {
		body = res.Body()
	}
	return
}

func Post(url string, values url.Values) (body []byte, err error) {
	req, res := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	args := req.PostArgs()
	for k := range values {
		args.Add(k, values.Get(k))
	}
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	err = client().Do(req, res)
	if err != nil {
		return
	}
	if res != nil {
		body = res.Body()
	}
	return
}
