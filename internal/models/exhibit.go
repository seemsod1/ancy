package models

import "time"

type Exhibit struct {
	ID          int    `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	TypeID      int
	Type        ExhibitType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Description string      `gorm:"size:255"`
	AssetPath   string      `gorm:"size:255;not null"`
	PreviewPath string      `gorm:"size:255;not null"`
	StatusID    int
	Status      ExhibitStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AuthorID    int
	Author      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExhibitType struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

type ExhibitStatus struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}
