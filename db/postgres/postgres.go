package postgres

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Postgres Postgres
type Postgres struct {
	db      *gorm.DB
	connStr string
	debug   bool
}

//NewPostgres New Postgres
func NewPostgres(connStr string, debug bool) *Postgres {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Postgres{
		db:      db,
		connStr: connStr,
		debug:   debug,
	}
}

//NewPostgres2 New Postgres
func NewPostgres2(userName, password, host, dbName string, port int, debug bool) *Postgres {
	return NewPostgres(fmt.Sprintf("postgresql://%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", userName, password, host, port, dbName), debug)
}

//Close 关闭
func (m *Postgres) Close() {

}

//DB DB
func (m *Postgres) DB() *gorm.DB {
	return m.db
}

//PageQuery 分页查询
func (m *Postgres) PageQuery(db *gorm.DB, pageSize, pageIndex int64, result interface{}) (records int64, pages int64, err error) {
	db = db.Count(&records)
	if err = db.Error; err != nil {
		return
	}
	if records == 0 {
		return
	}
	if pageSize == 0 {
		pageSize = 10
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pages = records / pageSize
	if records%pageSize > 0 {
		pages++
	}

	offset := pageSize * (pageIndex - 1)
	db = db.Offset(int(offset)).Limit(int(pageSize))

	db = db.Find(result)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return
	}
	err = db.Error
	return
}
