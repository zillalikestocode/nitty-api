package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/zillalikestocode/community-api/app/configs"
	"github.com/zillalikestocode/community-api/app/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/user", loadUserRoutes)
	router.Route("/community", loadCommunityRoutes)

	return router
}

func loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{}

	// protected
	router.With(jwtauth.Verifier(configs.UseJWT())).With(jwtauth.Authenticator(configs.UseJWT())).Group(func(router chi.Router) {
		// router.Use(jwtauth.Verifier(configs.UseJWT()))
		// router.Use(jwtauth.Authenticator(configs.UseJWT()))

		router.Get("/", userHandler.Get)
	})

	router.Group(func(router chi.Router) {
		router.Post("/create", userHandler.Create)
		router.Post("/login", userHandler.Login)
	})

}

func loadCommunityRoutes(router chi.Router) {
	communityHandler := &handler.Community{}
	router.With(jwtauth.Verifier(configs.UseJWT())).With(jwtauth.Authenticator(configs.UseJWT())).Group(func(router chi.Router) {

		router.Post("/create", communityHandler.Create)
		router.Post("/join", communityHandler.Join)
		router.Post("/leave", communityHandler.Leave)
		router.Post("/announcement/create", communityHandler.CreateAnnouncement)
		router.Post("/announcement/delete", communityHandler.DeleteAnnouncement)
		router.Post("/event/create", communityHandler.CreateEvent)
		router.Post("/event/delete", communityHandler.DeleteEvent)
		router.Post("/event/update", communityHandler.UpdateEvent)
	})
}
