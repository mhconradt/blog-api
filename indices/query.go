package indices

import (
	"fmt"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
	"net/url"
	"strconv"
	"strings"
)

type Direction int

const (
	DefaultCursor = iota
)

const (
	Ascending  Direction = 1
	Descending Direction = -1
)

type Index int

const (
	Date Index = iota
	FullTextSearch
	Topic
)

type Query struct {
	PageDirection Direction
	Index
	Cursor int64
	Term   string
	Limit  int
}

func DirectionFromQuery(v url.Values, f string) Direction {
	dir, found := StringAtField(v, f)
	if !found {
		return Ascending
	}
	switch strings.ToLower(dir) {
	case "desc":
	case "descending":
		return Descending
	}
	return Ascending
}

func IndexFromQuery(v url.Values) Index {
	index, found := StringAtField(v, "index")
	if !found {
		return Date
	}
	switch strings.ToLower(index) {
	case "text":
		return FullTextSearch
	case "topic":
		return Topic
	default:
		return Date
	}
}

func CursorFromQuery(v url.Values) int64 {
	curStr, found := StringAtField(v, "cursor")
	if !found {
		return DefaultCursor
	}
	cur, err := strconv.ParseInt(curStr, 10, 64)
	if err != nil {
		return DefaultCursor
	}
	return cur
}

func StringAtField(v url.Values, f string) (string, bool) {
	fv := v[f]
	if fv == nil || len(fv) == 0 {
		return "", false
	}
	return fv[0], true
}

func LimitFromQuery(v url.Values) int {
	ls, found := StringAtField(v, "limit")
	if !found {
		return config.DefaultLimit
	}
	l, _ := strconv.Atoi(ls)
	return l
}

func ParseQuery(v url.Values) Query {
	pd := DirectionFromQuery(v, "pageDir")
	fmt.Println("pd: ", pd)
	i := IndexFromQuery(v)
	fmt.Println("i:", i)
	t, _ := StringAtField(v, "term")
	c := CursorFromQuery(v)
	return Query{
		PageDirection: pd,
		Index:         i,
		Term:          t,
		Cursor:        c,
	}
}

func GetIndexForQuery(q Query, c *redis_client.RedisClient) ArticleIndex {
	switch q.Index {
	case Topic:
	default:
		return TopicIndex{c}
	}
	return TopicIndex{c}
}
