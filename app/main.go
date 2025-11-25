package main

import (
	"PRmanager/config"
	"PRmanager/internal/delivery"
	"PRmanager/internal/repository"
	"PRmanager/internal/usecase"
	"PRmanager/pkg/logs"
	"PRmanager/pkg/panic"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()
	repo := repository.NewDatabase(cfg)
	uc := usecase.NewUseCase(repo)
	handler := delivery.NewHandler(uc, cfg)

	r := chi.NewRouter()
	r.Use(panic.PanicMiddleware)
	r.Use(logs.LoggerMiddleware)

	r.Post("/team/add", handler.AddTeam)
	r.Get("/team/get", handler.GetTeam)

	r.Post("/users/setIsActive", handler.SetIsActive)
	r.Get("/users/getReview", handler.GetReview)

	r.Post("/pullRequest/create", handler.CreatePullRequest)
	r.Post("/pullRequest/merge", handler.MergePullRequest)
	r.Post("/pullRequest/reassign", handler.Reassign)

	log.Println("Servise started on :8080")
	log.Fatal(http.ListenAndServe(handler.AppPort, r))
}
