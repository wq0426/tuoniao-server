package repository

import (
    "context"
	"/internal/model"
)

type BattleRecordRepository interface {
	GetBattleRecord(ctx context.Context, id int64) (*model.BattleRecord, error)
}

func NewBattleRecordRepository(
	repository *Repository,
) BattleRecordRepository {
	return &battleRecordRepository{
		Repository: repository,
	}
}

type battleRecordRepository struct {
	*Repository
}

func (r *battleRecordRepository) GetBattleRecord(ctx context.Context, id int64) (*model.BattleRecord, error) {
	var battleRecord model.BattleRecord

	return &battleRecord, nil
}
