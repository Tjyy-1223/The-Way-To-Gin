package models

type Media struct {
	ID
	DiskType string `json:"disk_type" gorm:"size:20;index;not null;comment:存储类型"`
	ScrType  int8   `json:"src_type" gorm:"not null;comment:连接类型 1-相对路径 2-外链"`
	Src      string `json:"src" gorm:"not null;comment:资源链接"`
	Timestamp
}
