package preparable

type Repository interface {
	TableName() string
	Attributes() []string
	GetOneSQL() string
	GetOneForUpdateSQL() string
	GetManySQL() string
	CreateOneSQL() string
	UpdateOneSQL() string
	DeleteOneSQL() string
	RemoveOneSQL() string
}

type Query interface {
	SQL() string
}
