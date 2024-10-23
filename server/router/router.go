package router

import "net/http"

func API() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("6666"))
	})
	return router
}
