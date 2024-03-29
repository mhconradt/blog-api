package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/util"
	"log"
	"net/http"
)

type RedisHTTPHandler func(w http.ResponseWriter, r *http.Request, client *redis_client.RedisClient)

type VanillaHTTPHandler func(w http.ResponseWriter, r *http.Request)

func WrapRedisHandler(h RedisHTTPHandler) VanillaHTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		c := redis_client.GetRedisClient()
		h(w, r, c)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/articles", WrapRedisHandler(CreateArticle)).Methods("POST")
	r.HandleFunc("/articles", WrapRedisHandler(UpdateArticle)).Methods("PUT")
	r.HandleFunc("/articles/{id}", WrapRedisHandler(GetArticle)).Methods("GET")
	r.HandleFunc("/articles", WrapRedisHandler(ListArticles)).Methods("GET")
	port := util.LookupWithDefault("PORT", "3000")
	fmt.Println("listening for requests on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

/*
ENDPOINTS:
- POST /articles (create an article)
- PUT /articles (update an article)
- GET /articles (find articles based on search term, topic, date, etc.)
- DELETE /articles/{id} (remove an article from the site)
*/

// How to decouple rendering from data access?
// Have the rendering process call the API. Boom.

func NewArticleResponder(w http.ResponseWriter) func(a article.Article, status int) {
	return func(a article.Article, status int) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(a)
	}
}

func NewSearchResultsResponder(w http.ResponseWriter) func(r SearchResult, status int) {
	return func(r SearchResult, status int) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(r)
	}
}

func NewResponder(w http.ResponseWriter) func(msg string, status int) {
	return func(msg string, status int) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(msg))
	}
}

func NewErrorResponder(w http.ResponseWriter) func(err error, status int) {
	return func(err error, status int) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(err.Error()))
	}
}
