package repository

import (
	"context"
	"sort"
	"time"

	"app/internal/model"

	"go.uber.org/zap"
)

type UserEarningRepository interface {
	AddEarning(ctx context.Context, userID string, req model.AddEarningRequest) error
	GetEarningList(ctx context.Context, userID string, req model.QueryEarningRequest) (*model.EarningListResponse, error)
}

type userEarningRepository struct {
	*Repository
}

var _earningTypeMap = map[model.EarningType]string{
	model.EarningTypeEgg:      "游戏鸵鸟蛋收益",
	model.EarningTypeBird:     "商品鸟收益",
	model.EarningTypeBreeding: "种鸟收益",
}

var _earningTypeImageMap = map[model.EarningType]string{
	model.EarningTypeEgg:      "https://txtimages.oss-cn-beijing.aliyuncs.com/category/nosell.png",
	model.EarningTypeBird:     "/images/icons/tuoniao.png",
	model.EarningTypeBreeding: "/images/icons/tuoniao.png",
}

func NewUserEarningRepository(repository *Repository) UserEarningRepository {
	return &userEarningRepository{
		Repository: repository,
	}
}

// AddEarning 添加用户收益
func (r *userEarningRepository) AddEarning(ctx context.Context, userID string, req model.AddEarningRequest) error {
	// 处理日期，如果未提供则使用当前日期
	now := time.Now()
	earningDate := now.Format("2006-01-02")

	// 创建收益记录
	earning := model.UserEarning{
		UserID:      userID,
		EarningType: req.EarningType,
		Amount:      req.Amount,
		EarningDate: earningDate,
		Year:        now.Year(),
		Month:       int(now.Month()),
		TypeName:    _earningTypeMap[req.EarningType],
		Image:       _earningTypeImageMap[req.EarningType],
	}

	// 保存到数据库
	if err := r.DB(ctx).Create(&earning).Error; err != nil {
		r.logger.Error("添加用户收益失败", zap.Error(err))
		return err
	}

	return nil
}

// GetEarningList 获取用户收益列表
func (r *userEarningRepository) GetEarningList(ctx context.Context, userID string, req model.QueryEarningRequest) (*model.EarningListResponse, error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 初始化查询
	db := r.DB(ctx).Model(&model.UserEarning{}).Where("user_id = ?", userID)

	// 如果指定了日期，则按日期筛选
	if req.Date != "" {
		db = db.Where("earning_date = ?", req.Date)
	}

	// 计算总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		r.logger.Error("获取收益记录总数失败", zap.Error(err))
		return nil, err
	}

	// 获取收益列表
	var earnings []model.UserEarning
	if err := db.Order("created_at DESC").
		// Offset((page - 1) * pageSize).
		// Limit(pageSize).
		Find(&earnings).Error; err != nil {
		r.logger.Error("获取收益列表失败", zap.Error(err))
		return nil, err
	}

	// 将列表按照类型分组
	earningsMap := make(map[model.EarningType][]model.UserEarning)
	for _, earning := range earnings {
		earningsMap[earning.EarningType] = append(earningsMap[earning.EarningType], earning)
	}
	// 在earningsMap下将列表按照日期分组
	earningsDateMap := make(map[model.EarningType]map[string][]model.UserEarning)
	for earningType, earnings := range earningsMap {
		earningsDateMap[earningType] = make(map[string][]model.UserEarning)
		for _, earning := range earnings {
			// 只取earning.EarningDate的年和月
			date := earning.EarningDate[:7] // 取前7个字符，即年月部分 (格式为: "2024-01")
			earningsDateMap[earningType][date] = append(earningsDateMap[earningType][date], earning)
		}
	}
	// 构造出EarningListResponse
	earningList := []model.EarningListItem{}
	for earningType, earningsItem := range earningsDateMap {
		// 将earningsItem整理到DateList中
		dateList := []model.EarningListDateItem{}
		// 获取日期列表并按倒序排列
		dates := make([]string, 0, len(earningsItem))
		for date := range earningsItem {
			dates = append(dates, date)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(dates)))
		// 按照倒序日期遍历
		for _, date := range dates {
			earnings := earningsItem[date]
			amount := int64(0)
			amountItem := []model.EarningListDateAmountItem{}
			for _, earning := range earnings {
				amount += earning.Amount
				amountItem = append(amountItem, model.EarningListDateAmountItem{
					Name:   earning.TypeName,
					Image:  earning.Image,
					Date:   earning.EarningDate,
					Amount: earning.Amount,
				})
			}
			dateList = append(dateList, model.EarningListDateItem{
				Date:       date,
				Amount:     amount,
				AmountItem: amountItem,
			})
		}
		earningList = append(earningList, model.EarningListItem{
			ID:       earningType,
			TypeName: _earningTypeMap[earningType],
			DateList: dateList,
		})
	}

	// 构建响应
	response := &model.EarningListResponse{
		List: earningList,
	}

	return response, nil
}
