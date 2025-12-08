# run app
echo-run:
	go run app/echo-server/main.go

# api doc
swaggo-install:
	go install github.com/swaggo/swag/cmd/swag@v1.16.4

echo-swagger:
	cd app/echo-server && swag init

# mock
mock-install:
	go install github.com/golang/mock/mockgen@v1.6.0
mock-user:
	mockgen -source service/user/userRepo.go -destination service/user/mock/userMockRepo.go
mock-notification:
	mockgen -source service/notification/notificationRepo.go -destination service/notification/mock/notificationMockRepo.go

# test
test-service:
	go test -v ./service/... -coverprofile=coverage.out -cover -failfast
test-service-coverage:
	go test -v $$(go list ./service/... | grep -v '/mock') -coverprofile=coverage.out -cover -failfast && \
	go tool cover -html=coverage.out -o cover.html && \
	open cover.html