package middleware

import (
	"compress/gzip"
	"net/http"
)

func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Encoding") {
		case "gzip":
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer reader.Close()
			r.Body = reader
			next.ServeHTTP(w, r)
		default:
			next.ServeHTTP(w, r)
		}
	})
}
