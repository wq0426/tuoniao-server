package repository

import (
	"context"
	"errors"

	"app/internal/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductEvaluateRepository interface {
	CreateEvaluateReply(ctx context.Context, userID string, req model.CreateEvaluateReplyRequest) error
	CreateEvaluate(ctx context.Context, userID string, req model.CreateEvaluateRequest) error
	UpdateEvaluateAnonymous(ctx context.Context, userID string, req model.UpdateEvaluateAnonymousRequest) error
}

type productEvaluateRepository struct {
	*Repository
}

func NewProductEvaluateRepository(repository *Repository) ProductEvaluateRepository {
	return &productEvaluateRepository{
		Repository: repository,
	}
}

// CreateEvaluateReply 创建评价回复
func (r *productEvaluateRepository) CreateEvaluateReply(ctx context.Context, userID string, req model.CreateEvaluateReplyRequest) error {
	// 开启事务
	tx := r.DB(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error("开启事务失败", zap.Error(tx.Error))
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 用户是否存在
	var user model.Account
	if err := tx.Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		r.logger.Error("查询用户失败", zap.Error(err))
		return err
	}

	// 1. 查询被回复的评价是否存在
	var parentEvaluate model.ProductEvaluate
	if err := tx.Where("id = ?", req.EvaluateID).First(&parentEvaluate).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		r.logger.Error("查询评价失败", zap.Error(err))
		return err
	}

	// 2. 创建新的评价回复记录
	newEvaluate := model.ProductEvaluate{
		ParentID:  int8(req.EvaluateID), // 设置父级ID为被回复的评价ID
		ProductID: parentEvaluate.ProductID,
		ReviewID:  parentEvaluate.ReviewID,
		UserID:    userID,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Content:   req.Content,
	}

	if err := tx.Create(&newEvaluate).Error; err != nil {
		tx.Rollback()
		r.logger.Error("创建评价回复失败", zap.Error(err))
		return err
	}

	// 3. 更新评论数量 (evaluate_nums +1)
	if err := tx.Model(&model.ProductReview{}).
		Where("id = ?", parentEvaluate.ReviewID).
		Update("evaluate_nums", gorm.Expr("evaluate_nums + ?", 1)).Error; err != nil {
		tx.Rollback()
		r.logger.Error("更新评价数量失败", zap.Error(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.logger.Error("提交事务失败", zap.Error(err))
		return err
	}

	return nil
}

// CreateEvaluate 创建主评论
func (r *productEvaluateRepository) CreateEvaluate(ctx context.Context, userID string, req model.CreateEvaluateRequest) error {
	// 开启事务
	tx := r.DB(ctx).Begin()
	if tx.Error != nil {
		r.logger.Error("开启事务失败", zap.Error(tx.Error))
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 用户是否存在
	var user model.Account
	if err := tx.Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		r.logger.Error("查询用户失败", zap.Error(err))
		return err
	}

	// 1. 查询评价是否存在
	var review model.ProductReview
	if err := tx.Where("id = ?", req.ReviewID).First(&review).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评价不存在")
		}
		r.logger.Error("查询评价失败", zap.Error(err))
		return err
	}

	// 2. 创建新的评价记录
	newEvaluate := model.ProductEvaluate{
		ParentID:  0, // 设置父级ID为0，表示这是一级评论
		ProductID: review.ProductID,
		ReviewID:  req.ReviewID,
		UserID:    userID,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Content:   req.Content,
	}

	if err := tx.Create(&newEvaluate).Error; err != nil {
		tx.Rollback()
		r.logger.Error("创建评论失败", zap.Error(err))
		return err
	}

	// 3. 更新评论数量 (evaluate_nums +1)
	if err := tx.Model(&model.ProductReview{}).
		Where("id = ?", req.ReviewID).
		Update("evaluate_nums", gorm.Expr("evaluate_nums + ?", 1)).Error; err != nil {
		tx.Rollback()
		r.logger.Error("更新评价数量失败", zap.Error(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.logger.Error("提交事务失败", zap.Error(err))
		return err
	}

	return nil
}

// UpdateEvaluateAnonymous 更新评论匿名状态
func (r *productEvaluateRepository) UpdateEvaluateAnonymous(ctx context.Context, userID string, req model.UpdateEvaluateAnonymousRequest) error {
	// 查询评论是否存在并属于当前用户
	var review model.ProductReview
	if err := r.DB(ctx).Where("id = ? AND user_id = ?", req.ReviewId, userID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("评论不存在或您没有权限操作")
		}
		r.logger.Error("查询评论失败", zap.Error(err))
		return err
	}

	// 更新匿名状态
	isAnonymous := 0
	if req.IsAnonymous {
		isAnonymous = 1
	}

	if err := r.DB(ctx).Model(&model.ProductReview{}).
		Where("id = ?", req.ReviewId).
		Update("is_anonymous", isAnonymous).Error; err != nil {
		r.logger.Error("更新评论匿名状态失败", zap.Error(err))
		return err
	}

	return nil
}
