package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"dpm/services/auth"
)

func main() {
	flag.Parse()
	cmd := flag.Arg(0)

	cfg, err := auth.Load("")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	switch cmd {
	case "create-db":
		if err := auth.CreateDB(context.Background(), cfg); err != nil {
			log.Fatalf("create database: %v", err)
		}
		fmt.Println("database created successfully")
	case "migrate-up":
		if err := auth.MigrateUp(cfg.DBURL); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		fmt.Println("migrations applied successfully")
	default:
		fmt.Fprintf(os.Stderr, "usage: migrate <create-db|migrate-up>\n")
		os.Exit(1)
	}
}
