package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tigapilarmandiri/perkakas"
	"gorm.io/gorm/schema"
)

var allowedClause = map[string]string{
	"eq":        " = ?",
	"neq":       " != ?",
	"like":      " ILIKE ?",
	"startWith": " ILIKE ?",
	"endWith":   " ILIKE ?",
	"in":        " IN ?",
	"gt":        " > ? ",
	"gte":       " >= ? ",
	"lt":        " < ? ",
	"lte":       " <= ? ",
}

var listDatetime = []string{"tanggal", "dari_tanggal", "sampai_tanggal"}

func validateField(m any, fileds ...string) error {
	var tags []string
	val := reflect.ValueOf(m)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		field := t.Tag.Get("filter")

		tags = append(tags, field)
	}

	for _, f := range fileds {
		match := false

		// avoid if filter is empty or field equal "-"
		if len(f) > 0 && f != "-" {
			for _, t := range tags {
				if f == t {
					match = true
					break
				}
			}
		}

		if !match {
			return fmt.Errorf("invalid query filter for '%s'", f)
		}

	}

	return nil
}

func buildClause(field, clause, value string) (query string, arg any, err error) {
	c, ok := allowedClause[clause]
	if !ok {
		err = fmt.Errorf("invalid clause")
		return
	}

	query = field + c
	arg = value

	// var i int
	switch clause {
	case "like":
		arg = "%" + value + "%"
	case "startWith":
		arg = value + "%"
	case "endWith":
		arg = "%" + value
	case "eq":
		if value == "null" {
			query = field + " is null "
			arg = nil
		}
		// i, err = strconv.Atoi(value)
		// if err == nil {
		// 	arg = i
		// }
		// err = nil
		for _, v := range listDatetime {
			if len(field) < len(v) {
				continue
			}
			if v == field[len(field)-len(v):] {
				// field is table_name.field
				i, err := strconv.Atoi(value)
				if err != nil {
					return "", nil, err
				}
				t := time.UnixMilli(int64(i))
				arg = t.Format(time.RFC3339)
			}
		}
	case "in":
		s := strings.Split(value, ",")
		arg = s
		// _, err = strconv.Atoi(s[0])
		// if err == nil {
		// 	arrInt := make([]int, 0, len(s))
		// 	for _, v := range s {
		// 		i, err = strconv.Atoi(v)
		// 		if err != nil {
		// 			err = fmt.Errorf("invalid values")
		// 			return
		// 		}
		//
		// 		arrInt = append(arrInt, i)
		// 	}
		// 	arg = arrInt
		// } else {
		// 	err = nil
		// }
	case "gt", "gte", "lt", "lte":
		for _, v := range listDatetime {
			// field is table_name.field
			if v == field[len(field)-len(v):] {
				i, err := strconv.Atoi(value)
				if err != nil {
					return "", nil, err
				}
				t := time.UnixMilli(int64(i))
				arg = t.Format(time.RFC3339)
			}
		}

	}

	return
}

func getFilter(filter string) (filters []string, cond string) {
	// default
	filters = append(filters, filter)

	if strings.Contains(filter, ";AND;") {
		filters = strings.Split(filter, ";AND;")
		cond = "AND"
	} else if strings.Contains(filter, ";OR;") {
		filters = strings.Split(filter, ";OR;")
		cond = "OR"
	}

	return
}

// @m is a struct that contain filter tag
// @filter is query filter from client
// ex: name:eq:omama
func BuildFilterQuery(m any, filter string) (query string, args []any, err error) {
	filters, cond := getFilter(filter)
	var s *schema.Schema
	s, err = schema.Parse(m, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return
	}
	tableName := s.Table

	for _, f := range filters {
		qf := strings.Split(f, ":")
		if len(qf) != 3 {
			err = fmt.Errorf("invalid query filter")
			return
		}

		field, clause, value := qf[0], qf[1], qf[2]
		if err = validateField(m, field); err != nil {
			return
		}

		var q string
		var a any
		q, a, err = buildClause(tableName+"."+field, clause, value)
		if err != nil {
			return
		}

		query = fmt.Sprintf("%s %s %s", query, q, cond)
		if a != nil {
			args = append(args, a)
		}
	}

	// Cleanup condition at end of query
	if !perkakas.IsEmpty(cond) {
		query = strings.TrimSuffix(query, cond)
	}

	query = strings.Trim(query, " ")

	return
}

func GetQueryFilter(queryUrl, find string) (field string, operator string, value string) {
	filters, _ := getFilter(queryUrl)
	for _, v := range filters {
		qf := strings.Split(v, ":")
		if len(qf) != 3 {
			continue
		}
		if qf[0] == find {
			field = qf[0]
			operator = qf[1]
			value = qf[2]
			return
		}
	}

	return
}
