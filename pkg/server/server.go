package server

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tail12/prac-grpc-go/pkg/api"
	"github.com/tail12/prac-grpc-go/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Config struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBSchema   string
}

func RunServer() error {
	cfg := getConfig()
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatal("faild to listen: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBSchema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("faild to open database: %v", err)
	}
	defer db.Close()

	server := service.NewBookServiceServer(db)
	s := grpc.NewServer()

	api.RegisterBookServiceServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("starting gRPC server...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("faild to serve: %v", err)
		return err
	}

	return nil
}

func getConfig() Config {

	var cfg Config
	flag.StringVar(&cfg.Port, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	return cfg
}
