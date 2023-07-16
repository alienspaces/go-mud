package configurer

// Configurer provides an interface to get and set key/value strings. Typically
// implemented to ingest environment variables and make them available to the
// application.
type Configurer interface {
	Get(key string) string
	MustGet(key string) (string, error)
	GetAll() map[string]string
	Set(key string, value string)
	Add(key string, required bool) (err error)
}
