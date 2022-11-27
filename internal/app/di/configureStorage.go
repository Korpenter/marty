package di

import (
	"context"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/models"
	storage2 "github.com/Mldlr/marty/internal/app/storage"
	"github.com/samber/do"
	"go.uber.org/zap"
)

func configureStorage(i *do.Injector) {
	log := do.MustInvoke[*zap.Logger](i)
	cfg := do.MustInvoke[*config.Config](i)
	if cfg.PostgresURI != "" {
		r, err := storage2.NewPostgresRepo(cfg.PostgresURI)
		if err != nil {
			log.Fatal("Error initiating postgres connection", zap.Error(err))
		}

		err = r.Ping(context.Background())
		if err != nil {
			log.Fatal("Error reaching db", zap.Error(err))
		}

		err = r.NewTables()
		if err != nil {
			log.Fatal("Error creating tables", zap.Error(err))
		}

		do.Provide(
			i,
			func(i *do.Injector) (storage2.Repository, error) {
				return r, nil
			},
		)
		return
	}

	log.Fatal("configuring storage", zap.Error(models.ErrNoStorageSpecified))
}
