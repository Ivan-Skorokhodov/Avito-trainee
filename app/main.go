package main

import (
	"PRmanager/internal/delivery"
	"PRmanager/internal/repository"
	"PRmanager/internal/usecase"
	"PRmanager/pkg/logs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	repo := repository.NewDatabase()
	usecase := usecase.NewUseCase(repo)
	handler := delivery.NewHandler(usecase)

	r := chi.NewRouter()
	r.Use(logs.LoggerMiddleware)

	r.Post("/team/add", handler.AddTeam)
	r.Get("/team/get", handler.GetTeam)
	r.Post("/users/setIsActive", handler.SetIsActive)

	log.Println("Servise started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
