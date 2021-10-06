package configurer

// Configurer -
type Configurer interface {
	Get(key string) string
	Set(key string, value string)
	Add(key string, required bool) (err error)
}
