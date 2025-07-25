package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"Aicon-assignment/internal/domain/entity"
	"Aicon-assignment/internal/domain/errors"
)

type MySQLItemRepository struct {
	SqlHandler SqlHandler
}

// GetAll は全てのアイテムを取得します
func (r *MySQLItemRepository) GetAll(ctx context.Context) ([]*entity.Item, error) {
	query := `
		SELECT id, name, category, brand, purchase_price, purchase_date, created_at, updated_at
		FROM items
		ORDER BY id ASC
	`

	rows, err := r.SqlHandler.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []*entity.Item
	for rows.Next() {
		item, err := r.scanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}

	return items, nil
}

// GetByID は指定されたIDのアイテムを取得します
func (r *MySQLItemRepository) GetByID(ctx context.Context, id int64) (*entity.Item, error) {
	query := `
		SELECT id, name, category, brand, purchase_price, purchase_date, created_at, updated_at
		FROM items
		WHERE id = ?
	`

	row := r.SqlHandler.QueryRow(ctx, query, id)
	item, err := r.scanItemFromRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to scan item: %w", err)
	}

	return item, nil
}

// Create は新しいアイテムを作成します
func (r *MySQLItemRepository) Create(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	query := `
		INSERT INTO items (name, category, brand, purchase_price, purchase_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.SqlHandler.Execute(ctx, query,
		item.Name, item.Category, item.Brand, item.PurchasePrice, item.PurchaseDate,
		now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	item.ID = id
	item.CreatedAt = now
	item.UpdatedAt = now

	return item, nil
}

// Update は既存のアイテムを更新します
func (r *MySQLItemRepository) Update(ctx context.Context, item *entity.Item) error {
	query := `
		UPDATE items 
		SET name = ?, category = ?, brand = ?, purchase_price = ?, purchase_date = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.SqlHandler.Execute(ctx, query,
		item.Name, item.Category, item.Brand, item.PurchasePrice, item.PurchaseDate,
		now, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrItemNotFound
	}

	item.UpdatedAt = now
	return nil
}

// Delete は指定されたIDのアイテムを削除します
func (r *MySQLItemRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM items WHERE id = ?"

	result, err := r.SqlHandler.Execute(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrItemNotFound
	}

	return nil
}

// GetSummary はカテゴリー別のアイテム数を取得します
func (r *MySQLItemRepository) GetSummary(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT category, COUNT(*) as count
		FROM items
		GROUP BY category
	`

	rows, err := r.SqlHandler.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query summary: %w", err)
	}
	defer rows.Close()

	summary := make(map[string]int)
	
	// 全カテゴリーを0で初期化
	for _, category := range entity.GetValidCategories() {
		summary[category] = 0
	}

	// 実際のデータで上書き
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("failed to scan summary row: %w", err)
		}
		summary[category] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate summary rows: %w", err)
	}

	return summary, nil
}

// scanItem はRowsから1つのアイテムをスキャンします
func (r *MySQLItemRepository) scanItem(rows Rows) (*entity.Item, error) {
	var item entity.Item
	var createdAt, updatedAt time.Time
	var purchaseDate time.Time

	err := rows.Scan(
		&item.ID, &item.Name, &item.Category, &item.Brand,
		&item.PurchasePrice, &purchaseDate,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt
	item.PurchaseDate = purchaseDate.Format("2006-01-02")

	return &item, nil
}

// scanItemFromRow はRowから1つのアイテムをスキャンします
func (r *MySQLItemRepository) scanItemFromRow(row Row) (*entity.Item, error) {
	var item entity.Item
	var createdAt, updatedAt time.Time
	var purchaseDate time.Time

	err := row.Scan(
		&item.ID, &item.Name, &item.Category, &item.Brand,
		&item.PurchasePrice, &purchaseDate,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt
	item.PurchaseDate = purchaseDate.Format("2006-01-02")

	return &item, nil
} 