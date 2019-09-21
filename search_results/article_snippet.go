package search_results

import (
	"github.com/mhconradt/blog-api/article"
	"strings"
)

func GetPreview(body string) string {
	paragraphs := strings.Split(body, "\n")
	return paragraphs[0]
}

func SnippetFromArticle(a article.Article) ArticleSnippet {
	bodySnippet := GetPreview(a.Body)
	return ArticleSnippet{
		Id:        a.Id,
		Timestamp: a.Timestamp,
		Title:     a.Title,
		Body:      bodySnippet,
	}
}
