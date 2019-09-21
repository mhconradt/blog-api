package main

import (
	"encoding/json"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/indices"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
)

func UpdateArticle(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	a := article.Article{}
	er := NewErrorResponder(w)
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		er(err, 403)
		return
	}
	if err = indices.UpdateIndices(a, c); err != nil {
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

/*
TEST CASES:
1. Removing fields:
I should be able to send a PUT with the Id of an existing article without changing its data.
*/
