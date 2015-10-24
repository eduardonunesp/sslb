package request

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

type SSLBRequest struct {
	Status   int
	Header   http.Header
	Body     []byte
	Internal bool

	Upgraded bool
	Address  string
}

type SSLBRequestChan chan SSLBRequest

func NewWorkerRequestErr(status int, body []byte) SSLBRequest {
	return SSLBRequest{Status: status, Body: body, Internal: true}
}

func NewWorkerRequest(status int, header http.Header, body []byte) SSLBRequest {
	return SSLBRequest{Status: status, Header: header, Body: body}
}

func NewWorkerRequestUpgraded() SSLBRequest {
	return SSLBRequest{Upgraded: true}
}

func copy(dest *bufio.ReadWriter, src *bufio.ReadWriter) {
	buf := make([]byte, 40*1024)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 {
			return
		}
		dest.Write(buf[0:n])
		dest.Flush()
	}
}

func copyBidir(frontendConn io.ReadWriteCloser, rwFront *bufio.ReadWriter,
	backendConn io.ReadWriteCloser, rwBack *bufio.ReadWriter) {

	finished := make(chan bool)

	go func() {
		copy(rwBack, rwFront)
		backendConn.Close()
		finished <- true
	}()

	go func() {
		copy(rwFront, rwBack)
		frontendConn.Close()
		finished <- true
	}()

	<-finished
	<-finished
}

func (s *SSLBRequest) HijackWebSocket(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)

	if !ok {
		log.Println("Error: Webserver doesn't support hijacking")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	frontendConn, buffer, err := hj.Hijack()
	defer frontendConn.Close()

	URL := &url.URL{}
	UrlParsed, _ := URL.Parse(s.Address)

	backendConn, err := net.Dial("tcp", UrlParsed.Host)
	if err != nil {
		log.Println("Error: Couldn't connect to backend server")
		http.Error(w, "Internal Error", http.StatusServiceUnavailable)
		return
	}
	defer backendConn.Close()

	err = r.Write(backendConn)
	if err != nil {
		log.Printf("Writing WebSocket request to backend server failed: %v", err)
		return
	}

	copyBidir(frontendConn, buffer, backendConn,
		bufio.NewReadWriter(bufio.NewReader(backendConn), bufio.NewWriter(backendConn)))
}
