package response 

import (
	"fmt"
	
	"github.com/gregcozza-ai/httpfromtcp/internal/headers"
)



func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers 
}

