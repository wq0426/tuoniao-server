package repository

import (
	"context"

	"app/internal/model"
)

type FreeMarketMineRepository interface {
	GetUserEggsSummary(ctx context.Context, userID string) (*model.FreeMarketMineResponse, error)
	UpdateEggPrice(ctx context.Context, price float64, id int) error
}

func NewFreeMarketMineRepository(
	repository *Repository,
) FreeMarketMineRepository {
	return &freeMarketMineRepository{
		Repository: repository,
	}
}

type freeMarketMineRepository struct {
	*Repository
}

func (r *freeMarketMineRepository) GetUserEggsSummary(ctx context.Context, userID string) (*model.FreeMarketMineResponse, error) {
	// Get all records for this user
	var records []*model.FreeMarketMine
	if err := r.DB(ctx).Where("user_id = ?", userID).Find(&records).Error; err != nil {
		return nil, err
	}

	response := &model.FreeMarketMineResponse{
		SelledList: []model.SelledEggDTO{},
		NoSellList: []model.NoSellEggDTO{},
		Others:     []model.NoSellEggDTO{},
		Summary:    model.SummaryDTO{},
	}

	// Process records into appropriate lists
	for _, record := range records {
		// Skip records with null values
		if record.EggPrice == 0 || record.EggNum == nil || record.Status == nil || record.Date == nil {
			continue
		}

		dateStr := record.Date.Format("2006-01-02")

		// Check if sold or unsold
		if *record.Status == 1 {
			// Sold eggs
			total := record.EggPrice * float64(*record.EggNum)

			selledItem := model.SelledEggDTO{
				EggPrice: record.EggPrice,
				EggNum:   int(*record.EggNum),
				Date:     dateStr,
				Total:    total,
			}
			response.SelledList = append(response.SelledList, selledItem)
		} else {
			// Unsold eggs
			noSellItem := model.NoSellEggDTO{
				Id:       int(record.ID),
				EggPrice: record.EggPrice,
				EggNum:   int(*record.EggNum),
				Date:     dateStr,
			}
			response.NoSellList = append(response.NoSellList, noSellItem)
		}
	}

	summary := model.SummaryDTO{
		TotalSelled: len(response.SelledList),
		TotalNoSell: len(response.NoSellList),
	}
	response.Summary = summary

	var others []*model.FreeMarketMine
	if err := r.DB(ctx).Where("user_id != ?", userID).Find(&others).Error; err != nil {
		return nil, err
	}

	for _, record := range others {
		dateStr := record.Date.Format("2006-01-02")
		response.Others = append(response.Others, model.NoSellEggDTO{
			Id:       int(record.ID),
			EggPrice: record.EggPrice,
			EggNum:   int(*record.EggNum),
			Date:     dateStr,
		})
	}

	return response, nil
}

func (r *freeMarketMineRepository) UpdateEggPrice(ctx context.Context, price float64, id int) error {
	// Update price for all unsold eggs for this user
	result := r.DB(ctx).Model(&model.FreeMarketMine{}).
		Where("status = 0 AND id = ?", id).
		Update("egg_price", price)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
