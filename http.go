package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func newRequest() (req *Request) {
	return &Request{
		Packet: Packet{
			FromID:    "10003001",
			ServerKey: "omt5W465fjwlrtxcEco97kew2dkdrorqqq",
			Data:      Data{},
		},
	}
}

func executeRequest(url string, request *Request) (data []byte) {
	reqBody, _ := json.Marshal(request)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log(err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode > 200 {
		logf("StatusCode is %v", resp.StatusCode)
		return nil
	}

	respData, _ := ioutil.ReadAll(resp.Body)

	return respData
}
