package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ReportRequest struct {
	ID                uint           `gorm:"primarykey"`
	FileName          string         `gorm:"not null;size:32"`
	TableName         string         `gorm:"not null"`
	Query             string         `gorm:"not null"`
	Fields            datatypes.JSON `gorm:"not null"`
	ReportFileType    uint8
	IsCreatedReport   bool
	StorageDriver     string
	ReportFileAddress string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
