package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, rErr := http.NewRequest(r.Method, fmt.Sprintf("https:/%s", r.URL.String()), r.Body)
	if rErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(rErr.Error()))
		return
	}
	for k, v := range r.Header {
		for _, i := range v {
			req.Header.Add(k, i)
		}
	}
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(parseFormErr.Error()))
		return
	}
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(reqErr.Error()))
		return
	}
	w.WriteHeader(resp.StatusCode)
	respBody, respErr := ioutil.ReadAll(resp.Body)
	if respErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(respErr.Error()))
		return
	}
	w.Write(respBody)
}
