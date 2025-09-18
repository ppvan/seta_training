package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (me *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", me.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/posts", me.createPostHandler)
	router.HandlerFunc(http.MethodGet, "/v1/posts/search-by-tag", me.searchByTagHandler)
	router.HandlerFunc(http.MethodGet, "/v1/posts/:id", me.getPostHandler)

	return router
}
