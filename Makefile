BINARY_NAME=bcrypt_cost_factor_benchmark

.PHONY: dev
dev:
	go run main.go

.PHONY: build
build: ## Builds binary 
	@ printf "Building application...\n"
	@GOARCH=amd64 GOOS=linux go build -trimpath -o bin/${BINARY_NAME}_linux_amd64 main.go
	@GOARCH=arm64 GOOS=darwin go build -trimpath -o bin/${BINARY_NAME}_darwin_arm64 main.go
	@ echo "Done"