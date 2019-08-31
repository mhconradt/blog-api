package indices

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/util"
	"strings"
)

type DateIndex struct {
	*redis_client.RedisClient
}

func (d DateIndex) Populate(a article.Article, c *redis_client.RedisClient) error {
	z := redis.Z{
		Score:  float64(a.Timestamp),
		Member: a.WithHitPrefix(),
	}
	return d.ZAdd(config.DateIndexKey, z).Err()
}

func (d DateIndex) Update(_ article.Article, _ *redis_client.RedisClient) error {
	return nil
}

func (d DateIndex) Search(q Query) ([]string, Cursor, error) {
	result, err := d.EvalSha(config.SearchZIndex, []string{config.DateIndexKey}, q.Cursor, q.Limit).Result()
	if err != nil {
		fmt.Println(err)
		if strings.Index(err.Error(), "table expected") == -1 {
			return []string{}, Cursor{}, err
		}
		result = []interface{}{}
	}
	vs := util.ToStringSlice(result.([]interface{}))
	cur := NewCursor(q, vs)
	return vs, cur, nil
}

// SORTDIR can only be ASC atm
// each index must implement scanForward and scanBackward

/*
can scan a sorted set with an offset, this could be a valid cursor.
score based iteration is better than offset based iteration
*/
