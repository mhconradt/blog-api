package article

import (
	"fmt"
	"github.com/mhconradt/blog-api/config"
	"github.com/mhconradt/blog-api/util"
	"strings"
)

func (m Article) WithHitPrefix() string {
	return fmt.Sprintf(config.HitPrefix+"%v", m.Id)
}

func (m Article) ToRedis() map[string]interface{} {
	ma := map[string]interface{}{
		"id": m.Id,
	}
	if m.Timestamp != 0 {
		ma["timestamp"] = m.Timestamp
	}
	if len(m.Topics) != 0 {
		ma["topics"] = strings.Join(m.Topics, config.TopicSeparator)
	}
	if len(m.Title) != 0 {
		ma["title"] = m.Title
	}
	if len(m.Body) != 0 {
		ma["body"] = m.Body
	}
	return ma
}

func FromRedis(am map[string]interface{}) (*Article, error) {
	id := am["id"]
	if id == nil {
		return &Article{}, fmt.Errorf("article not found")
	}
	// Check all of these because this server will NEVER panic. Ever.
	a := Article{Id: util.ToInt(am["id"])}
	if timestamp, ok := am["timestamp"]; ok {
		a.Timestamp = util.ToInt(timestamp)
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
	return &a, nil
}
