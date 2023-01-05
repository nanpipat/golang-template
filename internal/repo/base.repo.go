package repo

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/nanpipat/golang-template-hexagonal/helper"
	"github.com/nanpipat/golang-template-hexagonal/utils"

	"gorm.io/gorm"
)

type BaseRepository[M any] struct {
	db            *gorm.DB
	isForceUpdate bool
}

func New[M any](db *gorm.DB) *BaseRepository[M] {
	item := new(M)
	return &BaseRepository[M]{db: db.Model(item)}
}

// FindAll find records that match given conditions
func (m *BaseRepository[M]) FindAll(conds ...interface{}) ([]M, error) {
	defer m.NewSession()
	list := make([]M, 0)
	err := m.getDBInstance().Find(&list, conds...).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// FindOne find first record that match given conditions, order by primary key
func (m *BaseRepository[M]) FindOne(conds ...interface{}) (*M, error) {
	defer m.NewSession()
	item := new(M)
	err := m.getDBInstance().First(item, conds...).Error
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return nil, errors.New("not found")
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (m *BaseRepository[M]) Count() (int64, error) {
	defer m.NewSession()
	var count int64
	err := m.getDBInstance().Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create insert the value into database
func (m *BaseRepository[M]) Create(values interface{}) error {
	defer m.NewSession()
	err := m.getDBInstance().Create(values).Error
	if errors.Is(err, gorm.ErrEmptySlice) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// Update update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
func (m *BaseRepository[M]) Update(values interface{}) error {
	defer m.NewSession()
	var err error
	if m.isForceUpdate {
		var newValues map[string]interface{}
		inrec, _ := json.Marshal(values)
		_ = json.Unmarshal(inrec, &newValues)
		err = m.getDBInstance().Updates(newValues).Error
	} else {
		err = m.getDBInstance().Updates(values).Error
	}
	if errors.Is(err, gorm.ErrEmptySlice) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// Delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (m *BaseRepository[M]) Delete(conds ...interface{}) error {
	defer m.NewSession()
	item := new(M)
	err := m.getDBInstance().Delete(item, conds...).Error
	if err != nil {
		return err
	}

	return nil
}

// HardDelete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (m *BaseRepository[M]) HardDelete(conds ...interface{}) error {
	defer m.NewSession()
	item := new(M)
	err := m.getDBInstance().Unscoped().Delete(item, conds...).Error
	if err != nil {
		return err
	}

	return nil
}

type Pagination[M any] struct {
	Page  int64 `json:"page" example:"1"`
	Total int64 `json:"total" example:"45"`
	Limit int64 `json:"limit" example:"30"`
	Count int64 `json:"count" example:"30"`
	Items []M   `json:"items"`
}

func (m *BaseRepository[M]) Pagination(pageOptions *helper.PageOptions) (*Pagination[M], error) {
	defer m.NewSession()
	list := make([]M, 0)
	pageRes, err := Paginate(m.getDBInstance(), &list, pageOptions)
	if err != nil {
		return nil, err
	}

	return &Pagination[M]{
		Limit: pageRes.Limit,
		Page:  pageRes.Page,
		Total: pageRes.Total,
		Count: pageRes.Count,
		Items: list,
	}, nil
}

func Paginate(db *gorm.DB, model interface{}, options *helper.PageOptions) (*helper.PageResponse, error) {
	if options.Page == 0 {
		options.Page = 1
	}

	offset := (options.Page - 1) * options.Limit

	if len(options.OrderBy) > 0 {
		for _, o := range options.OrderBy {
			db = db.Order(o)
		}
	}

	var totalCount int64
	err := db.Model(model).Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	err = db.Limit(int(options.Limit)).Offset(int(offset)).Find(model).Error
	if err != nil {
		return nil, err
	}

	return &helper.PageResponse{
		Total: totalCount,
		Limit: options.Limit,
		Count: int64(reflect.ValueOf(model).Elem().Len()),
		Page:  options.Page,
		Q:     options.Q,
	}, nil
}

func (m *BaseRepository[M]) Upsert(values interface{}) error {
	defer m.NewSession()
	old := *m
	value := map[string]interface{}{}
	_ = utils.MapToStruct(&values, &value)
	isExist := true

	r, err := m.FindOne()
	if err.Error() == "not found" {
		isExist = false
	}

	if err != nil && !(err.Error() == "not found") {
		return err
	}
	if r != nil {
		delete(value, "id")
	}

	if isExist {
		return old.Update(value)
	} else {
		return m.Create(value)
	}
}

func (m *BaseRepository[M]) Set(values interface{}) error {
	defer m.NewSession()
	model := new(M)
	err := m.getDBInstance().Model(model).Save(values).Error
	if errors.Is(err, gorm.ErrEmptySlice) {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (m *BaseRepository[M]) getDBInstance() *gorm.DB {
	return m.db
}

func (m *BaseRepository[M]) NewSession() *BaseRepository[M] {
	m.db = m.getDBInstance().Session(&gorm.Session{NewDB: true}).Model(new(M))
	return m
}

func (m *BaseRepository[M]) Where(query interface{}, args ...interface{}) *BaseRepository[M] {
	m.db = m.db.Where(query, args...)
	return m
}

func (m *BaseRepository[M]) Preload(query string, args ...interface{}) *BaseRepository[M] {
	m.db = m.db.Preload(query, args...)
	return m
}

func (m *BaseRepository[M]) Unscoped() *BaseRepository[M] {
	m.db = m.db.Unscoped()
	return m
}

// Exec execute raw sql
func (m *BaseRepository[M]) Exec(sql string, values ...interface{}) *BaseRepository[M] {
	m.db = m.db.Exec(sql, values...)
	return m
}

func (m *BaseRepository[M]) Group(name string) *BaseRepository[M] {
	m.db = m.db.Group(name)
	return m
}

func (m *BaseRepository[M]) Joins(query string, args ...interface{}) *BaseRepository[M] {
	m.db = m.db.Joins(query, args...)
	return m
}

func (m *BaseRepository[M]) Order(value interface{}) *BaseRepository[M] {
	m.db = m.db.Order(value)
	return m
}

// Distinct specify distinct fields that you want querying
func (m *BaseRepository[M]) Distinct(args ...interface{}) *BaseRepository[M] {
	m.db = m.db.Distinct(args...)
	return m
}

// force update for zero attributes
func (m *BaseRepository[M]) ForceUpdate(value bool) *BaseRepository[M] {
	m.isForceUpdate = value
	return m
}

func PaginationMap[T any, R any](pagination *Pagination[T], cb func(r T) R) *Pagination[R] {
	return &Pagination[R]{
		Page:  pagination.Page,
		Total: pagination.Total,
		Limit: pagination.Limit,
		Count: pagination.Count,
		Items: helper.ArrayMap(pagination.Items, cb),
	}
}
