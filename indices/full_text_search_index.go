package indices

import (
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/go-redis/redis"
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/redis_client"
	"github.com/mhconradt/blog-api/util"
	"strings"
)

type FullTextSearchIndex struct {
	*redis_client.RedisClient
}

func (fts FullTextSearchIndex) reverseIndexKey(a article.Article) string {
	return fmt.Sprintf(config.FullTextSearchReversePrefix+"%v", a.ID)
}

func (fts FullTextSearchIndex) forwardIndexKey(word string) string {
	return fmt.Sprintf(config.FullTextSearchForwardPrefix+"%v", word)
}

func (fts FullTextSearchIndex) getWords(body string) []string {
	cleaned := stopwords.CleanString(body, "en", false)
	words := strings.Split(cleaned, " ")
	notEmpty := make([]string, 0, len(words))
	for _, word := range words {
		if len(word) > 0 {
			notEmpty = append(notEmpty, word)
		}
	}
	return notEmpty
}

func (fts FullTextSearchIndex) countWords(words []string) map[string]int {
	ctr := make(map[string]int)
	// IF the word is not present, add it with count one. Otherwise, increment count.
	for _, word := range words {
		if val, ok := ctr[word]; !ok {
			ctr[word] = 1
		} else {
			ctr[word] = val + 1
		}
	}
	return ctr
}

func (fts FullTextSearchIndex) getWordCounts(body string) map[string]int {
	words := fts.getWords(body)
	return fts.countWords(words)
}

func (fts FullTextSearchIndex) getMembers(wc map[string]int) []redis.Z {
	num := len(wc)
	results := make([]redis.Z, 0, num)
	for word, count := range wc {
		results = append(results, redis.Z{
			Score:  float64(count),
			Member: word,
		})
	}
	return results
}

func UnzipMap(m map[string]int) ([]string, []int) {
	n := len(m)
	words, counts := make([]string, n), make([]int, n)
	i := 0
	for k, v := range m {
		words[i], counts[i] = k, v
		delete(m, k)
		i++
	}
	return words, counts
}

func Keys(m map[string]int) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		delete(m, k)
		i++
	}
	return keys
}

func (fts FullTextSearchIndex) GetExistingCounts(a article.Article, c *redis_client.RedisClient) (m map[string]int, err error) {
	opt := redis.ZRangeBy{
		Max: config.PositiveInfinity,
		Min: config.NegativeInfinity,
	}
	k := fts.reverseIndexKey(a)
	result, err := c.ZRangeByScoreWithScores(k, opt).Result()
	if err != nil {
		return m, err
	}
	m = make(map[string]int)
	for _, z := range result {
		m[z.Member.(string)] = int(z.Score)
	}
	return m, nil
}

func MapCountsToStrings(counts []int) []string {
	s := make([]string, len(counts))
	for i, c := range counts {
		s[i] = fmt.Sprintf("%v", c)
	}
	return s
}

func (fts FullTextSearchIndex) populateIndex(a article.Article, wc map[string]int, c *redis_client.RedisClient) error {
	// If no writes, don't do anything
	if len(wc) < 1 {
		return nil
	}
	mems := fts.getMembers(wc)
	k := fts.reverseIndexKey(a)
	err := c.ZAdd(k, mems...).Err()
	if err != nil {
		return err
	}
	words, counts := UnzipMap(wc)
	args := []string{fmt.Sprintf("%v", a.ID), config.SnippetPrefix, config.FullTextSearchForwardPrefix}
	args = append(args, MapCountsToStrings(counts)...)
	err = c.EvalSha(config.PopulateFullTextForwardIndex, words, args).Err()
	if err != nil {
		fmt.Println("hypothesis confirmed/")
		return err
	}
	return nil
}

func (fts FullTextSearchIndex) removeIndexEntries(a article.Article, keys []string, c *redis_client.RedisClient) error {
	return c.EvalSha(config.RemoveHitFromFTS, keys, a.ID, config.FullTextSearchForwardPrefix, config.FullTextSearchReversePrefix, config.SnippetPrefix).Err()
}

func (fts FullTextSearchIndex) Populate(a article.Article, c *redis_client.RedisClient) error {
	// Go (ha!) through the article body and add count the occurrence of all character combinations between 2 and 5 characters.
	wc := fts.getWordCounts(a.Body)
	return fts.populateIndex(a, wc, c)
}

func MergeMaps(a, b map[string]int) map[string]int {
	// write all keys from a onto b
	for k, v := range a {
		b[k] = v
		// Garbage collection, bitches.
		delete(a, k)
	}
	return b
}

func (fts FullTextSearchIndex) Update(a article.Article, c *redis_client.RedisClient) error {
	// Get all by score
	// can be NEW | UPDATED | DELETED | SAME
	// if same, do nothing
	// if new, add to reverse and forward indices
	// if updated, add to reverse and forward indices
	// if deleted, remove from forward and reverse indices
	exc, err := fts.GetExistingCounts(a, c)
	wc := fts.getWordCounts(a.Body)
	if err != nil {
		return err
	}
	additions := make(map[string]int)
	updates := make(map[string]int)
	for k, v := range wc {
		current, exists := exc[k]
		if !exists {
			additions[k] = v
			continue
		}
		if current != v {
			updates[k] = v
		}
		// What's left in exc will be deletes
		delete(exc, k)
	}
	deletes := Keys(exc)
	writes := MergeMaps(updates, additions)
	err = fts.populateIndex(a, writes, c)
	if err != nil {
		return err
	}
	return fts.removeIndexEntries(a, deletes, c)

	// group into write and delete
	// different scores? update
	// in exc but not wc? delete
	// in wc but not exc? insert
}

func (fts FullTextSearchIndex) Search(q Query) ([]string, Cursor, error) {
	result, err := fts.EvalSha(config.SearchZIndex, []string{fts.forwardIndexKey(q.Term)}, q.Cursor, q.Limit).Result()
	if err != nil {
		if strings.Index(err.Error(), "table expected") == -1 {
			return []string{}, Cursor{}, err
		}
		result = []interface{}{}
	}
	vs := util.ToStringSlice(result.([]interface{}))
	cur := NewCursor(q, vs)
	return vs, cur, nil
}

/*
How many windows? Length of the item - length of the substring plus one.
Intuition: can fit substr once at beginning, and then slide it over len(item) - len(substr) times
Full word only at first. Sounds good?
*/

// want to maintain digits
func init() {
	stopwords.DontStripDigits()
}
