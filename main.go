package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/takumi616/ielts-vocabularies-api/adapters/handlers"
	"github.com/takumi616/ielts-vocabularies-api/adapters/presenters"
	"github.com/takumi616/ielts-vocabularies-api/adapters/repositories"
	"github.com/takumi616/ielts-vocabularies-api/infrastructures"
	"github.com/takumi616/ielts-vocabularies-api/infrastructures/database"
	"github.com/takumi616/ielts-vocabularies-api/usecases/interactors"
)

func main() {
	ctx := context.Background()

	//Get config
	config, _ := infrastructures.NewConfig()

	//Initialize postgres db with Gorm
	gorm, err := database.Open(ctx, config.PgConfig)
	if err != nil {
		log.Fatal("failed to open db: %v", err)
	}

	//Initialize repository
	vocabRepository := &repositories.VocabRepository{Persistence: gorm}

	//Initialize presenter
	vocabPresenter := &presenters.VocabPresenter{}
	errPresenter := &presenters.ErrPresenter{}

	//Initialize interactor with repository and presenter
	vocabInteractor := &interactors.VocabInteractor{Repo: vocabRepository, VocabOutputPort: vocabPresenter, ErrOutputPort: errPresenter}

	//Initialize handler with service
	vocabHandler := &handlers.VocabHandler{VocabInputPort: vocabInteractor}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /vocabularies", vocabHandler.AddNewVocabulary)
	mux.HandleFunc("GET /vocabularies", vocabHandler.FetchAllVocabularies)
	mux.HandleFunc("GET /vocabularies/{id}", vocabHandler.FetchVocabularyById)
	mux.HandleFunc("PUT /vocabularies/{id}", vocabHandler.UpdateVocabularyById)
	mux.HandleFunc("DELETE /vocabularies/{id}", vocabHandler.DeleteVocabularyById)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("APP_CONTAINER_PORT"),
		Handler: mux,
	}

	log.Fatal(srv.ListenAndServe())
}
