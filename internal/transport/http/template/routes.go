package template

import (
	db "github.com/ckshitij/notification-srv/internal/db/mysql"
	"github.com/ckshitij/notification-srv/internal/domain/template"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Collection-level operations
	r.Post("/", h.CreateTemplate)

	// Resource-level operations
	r.Route("/{channel}/{name}", func(r chi.Router) {

		// Render / preview template (representation)
		r.Post("/", h.Render)

		// Versions sub-resource
		r.Route("/versions", func(r chi.Router) {
			r.Get("/", h.ListVersions)
			r.Post("/", h.AddVersion)
		})
	})

	return r
}

func NewTemplateRoutes(database *db.DB) chi.Router {
	repo := template.NewMySQLRepository(database.Conn())
	renderer := template.NewGoTemplateRenderer()
	service := template.NewService(repo, renderer)
	return NewHandler(service).Routes()
}
