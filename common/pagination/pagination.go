package pagination

import (
	"math"
	"strings"
	"sync"

	"github.com/tigapilarmandiri/perkakas"
	"github.com/tigapilarmandiri/perkakas/common/util"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Pagination struct {
	Model     any
	Option    *Option
	DBConn    *gorm.DB
	tableName string
}

type Option struct {
	Limit      int
	Page       int
	Sort       string
	Filter     string
	TotalRows  int64
	TotalPages int
}

func (p *Pagination) getOffset() int {
	return (p.getPage() - 1) * p.getLimit()
}

func (p *Pagination) getLimit() int {
	if p.Option.Limit == 0 {
		p.Option.Limit = 10
	}
	return p.Option.Limit
}

func (p *Pagination) getPage() int {
	if p.Option.Page == 0 {
		p.Option.Page = 1
	}
	return p.Option.Page
}

func (p *Pagination) getSort() string {
	if IsSortSave(p.Option.Sort) {
		return p.tableName + "." + p.Option.Sort
	}
	return p.tableName + ".Id ASC"
}

func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	var err error
	var args []any
	var query string

	if p.Model != nil {
		var s *schema.Schema
		s, err = schema.Parse(p.Model, &sync.Map{}, schema.NamingStrategy{})
		if err != nil {
			return nil
		}
		p.tableName = s.Table
	}

	if p.DBConn.Statement.Model == nil && p.DBConn.Statement.Table == "" {
		p.DBConn = p.DBConn.Model(p.Model)
	}
	if !perkakas.IsEmpty(p.Option.Filter) {
		query, args, err = util.BuildFilterQuery(p.Model, p.Option.Filter)
		if err == nil {
			p.DBConn.Where(query, args...).Count(&totalRows)
		} else {
			p.DBConn.Count(&totalRows)
		}
	} else {
		p.DBConn.Count(&totalRows)
	}

	p.Option.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.Option.Limit)))
	p.Option.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		if !perkakas.IsEmpty(p.Option.Filter) && err == nil {
			db = db.Where(query, args...)
		}

		return db.Offset(p.getOffset()).Limit(p.getLimit()).Order(p.getSort())
	}
}

func IsSortSave(payload string) bool {
	if len(payload) > 30 || len(payload) == 0 {
		return false
	}

	payload = strings.ToLower(payload)

	if strings.HasSuffix(payload, "desc") {
		arrs := strings.Split(payload, " ")
		if len(arrs) != 2 {
			return false
		}
		payload = arrs[0]
	}

	if strings.HasSuffix(payload, "asc") {
		arrs := strings.Split(payload, " ")
		if len(arrs) != 2 {
			return false
		}
		payload = arrs[0]
	}

	for _, v := range payload {
		if (v < 'a' || v > 'z') && v != '_' {
			return false
		}
	}
	return true
}
