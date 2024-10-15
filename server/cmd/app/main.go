package main

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
	"github.com/shlmvgleb/mtg-task/server/internal/config"
	"github.com/shlmvgleb/mtg-task/server/internal/database"
	"github.com/shlmvgleb/mtg-task/server/internal/repo"
	log "github.com/sirupsen/logrus"
)

func handleConnection(ctx context.Context, repo *repo.Repository, conn net.Conn) {
	defer conn.Close()
	log.Infof("client connected: %v", conn.RemoteAddr().String())

	buffer := make([]byte, 10240)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			log.Infof("client disconnected: %v", conn.RemoteAddr().String())
			break
		} else if err != nil {
			log.Errorf("error reading from connection: %v\n", err)
			break
		}

		log.Infof("new data accepted: %v...", string(buffer[:n])[:10])

		socketID := fmt.Sprintf("%s_%s", conn.RemoteAddr().String(), uuid.NewString())
		if err := repo.InsertClientData(ctx, string(buffer[:n]), socketID); err != nil {
			log.Errorf("failed to insert data: %v", err)
		}
	}
}

func main() {
	config := config.ReadFromEnv()

	ctx := context.Background()
	db, err := database.New(ctx, config.Postgres)
	if err != nil {
		log.Fatalf("could not create database connection: %v", err)
	}

	repo := repo.NewPostgresRepo(db)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}

	defer listener.Close()

	log.Infof("server started on :%d", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(ctx, repo, conn)
	}
}
