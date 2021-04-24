package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"labsystem/db"
	"labsystem/model"
	"sync"
	"time"
)

// order_by
type OrderTyp string

const (
	DESC OrderTyp = "desc"
	ASC           = "asc"
)

type orderBy struct {
	field string
	sort  OrderTyp
}

func NewOrderBy(field string, sort OrderTyp) *orderBy {
	return &orderBy{field, sort}
}

// filter
type Filter interface {
	OrderBy() string
	PageScope(db *gorm.DB) *gorm.DB
}

type BaseFilter struct {
	Sort          *orderBy
	createdEarly  *time.Time
	createdLatest *time.Time
	Page          uint
	PerPage       uint
	TotalPage     uint
	TotalCount    uint
}

func (f *BaseFilter) OrderBy() string {
	return f.Sort.field + " " + string(f.Sort.sort)
}

func (f *BaseFilter) SetCreatedAtRange(min *time.Time, max *time.Time) {
	f.createdLatest = max
	f.createdEarly = min
}

func (f *BaseFilter) PageScope(db *gorm.DB) *gorm.DB {
	if f.createdEarly != nil {
		db = db.Where("created_at > ?", f.createdEarly)
	}
	if f.createdLatest != nil {
		db = db.Where("created_at < ?", f.createdLatest)
	}
	if f.Page == 0 || f.PerPage == 0 {
		return db
	}
	var length int64
	db.Where("deleted_at is null").Count(&length)
	f.TotalCount = uint(length)
	if f.TotalCount == 0 {
		f.PerPage = 0
	} else {
		f.TotalPage = (f.TotalCount-1)/f.PerPage + 1
	}
	return db.Offset(int((f.Page - 1) * f.PerPage)).Limit(int(f.PerPage))
}

// DAO
type DAO struct {
	SQL   *gorm.DB
	cache *redis.Client
	ctx   context.Context
	once  sync.Once
}

var _db DAO

func NewDAO() *DAO {
	_db.once.Do(func() {
		_db.SQL = db.NewMySQL().DB
		_db.cache = db.NewRedis().Cli
		_db.ctx = context.Background()
	})

	return &_db
}

// clear deleted
func (d *DAO) Clear(ts []string) {
	db := d.SQL.Unscoped().Where("deleted_at is not null")
	switch ts {
	case nil:
		for _, v := range model.Models {
			db.Delete(v)
		}
	default:
		for _, v := range ts {
			db.Delete(model.Models[v])
		}
	}
}

func (d *DAO) CGet(key string) string {
	return d.cache.Get(d.ctx, key).Val()
}

func (d DAO) CSet(key, val string, exp time.Duration) error {
	return d.cache.Set(d.ctx, key, val, exp).Err()
}

func (d DAO) CHashAdd(k, f, v string) error {
	return d.cache.HSetNX(d.ctx, k, f, v).Err()
}

func (d DAO) CHashRem(k, f string) error {
	return d.cache.HDel(d.ctx, k, f).Err()
}

func (d DAO) CHashList(k string) (map[string]string, error) {
	obj := d.cache.HGetAll(d.ctx, k)
	return obj.Val(), obj.Err()
}

func (d DAO) CHashGetV(k, f string) (v string, err error) {
	obj := d.cache.HGet(d.ctx, k, f)
	return obj.Val(), obj.Err()
}

func (d DAO) CHashDelV(k, f string) error {
	return d.cache.HDel(d.ctx, k, f).Err()
}

func (d DAO) CDelete(key string) error {
	return d.cache.Del(d.ctx, key).Err()
}

func (d DAO) FlushDB() error {
	return d.cache.FlushDB(d.ctx).Err()
}

type BaseCacheDao interface {
	CGet(string) string
	CSet(k, v string, exp time.Duration) error // 0 -> forever
	CHashAdd(k, f, v string) error
	CHashRem(k, f string) error
	CHashList(k string) (map[string]string, error)
	CHashGetV(k, f string) (v string, err error)
	CHashDelV(k, f string) error
	CDelete(key string) error
}

type BaseDao interface {
	Create(interface{}) error
	Query(Filter) (interface{}, error)
	Update(map[string]interface{}, map[string]interface{}) error
	Delete(map[string]interface{}) error
	Clear()
}

type TestDao interface {
	Truncate() error
	FlushDB() error
}
