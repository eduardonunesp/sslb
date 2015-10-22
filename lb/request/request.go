package request

import (
	"net/http"
)

type SSLBRequest struct {
	Status   int
	Header   http.Header
	Body     []byte
	Internal bool
}

type SSLBRequestChan chan SSLBRequest

func NewWorkerRequestErr(status int, body []byte) SSLBRequest {
	return SSLBRequest{Status: status, Body: body, Internal: true}
}

func NewWorkerRequest(status int, header http.Header, body []byte) SSLBRequest {
	return SSLBRequest{Status: status, Header: header, Body: body}
}
