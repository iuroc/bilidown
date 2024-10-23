package router

import "net/http"

func API() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/checkLogin", func(w http.ResponseWriter, r *http.Request) {

	})

	return router
}
