package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func httpPost(url string, data []byte) (string, error) {
	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logrus.Error("Problem to create request", err.Error())
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error("Problem making request: ", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("Problem to read response", err.Error())
		return "", err
	}
	return string(respBody), nil
}

func proxyHandler(rw http.ResponseWriter, r *http.Request) {
	// Request payload
	var request_payload Request

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message.
	err := json.NewDecoder(r.Body).Decode(&request_payload)

	if err != nil {
		logrus.Error("wrong method, only POST available.")
		cbody := json.RawMessage(`{"code":-32000,"message":"not able to decode body"}`)
		var cresp = Response{
			JSONRPC: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(rw, string(cresb))

		return
	}
	if r.Method != "POST" {
		logrus.Error("wrong method, only POST available.")
		cbody := json.RawMessage(`{"code":-32000,"message":"wrong method, only POST available"}`)
		var cresp = Response{
			JSONRPC: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(rw, string(cresb))
		return
	}

	// ####### bunch of things we make here ################
	data, _ := json.Marshal(request_payload)

	// #### Passing to the client the req ######
	respBody, err := httpPost(quiknode_proxy.downstream_client_addr, data)
	if err != nil {
		logrus.Error("Problem to connect to the local client: ", err.Error())
		cbody := json.RawMessage(err.Error())
		var cresp = Response{
			JSONRPC: "2.0",
			Error:   cbody,
		}
		cresb, _ := json.Marshal(cresp)
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(rw, string(cresb))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(rw, string(respBody))
}
