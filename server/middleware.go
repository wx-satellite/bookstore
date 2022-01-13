package server

import (
	"log"
	"mime"
	"net/http"
)

// Logger 打印请求日志
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("recv a %s request from %s", req.Method, req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}

// Validating post/put/patch/delete 请求体必须是 application/json
func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !(req.Method == http.MethodPost ||
			req.Method == http.MethodPut ||
			req.Method == http.MethodDelete ||
			req.Method == http.MethodPatch) {
			return
		}
		// 以上方法请求头必须是 application/json
		mediaType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if mediaType != "application/json" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, req)
	})
}
