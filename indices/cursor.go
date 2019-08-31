package indices

type Cursor struct {
	Forward int64 `json:"forward"`
	Reverse int64 `json:"reverse"`
	Count int64 `json:"count"`
}

func NewCursor(q Query, results []string) Cursor {
	cur := Cursor{}
	cur.Count = int64(len(results))
	cur.Forward = q.Cursor + cur.Count
	if cur.Count < int64(q.Limit) {
		cur.Forward = -1
	}
	cur.Reverse = q.Cursor - int64(q.Limit)
	if q.Cursor == 0 {
		cur.Reverse = -1
	}
	return cur
}
