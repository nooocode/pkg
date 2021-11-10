package postgres

import (
	"testing"

	"gorm.io/gorm"
)

var pg *Postgres

func init() {
	pg = NewPostgres("postgresql://admin:123456@127.0.0.1:5432/mock", true)
}

func TestAutoMigrate(t *testing.T) {
	pg.db.AutoMigrate(&Location{}, &LocationTag{}, &LocationFile{})
}

func TestAdd(t *testing.T) {
	l := &Location{
		Name:     "B01",
		Code:     "1001",
		ParentID: 1,
		Tags: []LocationTag{
			{
				Name: "tag1",
			}, {
				Name: "tag2",
			},
		},
	}
	err := pg.db.Create(l).Error

	if err != nil {
		t.Fatal(err)
	}

	l.Name = "mock_update"
	l.Tags[0].Name = "tag3"
	l.Tags = append(l.Tags, LocationTag{
		Name: "tag4",
	})
	// 同时更新关联表
	err = pg.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(l).Error
	if err != nil {
		t.Fatal(err)
	}
}

type Location struct {
	gorm.Model
	Name     string
	Code     string
	ParentID uint
	Tags     []LocationTag
	Files    []LocationFile
}

type LocationTag struct {
	gorm.Model
	LocationID uint
	Name       string
}

type LocationFile struct {
	gorm.Model
	LocationID uint
	FileID     string
}
