# cookit

CookIt is an app for creating recipes, searching for already created ones and managing recipe comments.
The app consists of a server and a cli client.

Server can be started by running
`go run ./cmd/cookit.go`
and can be accessed on localhost:8080

Client can be started by running
`go run ./client/cmd/client.cookit.go`

Mysql instance is required to run the server and configurations for it can be found in the configs/app.env
