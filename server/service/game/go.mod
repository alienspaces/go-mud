module gitlab.com/alienspaces/go-mud/server/service/game

go 1.15

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/alienspaces/go-mud/server/core v1.0.0
	gitlab.com/alienspaces/go-mud/server/schema v1.0.0
)

replace (
	gitlab.com/alienspaces/go-mud/server/constant => ../../constant
	gitlab.com/alienspaces/go-mud/server/core => ../../core
	gitlab.com/alienspaces/go-mud/server/schema => ../../schema

	gitlab.com/alienspaces/go-mud/server/service/game => ../../service/game
)
