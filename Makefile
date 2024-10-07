all: build

build: build-api build-user build-product build-order

build-api:
	@echo "building API..."
	@go build -o api-bin cmd/api/main.go

build-user:
	@echo "building user service..."
	@go build -o user-service-bin cmd/user-service/main.go

build-product:
	@echo "building product service..."
	@go build -o product-service-bin cmd/product-service/main.go

build-order:
	@echo "building order service..."
	@go build -o order-service-bin cmd/order-service/main.go

run-api: build-api
	@echo "running API..."
	@./api-bin

run-user: build-user
	@echo "running user service..."
	@./user-service-bin

run-product: build-product
	@echo "running product service..."
	@./product-service-bin

run-order: build-order
	@echo "running order service..."
	@./order-service-bin

clean:
	@echo "cleaning binaries..."
	@rm -f api-bin user-service-bin product-service-bin order-service-bin

.PHONY: all build build-api build-user build-product build-order run-user run-product run-order run-api clean
