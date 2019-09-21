package indices

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/search_results"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	Cursor int32
	Term   string
	Limit  int32
}

func DirectionFromQuery(v url.Values, f string) Direction {
	dir, found := StringAtField(v, f)
	if !found {
		return Ascending
	}
	switch strings.ToLower(dir) {
	case "desc":
		fallthrough
	case "descending":
		return Descending
	case "asc":
		fallthrough
	case "ascending":
		fallthrough
	default:
		return Ascending
	}
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
	case "date":
		fallthrough
	case "time":
		fallthrough
	default:
		return Date
	}
}

func CursorFromQuery(v url.Values) int {
	curStr, found := StringAtField(v, "cursor")
	if !found {
		return DefaultCursor
	}
	cur, err := strconv.ParseInt(curStr, 10, 64)
	if err != nil {
		return DefaultCursor
	}
	return int(cur)
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
	i := IndexFromQuery(v)
	t, _ := StringAtField(v, "term")
	c := CursorFromQuery(v)
	l := LimitFromQuery(v)
	return Query{
		PageDirection: pd,
		Index:         i,
		Term:          t,
		Cursor:        int32(c),
		Limit:         int32(l),
	}
}

func GetIndexForQuery(q Query, c *redis_client.RedisClient) ArticleIndex {
	switch q.Index {
	case Topic:
		return TopicIndex{c}
	case FullTextSearch:
		return FullTextSearchIndex{c}
	case Date:
		fallthrough
	default:
		return DateIndex{c}
	}
}

func MarshalProtos(bufs []string) ([]*search_results.ArticleSnippet, error) {
	count := len(bufs)
	snippets := make([]*search_results.ArticleSnippet, count, count)
	for i, b := range bufs {
		snippets[i] = &search_results.ArticleSnippet{}
		if err := proto.Unmarshal([]byte(b), snippets[i]); err != nil {
			return snippets, err
		}
	}
	fmt.Println(snippets)
	return snippets, nil
}

func Search(q Query, c *redis_client.RedisClient) (search_results.SearchResults, error) {
	i := GetIndexForQuery(q, c)
	results, cur, err := i.Search(q)
	if err != nil {
		return search_results.SearchResults{}, err
	}
	start := time.Now().UnixNano()
	snippets, err := MarshalProtos(results)
	end := time.Now().UnixNano()
	fmt.Println("decode duration: ", end-start)
	if err != nil {
		return search_results.SearchResults{}, err
	}
	return search_results.SearchResults{
		Results: snippets,
		Cursor:  &cur,
	}, nil
}
