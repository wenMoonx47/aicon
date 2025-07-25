package database

import (
	"context"

	"Aicon-assignment/internal/domain/entity"
)

// SQLHandler インターフェース
type SqlHandler interface {
	Execute(ctx context.Context, statement string, args ...interface{}) (Result, error)
	Query(ctx context.Context, statement string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, statement string, args ...interface{}) Row
	Close() error
}

// Result インターフェース
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Rows インターフェース
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// Row インターフェース
type Row interface {
	Scan(dest ...interface{}) error
}

// ItemRepository インターフェース
type ItemRepository interface {
	GetAll(ctx context.Context) ([]*entity.Item, error)
	GetByID(ctx context.Context, id int64) (*entity.Item, error)
	Create(ctx context.Context, item *entity.Item) (*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) error
	Delete(ctx context.Context, id int64) error
	GetSummary(ctx context.Context) (map[string]int, error)
} 