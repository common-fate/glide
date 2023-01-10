package main

import (
	"context"
	"log"

	"github.com/common-fate/common-fate/governance/pkg/server"
	"github.com/common-fate/common-fate/pkg/config"
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
	_ = godotenv.Load("../../.env")

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
