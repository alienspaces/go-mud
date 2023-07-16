package preparable

type Repository interface {
	TableName() string
	Attributes() []string
	GetOneSQL() string
	GetManySQL() string
	CreateOneSQL() string
	UpdateOneSQL() string
	UpdateManySQL() string
	DeleteOneSQL() string
	DeleteManySQL() string
	RemoveOneSQL() string
	RemoveManySQL() string
}

type Query interface {
	Name() string
	SQL() string
}
