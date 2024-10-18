package db

import "github.com/tigapilarmandiri/perkakas/common/pagination"

// Generic Data type for Query
type Query[T comparable] struct {
	Data T
	Meta pagination.Option
}
