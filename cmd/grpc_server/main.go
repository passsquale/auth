package main

import (
	"context"
	"github.com/passsquale/auth/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err.Error())
	}
}
