package repo

import (
	"log/slog"

	"github.com/eduardooliveira/stLib/v2/database"
	"github.com/eduardooliveira/stLib/v2/library/entities"
)

type AssetRepo struct {
	l *slog.Logger
}

func New() (*AssetRepo, error) {

	if err := database.DB.AutoMigrate(&entities.Asset{}); err != nil {
		slog.Error("failed to auto migrate asset", "error", err)
		return nil, err
	}

	return &AssetRepo{
		l: slog.With("module", "repo"),
	}, nil
}
