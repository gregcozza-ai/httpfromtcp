package response 

import (
	"fmt"
	"io"
	
	"github.com/gregcozza-ai/httpfromtcp/internal/headers"
)



func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers 
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err 
		}
	}
	_, err :=w.Write([]byte("\r\n"))
	return err 
}