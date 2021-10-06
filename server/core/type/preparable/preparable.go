package preparable

// Preparable -
type Preparable interface {
	TableName() string
	GetOneSQL() string
	GetOneForUpdateSQL() string
	GetManySQL() string
	CreateOneSQL() string
	UpdateOneSQL() string
	UpdateManySQL() string
	DeleteOneSQL() string
	DeleteManySQL() string
	RemoveOneSQL() string
	RemoveManySQL() string
}
