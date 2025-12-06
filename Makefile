echo-run:
	go run app/echo-server/main.go

swaggo-install:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

echo-swagger:
	cd app/echo-server && swag init