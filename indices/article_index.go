package indices

import (
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/search_results"
)

type ArticleIndex interface {
	Populate(a article.Article, c *redis_client.RedisClient) error
	Update(article article.Article, c *redis_client.RedisClient) error
	Search(q Query) ([]string, search_results.Cursor, error)
}

func PopulateIndices(a article.Article, c *redis_client.RedisClient) error {
	i := GetIndices(c)
	for _, index := range i {
		err := index.Populate(a, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateIndices(a article.Article, c *redis_client.RedisClient) error {
	i := GetIndices(c)
	for _, index := range i {
		err := index.Update(a, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetIndices(c *redis_client.RedisClient) []ArticleIndex {
	return []ArticleIndex{TopicIndex{c}, DateIndex{c}, FullTextSearchIndex{c}}
}
