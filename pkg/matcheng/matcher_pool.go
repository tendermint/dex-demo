package matcheng

import "sync"

var pool *sync.Pool

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return NewMatcher()
		},
	}
}

func GetMatcher() *Matcher {
	return pool.Get().(*Matcher)
}

func ReturnMatcher(m *Matcher) {
	m.Reset()
	pool.Put(m)
}
