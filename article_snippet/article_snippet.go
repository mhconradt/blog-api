package article_snippet

import (
	"github.com/mhconradt/blog-api/article"
	"github.com/mhconradt/blog-api/config"
)

type ArticleSnippet struct {
	ID        int    `json:"id"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
}

func SnippetFromArticle(a article.Article) ArticleSnippet {
	bodySnippet := func(body string) string {
		if len(body) > config.SnippetLength {
			return body[:config.SnippetLength]
		} else {
			return body
		}
	}(a.Body)
	return ArticleSnippet{
		a.ID,
		a.Timestamp,
		a.Title,
		bodySnippet,
	}
}
