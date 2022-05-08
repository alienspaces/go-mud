module gitlab.com/alienspaces/go-mud/server/service/template

go 1.15

require (
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

	gitlab.com/alienspaces/go-mud/server/service/template => ../../service/template
)
