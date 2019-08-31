package redis_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/article_snippet"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/util"
)

type RedisClient struct {
	*redis.Client
}

func ArticleKeyFromId(id int) string {
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
	buf := bytes.NewBuffer([]byte{})
	snip := article_snippet.SnippetFromArticle(a)
	if err := json.NewEncoder(buf).Encode(snip); err != nil {
		return nil
	}
	jsonStr := string(buf.Bytes())
	key := fmt.Sprintf(config.SnippetPrefix+"%v", a.ID)
	if err := c.Set(key, jsonStr, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisClient) NewId() int {
	return int(c.Incr(config.IDKey).Val())
}

func (c *RedisClient) SetArticle(a article.Article) error {
	return c.HMSet(ArticleKeyFromId(a.ID), a.ToRedis()).Err()
}

func (c *RedisClient) GetArticle(id int) (article.Article, error) {
	keys := []string{"body", "views", "title", "topics", "timestamp", "id"}
	va, err := c.HMGet(ArticleKeyFromId(id), keys...).Result()
	if err != nil {
		return article.Article{}, err
	}
	vm := util.ZipMap(keys, va)
	return article.FromRedis(vm)
}

func GetConfig() *redis.Options {
	return &redis.Options{
		Addr: "localhost:6379",
	}
}

func GetRedisClient() *RedisClient {
	cfg := GetConfig()
	c := redis.NewClient(cfg)
	return &RedisClient{c}
}
