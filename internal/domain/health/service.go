package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB *pgxpool.Pool
}

func NewService(DB *pgxpool.Pool) *Service {
	return &Service{DB: DB}
}

func (s *Service) CheckDB(ctx context.Context) error {
	row := s.DB.QueryRow(ctx, "SELECT 1")
	var tmp int
	return row.Scan(&tmp)
}
