package mysql

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//Mysql Mysql
type Mysql struct {
	db      *gorm.DB
	connStr string
	debug   bool
}

//NewMysql New Mysql
func NewMysql(connStr string, debug bool) *Mysql {
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &Mysql{
		db:      db,
		connStr: connStr,
		debug:   debug,
	}
}

//NewMysql2 New Mysql
func NewMysql2(userName, password, host, dbName string, port int, debug bool) *Mysql {
	return NewMysql(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, host, port, dbName), debug)
}

//Close 关闭
func (m *Mysql) Close() {
}

//DB DB
func (m *Mysql) DB() *gorm.DB {
	if m.debug {
		return m.db.Session(&gorm.Session{}).Debug()
	}
	return m.db.Session(&gorm.Session{})
}

//PageQuery 分页查询
func (m *Mysql) PageQuery(db *gorm.DB, pageSize, pageIndex int64, order string, result interface{}) (records int64, pages int64, err error) {
	err = db.Count(&records).Error
	if err != nil {
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
	db = db.Order(order).Offset(int(offset)).Limit(int(pageSize))

	db = db.Find(result)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return
	}
	err = db.Error
	return
}

//PageQueryWithPreload 分页查询
func (m *Mysql) PageQueryWithPreload(db *gorm.DB, pageSize, pageIndex int64, order string, preload []string, result interface{}) (records int64, pages int64, err error) {
	err = db.Count(&records).Error
	if err != nil {
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
	for _, s := range preload {
		db = db.Preload(s)
	}
	db = db.Order(order)
	db = db.Find(result)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return
	}
	err = db.Error
	return
}

//PageQueryWithAssociations 分页查询
func (m *Mysql) PageQueryWithAssociations(db *gorm.DB, pageSize, pageIndex int64, order string, result interface{}) (records int64, pages int64, err error) {
	err = db.Count(&records).Error
	if err != nil {
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
	db = db.Order(order)
	db = db.Offset(int(offset)).Limit(int(pageSize))

	db = db.Preload(clause.Associations).Find(result)

	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return
	}
	err = db.Error
	return
}

func (m *Mysql) CheckDuplication(db *gorm.DB, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	fmt.Println("count=", count)
	return count > 0, nil
}

func (m *Mysql) CheckDuplicationByTableName(db *gorm.DB, tableName string, query string, args ...interface{}) (bool, error) {
	var result = make(map[string]int64)
	err := db.Raw(fmt.Sprintf("select count(1) as count from %s where %s", tableName, query), args...).Scan(&result).Error
	if err != nil {
		return true, err
	}
	return result["count"] > 0, nil
}

func (m *Mysql) CreateWithCheckDuplication(info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := m.DB().Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}

	if count > 0 {
		return true, nil
	}
	err = m.DB().Create(info).Error
	return false, err
}

func (m *Mysql) CreateWithCheckDuplicationWithDB(db *gorm.DB, info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}

	if count > 0 {
		return true, nil
	}
	err = m.DB().Create(info).Error
	return false, err
}

func (m *Mysql) CreateWithCheckDuplicationByTableName(tableName string, info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := m.DB().Table(tableName).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = m.DB().Create(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplication(info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := m.DB().Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = m.DB().Session(&gorm.Session{FullSaveAssociations: true}).Save(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplication2(db *gorm.DB, info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Save(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplicationAndOmit(info interface{}, omit []string, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := m.DB().Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = m.DB().Session(&gorm.Session{FullSaveAssociations: true}).Omit(omit...).Save(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplicationAndOmit2(db *gorm.DB, info interface{}, omit []string, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Omit(omit...).Save(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplicationByTableName(tableName string, info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := m.DB().Table(tableName).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = m.DB().Session(&gorm.Session{FullSaveAssociations: true}).Save(info).Error
	return false, err
}

func (m *Mysql) UpdateWithCheckDuplication3(db *gorm.DB, info, query interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := db.Model(info).Where(query, args...).Count(&count).Error
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Save(info).Error
	return false, err
}
