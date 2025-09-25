package services

import (
	"errors"

	"gorm.io/gorm"
)

// BaseService 基础服务，包含通用的CRUD操作
type BaseService struct {
	DB *gorm.DB
}

// NewBaseService 创建基础服务
func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{DB: db}
}

// GetByID 根据ID获取记录
func (bs *BaseService) GetByID(model interface{}, id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	result := bs.DB.First(model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("记录不存在")
		}
		return result.Error
	}

	return nil
}

// Create 创建记录
func (bs *BaseService) Create(model interface{}) error {
	result := bs.DB.Create(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update 更新记录
func (bs *BaseService) Update(model interface{}, id uint, updates interface{}) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	result := bs.DB.Model(model).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("记录不存在或未发生更改")
	}

	return nil
}

// Delete 软删除记录
func (bs *BaseService) Delete(model interface{}, id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	result := bs.DB.Delete(model, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("记录不存在")
	}

	return nil
}

// GetList 获取分页列表
func (bs *BaseService) GetList(model interface{}, page, pageSize int, conditions ...interface{}) (int64, error) {
	var total int64

	query := bs.DB.Model(model)

	// 添加查询条件
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(model).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Exists 检查记录是否存在
func (bs *BaseService) Exists(model interface{}, conditions ...interface{}) (bool, error) {
	var count int64

	query := bs.DB.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetCount 获取记录总数
func (bs *BaseService) GetCount(model interface{}, conditions ...interface{}) (int64, error) {
	var count int64

	query := bs.DB.Model(model)
	if len(conditions) > 0 {
		query = query.Where(conditions[0], conditions[1:]...)
	}

	err := query.Count(&count).Error
	return count, err
}

// BatchCreate 批量创建记录
func (bs *BaseService) BatchCreate(models interface{}) error {
	result := bs.DB.Create(models)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateColumn 更新单个字段
func (bs *BaseService) UpdateColumn(model interface{}, id uint, column string, value interface{}) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}

	result := bs.DB.Model(model).Where("id = ?", id).Update(column, value)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("记录不存在或未发生更改")
	}

	return nil
}