package article

import (
	"encoding/json"
	"net/http"
	"website/internal/app/response"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

func getArticlesEndpoint(w http.ResponseWriter, r *http.Request) {
	as, err := GetArticles()
	if err != nil {
		logrus.WithError(err).Error("Failed executing template")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&as)
}

func getArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	a, err := GetArticle(mux.Vars(r)["article"])
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Write404(w)
			return
		}
		logrus.WithError(err).Error("Failed retrieving article")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&a)
}

func createArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	var a Article
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		response.Write400(w)
		return
	}

	err = a.Save()
	if err != nil {
		logrus.WithError(err).Error("Failed saving article")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&a)
}

func updateArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	a, err := GetArticle(mux.Vars(r)["article"])
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Write404(w)
			return
		}
		logrus.WithError(err).Error("Failed retrieving article")
		response.Write500(w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		response.Write400(w)
		return
	}

	err = a.Save()
	if err != nil {
		logrus.WithError(err).Error("Failed saving article")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&a)
}

func deleteArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	a, err := GetArticle(mux.Vars(r)["article"])
	if err != nil {
		if err == pgx.ErrNoRows {
			response.Write404(w)
			return
		}
		logrus.WithError(err).Error("Failed retrieving article")
		response.Write500(w)
		return
	}
	err = a.Delete()
	if err != nil {
		logrus.WithError(err).Error("Failed deleting article")
		response.Write500(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}