package config

// Item defines a valid environment variable and whether it is required
type Item struct {
	Key      string
	Required bool
}

func NewItems(keys []string, required bool) []Item {
	items := make([]Item, len(keys))
	for i, k := range keys {
		items[i] = NewItem(k, required)
	}
	return items
}

func NewItem(key string, required bool) Item {
	return Item{Key: key, Required: required}
}
