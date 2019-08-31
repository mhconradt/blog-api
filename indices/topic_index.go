package indices

import (
	"fmt"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/util"
)

type TopicIndex struct {
	*redis_client.RedisClient
}

func (t TopicIndex) Populate(a article.Article, c *redis_client.RedisClient) error {
	cmd := c.EvalSha(config.PopulateIndex, a.Topics, "topics", a.WithHitPrefix())
	return cmd.Err()
}

func (t TopicIndex) Update(a article.Article, c *redis_client.RedisClient) error {
	fmt.Println("topics index upd8: ", a.Topics)
	if len(a.Topics) == 0 {
		return nil
	}
	// get current topics
	// only need to add and delete necessary ones
	reverseIndexKey := fmt.Sprintf("topics.reverse.%v", a.ID)
	old := c.LRange(reverseIndexKey, 0, -1).Val()
	additions, removals := Diff(old, a.Topics)
	fmt.Println("a: ", additions)
	fmt.Println("r: ", removals)
	if len(additions) > 0 {
		if cmd := c.EvalSha(config.PopulateIndex, additions, "topics", a.WithHitPrefix()); cmd.Err() != nil {
			return cmd.Err()
		}
	}
	if len(removals) > 0 {
		if cmd := c.EvalSha(config.RemoveIndexEntries, removals, "topics", a.WithHitPrefix()); cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func (t TopicIndex) Search(q Query) ([]string, Cursor, error) {
	end := q.Cursor + (q.Cursor + int64(int(q.PageDirection)*q.Limit))
	min, max := func(a, b int64) (int64, int64) {
		if a > b {
			return b, a
		}
		return a, b
	}(q.Cursor, end)
	result, err := t.EvalSha(config.SearchListIndex, []string{}, "topics", q.Term, min, max).Result()
	if err != nil {
		return []string{}, Cursor{}, err
	}
	vs := util.ToStringSlice(result.([]interface{}))
	cur := NewCursor(q, vs)
	return vs, cur, nil
}

func Incrementer(m map[string]uint8) func(by uint8) func(field string) {
	return func(i uint8) func(field string) {
		return func(field string) {
			if _, ok := m[field]; !ok {
				m[field] = i
			} else {
				m[field] += i
			}
			return
		}
	}
}

func Diff(old, current []string) ([]string, []string) {
	m := make(map[string]uint8)
	incrementer := Incrementer(m)
	byOne := incrementer(uint8(1))
	for _, val := range old {
		byOne(val)
	}
	byTwo := incrementer(uint8(2))
	for _, val := range current {
		byTwo(val)
	}
	removals, additions := make([]string, 0, len(old)), make([]string, 0, len(current))
	for k, v := range m {
		if v == uint8(1) {
			removals = append(removals, k)
		}
		if v == uint8(2) {
			additions = append(additions, k)
		}
	}
	fmt.Println("m: ", m)
	return additions, removals
}
