package indices

import (
	"github.com/mhconradt/blog-api/search_results"
)

func NewCursor(q Query, results []string) search_results.Cursor {
	cur := search_results.Cursor{}
	cur.Count = int32(len(results))
	cur.Forward = q.Cursor + cur.Count
	if cur.Count < int32(q.Limit) {
		cur.Forward = -1
	}
	cur.Reverse = q.Cursor - int32(q.Limit)
	if q.Cursor == 0 {
		cur.Reverse = -1
	}
	cur.Term = q.Term
	return cur
}
