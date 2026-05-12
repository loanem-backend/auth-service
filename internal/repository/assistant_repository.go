package repository

import (
	"context"

	"github.com/loanem-backend/auth-service/infra/database/sqlc"
	"github.com/loanem-backend/auth-service/internal/entity"
)

type AssistantRepository interface {
	FindByPhone(ctx context.Context, phone string) (*entity.Assistant, error)
}

type assistantRepository struct {
	db *sqlc.Queries
}

func NewAssistantRepository(q *sqlc.Queries) AssistantRepository {
	return &assistantRepository{
		db: q,
	}
}

func (r *assistantRepository) FindByPhone(ctx context.Context, phone string) (*entity.Assistant, error) {
	row, err := r.db.FindAssistantByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	return toAssistant(row), nil
}

func toAssistant(row sqlc.Assistant) *entity.Assistant {
	return &entity.Assistant{
		ID:           int(row.ID),
		Name:         row.Name,
		Phone:        row.Phone,
		HashPassword: row.Password.String,
		Active:       row.Active,
		Period:       int(row.Period),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
