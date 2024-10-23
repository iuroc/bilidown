package router

import "net/http"

func API() *http.ServeMux {
	router := http.NewServeMux()
	
	return router
}
