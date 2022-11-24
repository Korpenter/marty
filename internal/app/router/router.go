package router

import (
	"github.com/Mldlr/mart/marty/internal/app/config"
	"github.com/Mldlr/mart/marty/internal/app/container"
	handlers2 "github.com/Mldlr/mart/marty/internal/app/router/handlers"
	"github.com/Mldlr/mart/marty/internal/app/router/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func NewRouter() chi.Router {
	cfg := container.Container.Get("cfg").(*config.Config)
	tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "application/json", "text/plain"))
	r.Use(middleware.Decompress)
	r.Group(func(r chi.Router) {
		r.Use(chiMiddleware.AllowContentType("application/json"))
		r.Post("/api/user/register", handlers2.Register)
		r.Post("/api/user/login", handlers2.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		//r.Get("/api/user/orders", )
		//r.Get("/api/user/balance", )
		//r.Get("/api/user/balance/withdrawals", )

		r.Group(func(r chi.Router) {
			r.Use(chiMiddleware.AllowContentType("text/plain"))
			r.Post("/api/user/orders", handlers2.AddOrder)
		})

		r.Group(func(r chi.Router) {
			r.Use(chiMiddleware.AllowContentType("application/json"))
			//r.Post("/api/user/balance/withdraw", )
		})
	})
	return r
}
