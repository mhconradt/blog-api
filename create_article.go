package main

import (
	"encoding/json"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/indices"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
	"time"
)

func CreateArticle(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	a := article.Article{}
	er := NewErrorResponder(w)
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		er(err, 403)
		return
	}
	a.ID = c.NewId()
	a.Timestamp = time.Now().Unix()
	if err = indices.PopulateIndices(a, c); err != nil {
		er(err, 500)
		return
	}
	if err = c.WriteArticle(a); err != nil {
		er(err, 500)
		return
	}
	NewResponder(w)("success!", 200)
	return
}
