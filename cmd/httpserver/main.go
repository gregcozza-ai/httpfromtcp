package main

import (
	"log"
	"os"
	"io"
	"os/signal"
	"syscall"
	"net/http"
	"strings"
	"fmt"

	"github.com/gregcozza-ai/httpfromtcp/internal/request"
	"github.com/gregcozza-ai/httpfromtcp/internal/response"
	"github.com/gregcozza-ai/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		handlerProxy(w, req)
		return 
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusCodeSuccess)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	
}

func handlerProxy(w *response.Writer, req *request.Request) {
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	if path == req.RequestLine.RequestTarget {
		w.WriteStatusLine(response.StatusCode(404))
		return 
	}
	url := "https://httpbin.org/" + path 
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return 
	}
	defer resp.Body.Close()

	//resp.Header.Del("Content-Length")
	//resp.Header.Set("Transfer-Encoding", "chunked")

	w.WriteStatusLine(response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(0)
	headers.Override("Transfer-Encoding", "chunked")
	headers.Remove("Content-Length")
	w.WriteHeaders(headers)

	const maxChunkSize = 1024
	buf := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF{
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
}