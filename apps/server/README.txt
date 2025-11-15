## To run this server:
- go run .
runs on port 1323

## Echo:
Echo is the web server
Echode docs: https://echo.labstack.com/docs/quick-start

## SQL and sqlc:
sqlc is a tool that generates go code to interact with our database
sqcl generates this go code from sql schema and queries (/database/schema.sql, /database/queries.sql)
This generated code lives in /database/gen and must not be edited - any edits will be lost when sqlc is rerun
- sqlc can be installed with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
- to run simply do: sqlc generate

sqlc docs: https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html

file: sqlc.yaml
- a config file for sqlc

## Openapi
file: swagger.yaml
- openapi spec of api we're supposed to implement
- https://editor.swagger.io - use this to visualize and interact with the spec

## Other files
file: go.mod
file: go.sum
these are go project files that shoudn't be manually edited
packages installed (with go get) get added to them

## TODO
add hot reload