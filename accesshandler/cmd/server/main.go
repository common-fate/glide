package main

import (
	"context"
	"log"

	"github.com/common-fate/common-fate/accesshandler/pkg/config"

	"github.com/common-fate/common-fate/accesshandler/pkg/server"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg config.Config
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	s, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}

	return s.Start(ctx)
}
