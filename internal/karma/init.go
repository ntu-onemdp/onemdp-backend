package karma

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func Init(db *pgxpool.Pool) {
	repo := &karmaRepository{db: db}
	Service = NewKarmaService(repo)

	utils.Logger.Info().Msg("Karma service initialized")
}
