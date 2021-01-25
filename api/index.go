package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Print(r.URL.String())
	Try(func() {
		client := &http.Client{}
		req, rErr := http.NewRequest(r.Method, fmt.Sprintf("https:/%s", r.URL.String()), r.Body)
		ThrowIfError(rErr)

		for k, v := range r.Header {
			for _, i := range v {
				req.Header.Add(k, i)
			}
		}

		ThrowIfError(req.ParseForm())

		resp, reqErr := client.Do(req)
		ThrowIfError(reqErr)

		respBody, respErr := ioutil.ReadAll(resp.Body)
		ThrowIfError(respErr)

		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)

	}).Catch(func(i interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
		if e, ok := i.(error); ok {
			w.Write([]byte(e.Error()))
		} else {
			w.Write([]byte("Unknown error"))
		}
	})
}

type CatchHandler interface {
	 Catch(handler func(interface{}))
}

type catchHandler struct {
	err interface{}
}

func (c *catchHandler) Catch(handler func(interface{}))  {
	if c.err != nil {
		handler(c.err)
	}
}

func Try(block func()) *catchHandler {
	c := &catchHandler{}
	defer func() {
		defer func() {
			if v := recover(); v != nil {
				c.err = c
			}
		}()
		block()
	}()
	return c
}

func ThrowIfError(e error)  {
	if e != nil {
		panic(e)
	}
}