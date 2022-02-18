package middleware

import (
	"log"
	"mime"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
)

func Logging(next http.Handler) http.Handler {
	function := func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("recv a %s request from %s", request.Method, request.RemoteAddr)
		next.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(function)
}

func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get(ContentType)
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if mediaType != ApplicationJson {
			http.Error(writer, "invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
