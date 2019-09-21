package redis_client

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/search_results"
	"github.com/mhconradt/blog-api/util"
	"os"
)

var addr string

type RedisClient struct {
	*redis.Client
}

func ArticleKeyFromId(id int32) string {
	return fmt.Sprintf(config.ArticlePrefix+"%v", id)
}

func (c *RedisClient) WriteArticle(a article.Article) error {
	if err := c.SetArticle(a); err != nil {
		return err
	}
	if err := c.WriteCache(a); err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) WriteCache(a article.Article) error {
	// proto.Marshal
	snip := search_results.SnippetFromArticle(a)
	b, err := proto.Marshal(&snip)
	if err != nil {
		return err
	}
	key := fmt.Sprintf(config.SnippetPrefix+"%v", a.Id)
	if err := c.Set(key, string(b), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) NewId() int32 {
	return int32(c.Incr(config.IDKey).Val())
}

func (c *RedisClient) SetArticle(a article.Article) error {
	return c.HMSet(ArticleKeyFromId(a.Id), a.ToRedis()).Err()
}

func (c *RedisClient) GetArticle(id int32) (*article.Article, error) {
	keys := []string{ArticleKeyFromId(id), "body", "views", "title", "topics", "timestamp", "id"}
	va, err := c.EvalSha(config.GetArticle, keys).Result()
	if err != nil {
		return new(article.Article), err
	}
	vm := util.ZipMap(keys[1:], va.([]interface{}))
	return article.FromRedis(vm)
}

func GetConfig() *redis.Options {
	return &redis.Options{
		Addr: addr,
	}
}

func GetRedisClient() *RedisClient {
	cfg := GetConfig()
	c := redis.NewClient(cfg)
	return &RedisClient{c}
}

func init() {
	if a, found := os.LookupEnv("REDIS_ADDRESS"); !found {
		addr = "localhost:6379"
	} else {
		addr = a
	}
}
