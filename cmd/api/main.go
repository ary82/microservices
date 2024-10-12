package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ary82/microservices/internal/api"
	"github.com/ary82/microservices/internal/api/gql"
	"github.com/ary82/microservices/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultPort = "8000"

func main() {
	mode := os.Getenv("MODE")
	if mode != "local" && mode != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env:", err)
		}
	}

	port := os.Getenv("API_SERVICE_PORT")
	if port == "" {
		port = defaultPort
	}

	creds := credentials.NewTLS(&tls.Config{})
	if mode != "prod" {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(
			creds,
		),
	}

	UserServiceConn, err := grpc.NewClient(os.Getenv("USERS_SERVICE_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	ProductServiceConn, err := grpc.NewClient(os.Getenv("PRODUCTS_SERVICE_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	OrderServiceConn, err := grpc.NewClient(os.Getenv("ORDERS_SERVICE_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}

	userService := proto.NewUserServiceClient(UserServiceConn)
	productService := proto.NewProductsServiceClient(ProductServiceConn)
	ordersService := proto.NewOrdersServiceClient(OrderServiceConn)

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &api.Resolver{
		UserService:    userService,
		ProductService: productService,
		OrderService:   ordersService,
	}}))

	router := http.NewServeMux()

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", api.Middleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
