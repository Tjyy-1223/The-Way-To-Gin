package models

import (
	"gorm.io/gorm"
	"time"
)

// ID 自增主键 ID
type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// Timestamp 创建、更新时间
type Timestamp struct {
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SoftDeletes 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
