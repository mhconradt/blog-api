package article

import (
	"fmt"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/util"
	"strings"
)

type Article struct {
	ID        int      `json:"id"`
	Timestamp int64    `json:"timestamp,omitempty"`
	Title     string   `json:"title,omitempty"`
	Body      string   `json:"body,omitempty"`
	Topics    []string `json:"topics,omitempty"`
	Views     int      `json:"views,omitempty"`
}

func (a Article) WithHitPrefix() string {
	return fmt.Sprintf(config.HitPrefix+"%v", a.ID)
}

func (a Article) ToRedis() map[string]interface{} {
	m := map[string]interface{}{
		"id": a.ID,
	}
	if a.Timestamp != 0 {
		m["timestamp"] = a.Timestamp
	}
	if len(a.Topics) != 0 {
		m["topics"] = strings.Join(a.Topics, config.TopicSeparator)
	}
	if len(a.Title) != 0 {
		m["title"] = a.Title
	}
	if len(a.Body) != 0 {
		m["body"] = a.Body
	}
	return m
}

func FromRedis(am map[string]interface{}) (Article, error) {
	id := am["id"]
	if id == nil {
		return Article{}, fmt.Errorf("article not found")
	}
	// Check all of these because this server will NEVER panic. Ever.
	a := Article{ID: util.ToInt(am["id"])}
	if timestamp, ok := am["timestamp"]; ok {
		a.Timestamp = util.ToInt64(timestamp)
	}
	if views, ok := am["views"]; ok {
		a.Views = util.ToInt(views)
	}
	if topics, ok := am["topics"]; ok {
		a.Topics = util.ToArray(topics)
	}
	if body, ok := am["body"]; ok {
		a.Body = body.(string)
	}
	if title, ok := am["title"]; ok {
		a.Title = title.(string)
	}
	return a, nil
}
