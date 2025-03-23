package repository

import (
	"context"

	"app/internal/model"
)

type NewsRepository interface {
	GetNewsList(ctx context.Context, keyword string) ([]*model.NewsType, error)
}

func NewNewsRepository(
	repository *Repository,
) NewsRepository {
	return &newsRepository{
		Repository: repository,
	}
}

type newsRepository struct {
	*Repository
}

func (r *newsRepository) GetNewsList(ctx context.Context, keyword string) ([]*model.NewsType, error) {
	var news []*model.News
	if keyword != "" {
		if err := r.DB(ctx).Where("title LIKE ? or content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&news).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.DB(ctx).Find(&news).Error; err != nil {
			return nil, err
		}
	}

	// Group news by type
	newsMap := make(map[int8]*model.NewsType)

	for _, item := range news {
		// Create newsResponse
		newsResponse := &model.NewsResponse{
			ID:      item.ID,
			Title:   item.Title,
			Content: item.Content,
			Img:     item.Img,
			Url:     item.Url,
			Date:    item.Date.Format("2006-01-02"),
		}

		// Get or create newsType
		if _, ok := newsMap[item.Type]; !ok {
			// Get title for this type
			var title string
			switch item.Type {
			case 1:
				title = "鸵鸟信息"
			case 2:
				title = "农场信息"
			case 3:
				title = "市场咨询"
			case 4:
				title = "供需发布"
			default:
				title = "其他信息"
			}

			newsMap[item.Type] = &model.NewsType{
				Type:  item.Type,
				Title: title,
				List:  []*model.NewsResponse{},
			}
		}

		// Append newsResponse to this type
		newsType := newsMap[item.Type]
		newsType.List = append(newsType.List, newsResponse)
	}

	// Convert map to slice
	var newsTypes []*model.NewsType
	for _, v := range newsMap {
		newsTypes = append(newsTypes, v)
	}

	return newsTypes, nil
}
