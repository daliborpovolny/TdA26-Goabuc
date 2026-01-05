## To run this server:
- use the start scripts in the root folder


## .env
Enviromental variables are used to control apps, these variables can be set when running the app command,
or thanks to a library we use from a .env file situated in the same folder as main.go
- example file is included in the repo as .example_env, to use it rename it to .env

## Architecture
cmd/tourbackend/main.go is the entry point of the app, this is done by convention
static/ is a folder for all static assets as well as user created assests - such as course materials. These assests will be availabe at: api/static/{filename}
internal/ is a folder with all the components of the app, such as database, auth, courses, etc. main.go imports these components and plugs them into routes
- each component (in praxis not all have been refactored yet) consists of three files
    - errors.go - here all the custom signal errors are declared
    - service.go - the middle math between handler and database, offers clean intereface to interact with the database and returns errors when something goes wrong
    - handler.go - this is a collection of functions that get triggered by a request to a specific endpoint, each endpoint handler creates its own context object, this gives the handler access to the service
    
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

## Testing
ideal would be to run the tests on tourdecloud
but these take time, test.py is a script mimicking some the of the tests, to run simply do: `python test.py`