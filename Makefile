swag:
	swag init
build:
	GOOS=linux GOARCH=amd64 go build -o go-cms main.go