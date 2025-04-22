package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type ReportRequest struct {
	ID                uint `gorm:"primarykey"`
	TableName         string
	Query             string
	Fields            datatypes.JSON
	ReportFileType    uint8
	IsCreatedReport   bool
	ReportFileAddress string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
