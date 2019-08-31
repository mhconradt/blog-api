package indices

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
)

type DateIndex struct {
	*redis_client.RedisClient
}

func (d DateIndex) Populate(a article.Article) error {
	z := redis.Z{
		Score:  float64(a.Timestamp),
		Member: a.WithHitPrefix(),
	}
	return d.ZAdd(config.DateIndexKey, z).Err()
}

func (d DateIndex) Update() error {
	return nil
}

func (d DateIndex) Search(q Query) ([]string, Cursor, error) {
	opt := redis.ZRangeBy{
		Offset: q.Cursor,
		Count:  int64(q.Limit),
	}
	result, err := d.ZRevRangeByScoreWithScores(config.DateIndexKey, opt).Result()
	if err != nil {
		return []string{}, Cursor{}, err
	}
	fmt.Println(result)
	return []string{}, Cursor{}, nil
}

// SORTDIR can only be ASC atm
// each index must implement scanForward and scanBackward

/*
can scan a sorted set with an offset, this could be a valid cursor.
score based iteration is better than offset based iteration
*/
