package main

import (
	"encoding/json"
	"github.com/mhconradt/blog-api/article_snippet"
	"github.com/mhconradt/blog-api/indices"
	"github.com/mhconradt/blog-api/redis_client"
	"net/http"
	"strings"
)

type SearchResult struct {
	Results []article_snippet.ArticleSnippet `json:"results"`
	Cursor  indices.Cursor                   `json:"cursor"`
}

func ListArticles(w http.ResponseWriter, r *http.Request, c *redis_client.RedisClient) {
	er := NewErrorResponder(w)
	q := indices.ParseQuery(r.URL.Query())
	i := indices.GetIndexForQuery(q, c)
	results, cur, err := i.Search(q)
	if err != nil {
		er(err, 500)
		return
	}
	j := "[" + strings.Join(results, ",") + "]"
	jb := []byte(j)
	snippets := make([]article_snippet.ArticleSnippet, len(results))
	err = json.Unmarshal(jb, &snippets)
	if err != nil {
		er(err, 500)
		return
	}
	sr := SearchResult{snippets, cur}
	NewSearchResultsResponder(w)(sr, 200)
	return
}

/*
BENCHMARKS:

No Cursor
  5.2ms
Yes Cursor (json)
	6.2ms

*/
