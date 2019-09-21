package main

import (
	"github.com/mhconradt/blog-api/indices"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
)

func ListArticles(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	er := NewErrorResponder(w)
	q := indices.ParseQuery(r.URL.Query())
	results, err := indices.Search(q, c)
	if err != nil {
		er(err, 500)
		return
	}
	NewSearchResultsResponder(w)(&results, 200)
	return
}

/*
BENCHMARKS:

No Cursor
  5.2ms
Yes Cursor (json)
	6.2ms

*/
