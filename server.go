package main

import (
	"log"
	"time"

	routes "gateway/src/api/routes"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func createGRPCConnection(serviceAddress string) *grpc.ClientConn {
    var conn *grpc.ClientConn
    var err error

    for {
        conn, err = grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
        if err == nil && conn.GetState() == connectivity.Ready {
            break
        }

        log.Printf("Waiting for connection to %s to become ready", serviceAddress)
        time.Sleep(1 * time.Second) // Exponential backoff could be implemented here
    }

    return conn
}

func main() {
    if err := env.Load(".env"); err != nil {
        log.Fatalf("Failed to load environment variables: %v", err)
    }

    authAddress := env.Get("AUTH_IP", "")
    embeddingAddress := env.Get("EMBEDDING_IP", "")

    authConn := createGRPCConnection(authAddress)
    defer authConn.Close()

    embeddingConn := createGRPCConnection(embeddingAddress)
    defer embeddingConn.Close()

    r := gin.Default()

    authRouter := routes.NewAuthRouter(authConn)
    authRouter.RegisterRoutes(r)

    embeddingRouter := routes.NewEmbeddingRouter(embeddingConn)
    embeddingRouter.RegisterRoutes(r)

    r.Run(":8080")
}
