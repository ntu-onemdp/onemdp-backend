package semester

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

func Init(db *pgxpool.Pool) {
	repo = &SemesterRepository{db: db}
	Service = NewSemesterService(repo)

	utils.Logger.Info().Msg("Semester service initialized")
}
