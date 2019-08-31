package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
	"strconv"
)

func GetArticle(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	er := NewErrorResponder(w)
	v := mux.Vars(r)
	id, _ := strconv.Atoi(v["id"])
	a, err := c.GetArticle(id)
	if err != nil {
		er(err, 404)
		return
	}
	NewArticleResponder(w)(a, 200)
	fmt.Println(w.Header())
}
