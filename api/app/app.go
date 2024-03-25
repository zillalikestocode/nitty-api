package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("failed to start: %v", err)
	} else {
		fmt.Print("server running")
	}

	return err
}