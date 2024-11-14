package main

import (
	"context"
	"flag"
	"log"

	"github.com/ipv02/auth/internal/app"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "local.env", "path to config file")
}

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
