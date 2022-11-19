package repositor

// Repositor -
type Repositor interface {
	Init() error
	TableName() string
}
