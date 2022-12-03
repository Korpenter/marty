package router

import (
	"github.com/Mldlr/marty/internal/app/config"
	controllers "github.com/Mldlr/marty/internal/app/controllers/order"
	"github.com/Mldlr/marty/internal/app/controllers/user"
	"github.com/Mldlr/marty/internal/app/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/samber/do"
)

func NewRouter(i *do.Injector) chi.Router {
	orderController := do.MustInvoke[*controllers.OrderController](i)
	userController := do.MustInvoke[*user.UserController](i)
	cfg := do.MustInvoke[*config.Config](i)
	tokenAuth := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.AllowContentEncoding("gzip"))
	r.Use(chiMiddleware.Compress(5, "application/json", "text/plain"))
	r.Use(middleware.Decompress)
	r.Group(func(r chi.Router) {
		r.Use(chiMiddleware.AllowContentType("application/json"))
		r.Use(middleware.Unauthorized)
		r.Post("/api/user/register", userController.Register)
		r.Post("/api/user/login", userController.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middleware.Authenticator)
		r.Get("/api/user/orders", orderController.OrdersByUser)
		r.Get("/api/user/balance", userController.Balance)
		r.Get("/api/user/withdrawals", userController.UserWithdrawals)

		r.Group(func(r chi.Router) {
			r.Use(chiMiddleware.AllowContentType("text/plain"))
			r.Post("/api/user/orders", orderController.AddOrder)
		})

		r.Group(func(r chi.Router) {
			r.Use(chiMiddleware.AllowContentType("application/json"))
			r.Post("/api/user/balance/withdraw", userController.Withdraw)
		})
	})
	return r
}
