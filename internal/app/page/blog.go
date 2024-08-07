package page

import (
	"database/sql"
	"net/http"
	"website/internal/app/article"
	"website/internal/app/response"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func renderBlog(w http.ResponseWriter, r *http.Request) {
	as, err := article.GetArticles()
	if err != nil {
		logrus.WithError(err).Error("Failed executing template")
		response.Write500(w)
		return
	}

	render(w, "page", pageInput{
		URI: r.URL.Path,
		Title: "Blog",
		CSSFile: "blog",
		PageTemplateName: "blog",
		PageTemplateData: as,
	})
}

func renderArticle(w http.ResponseWriter, r *http.Request) {
	a, err := article.GetArticle(mux.Vars(r)["article"])
	if err != nil {
		if err == sql.ErrNoRows {
			response.Write404(w)
			return
		}
		logrus.WithError(err).Error("Failed executing template")
		response.Write500(w)
		return
	}
	if !a.IsReleased() {
		response.Write404(w)
		return
	}

	render(w, "page", pageInput{
		URI: r.URL.Path,
		Title: a.Title,
		CSSFile: "article",
		PageTemplateName: "article",
		PageTemplateData: a,
	})
}