package api

import (
	"log/slog"

	"github.com/go-fuego/fuego"
	"github.com/yyewolf/rodent/mischief"
)

type CleanupRepository struct {
	mischief *mischief.Mischief
	logger   *slog.Logger
}

func NewCleanupRepository(mischief *mischief.Mischief, logger *slog.Logger) *CleanupRepository {
	return &CleanupRepository{
		mischief: mischief,
		logger:   logger,
	}
}

func (c *CleanupRepository) Group() string {
	return "/cleanup"
}

func (c *CleanupRepository) Register(server *fuego.Server) {
	fuego.Get(server, "", c.DoCleanup)
}

func (c *CleanupRepository) DoCleanup(ctx fuego.ContextNoBody) (bool, error) {
	return true, c.mischief.Cleanup(ctx)
}

var _ Repository = &CleanupRepository{}
