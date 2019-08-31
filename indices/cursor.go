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
	cur.Reverse = func(a, b int64) int64 {
		if a > b {
			return a
		}
		return b
	}(q.Cursor, q.Cursor - int64(q.Limit))
	return cur
}
