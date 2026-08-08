package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/loanem-backend/auth-service/infra/database/sqlc"
	"github.com/loanem-backend/auth-service/internal/entity"
)

type AssistantRepository interface {
	FindByID(ctx context.Context, id int) (*entity.Assistant, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Assistant, error)
	Insert(ctx context.Context, a *entity.Assistant) (int16, error)
	UpdatePassword(ctx context.Context, a *entity.Assistant) error

	FindActiveAssistants(ctx context.Context) ([]*entity.Assistant, error)
}

type assistantRepository struct {
	db *sqlc.Queries
}

func NewAssistantRepository(q *sqlc.Queries) AssistantRepository {
	return &assistantRepository{
		db: q,
	}
}

func (r *assistantRepository) FindByID(ctx context.Context, id int) (*entity.Assistant, error) {
	row, err := r.db.FindAssistantByID(ctx, int16(id))
	if err != nil {
		return nil, err
	}

	return toAssistant(row), nil
}

func (r *assistantRepository) FindByPhone(ctx context.Context, phone string) (*entity.Assistant, error) {
	row, err := r.db.FindAssistantByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFindByPhoneNotFound
		}
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

func (r *assistantRepository) Insert(ctx context.Context, a *entity.Assistant) (int16, error) {
	result, err := r.db.InsertAssistant(ctx, sqlc.InsertAssistantParams{
		Name:     a.Name,
		Phone:    a.Phone,
		Password: pgtype.Text{String: a.HashPassword, Valid: true},
		Period:   int16(a.Period),
	})
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *assistantRepository) UpdatePassword(ctx context.Context, a *entity.Assistant) error {
	if err := r.db.SetAssistantPassword(ctx, sqlc.SetAssistantPasswordParams{
		ID:        int16(a.ID),
		Password:  pgtype.Text{String: a.HashPassword, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: a.UpdatedAt, Valid: true},
	}); err != nil {
		return err
	}

	return nil
}

func (r *assistantRepository) FindActiveAssistants(ctx context.Context) ([]*entity.Assistant, error) {
	rows, err := r.db.FindActiveAssistants(ctx)
	if err != nil {
		return []*entity.Assistant{}, err
	}

	assistants := make([]*entity.Assistant, len(rows))
	for i, row := range rows {
		assistants[i] = toAssistant(row)
	}

	return assistants, nil
}
