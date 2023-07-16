package counter

import (
	"strconv"
)

type Counter map[string]int

func New() Counter {
	return Counter{}
}

func (m Counter) Increment(key string) {
	m[key]++
}

func (m Counter) Decrement(key string) {
	m[key]--
}

func (m Counter) CountToString(key string) string {
	c := m[key]
	return strconv.Itoa(c)
}
