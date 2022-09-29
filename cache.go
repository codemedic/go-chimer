package chimer

import "errors"

func NewCache[K comparable, V any](f func(K) (V, error)) *Cache[K, V] {
	if f == nil {
		panic(errors.New("bad callback given"))
	}

	return &Cache[K, V]{
		m: make(map[K]V),
		f: f,
	}
}

type Cache[K comparable, V any] struct {
	m map[K]V
	f func(s K) (V, error)
}

func (c *Cache[K, V]) Get(k K) (V, error) {
	var v V
	var found bool
	var err error
	if v, found = c.m[k]; found {
		return v, nil
	}

	v, err = c.f(k)
	if err != nil {
		return v, err
	}

	c.m[k] = v

	return v, nil
}
