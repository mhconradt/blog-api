package main

import (
	"fmt"
	"github.com/mhconradt/blog-api/indices"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
	"strings"
)

func ListArticles(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	er := NewErrorResponder(w)
	q := indices.ParseQuery(r.URL.Query())
	i := indices.GetIndexForQuery(q, c)
	results, _, err := i.Search(q)
	if err != nil {
		er(err, 500)
		return
	}
	json := "[" + strings.Join(results, ",") + "]"
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if _, err = w.Write([]byte(json)); err != nil {
		fmt.Println("error writing response: ", err)
	}
}

/*
BENCHMARKS:

No Cursor
  5.2ms
Yes Cursor

 */
