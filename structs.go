package main

import "encoding/json"

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"-"`
	IDs     string          `json:"-"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage `json:"error,omitempty"`
}

type Request struct {
	JSONRPC string
	Method  string
	ID      int
	IDs     string
	Params  []json.RawMessage
}
