module gitlab.com/alienspaces/go-mud/server/service/game

go 1.18

require (
	github.com/brianvoe/gofakeit v3.18.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/alienspaces/go-mud/server/core v1.0.0
	gitlab.com/alienspaces/go-mud/server/schema v1.0.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace (
	gitlab.com/alienspaces/go-mud/server/constant => ../../constant
	gitlab.com/alienspaces/go-mud/server/core => ../../core
	gitlab.com/alienspaces/go-mud/server/schema => ../../schema

	gitlab.com/alienspaces/go-mud/server/service/game => ../../service/game
)
