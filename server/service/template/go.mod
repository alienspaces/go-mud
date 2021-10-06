module gitlab.com/alienspaces/go-boilerplate/server/service/template

go 1.15

require (
	github.com/google/uuid v1.1.1
	github.com/jmoiron/sqlx v1.3.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/alienspaces/go-boilerplate/server/core v1.0.0
	gitlab.com/alienspaces/go-boilerplate/server/schema v1.0.0
)

replace (
	gitlab.com/alienspaces/go-boilerplate/server/constant => ../../constant
	gitlab.com/alienspaces/go-boilerplate/server/core => ../../core
	gitlab.com/alienspaces/go-boilerplate/server/schema => ../../schema

	gitlab.com/alienspaces/go-boilerplate/server/service/template => ../../service/template
)
