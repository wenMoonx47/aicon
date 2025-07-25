package usecase

import (
	"context"
	"fmt"

	"Aicon-assignment/internal/domain/entity"
	"Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/interfaces/database"
)

type ItemUsecase struct {
	itemRepo database.ItemRepository
}

func NewItemUsecase(itemRepo database.ItemRepository) *ItemUsecase {
	return &ItemUsecase{
		itemRepo: itemRepo,
	}
}

// GetAllItems は全てのアイテムを取得します
func (u *ItemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	items, err := u.itemRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all items: %w", err)
	}
	return items, nil
}

// GetItemByID は指定されたIDのアイテムを取得します
func (u *ItemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	item, err := u.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item by id: %w", err)
	}
	return item, nil
}

// CreateItem は新しいアイテムを作成します
func (u *ItemUsecase) CreateItem(ctx context.Context, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error) {
	// ドメインエンティティでバリデーション
	item, err := entity.NewItem(name, category, brand, purchasePrice, purchaseDate)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// データベースに保存
	createdItem, err := u.itemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return createdItem, nil
}

// UpdateItem は既存のアイテムを完全更新します
func (u *ItemUsecase) UpdateItem(ctx context.Context, id int64, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	// 既存のアイテムを取得
	existingItem, err := u.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing item: %w", err)
	}

	// アイテムを更新（バリデーション含む）
	if err := existingItem.Update(name, category, brand, purchasePrice, purchaseDate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// データベースに保存
	if err := u.itemRepo.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return existingItem, nil
}

// PatchItem は既存のアイテムを部分更新します（PATCH用）
func (u *ItemUsecase) PatchItem(ctx context.Context, id int64, updates map[string]interface{}) (*entity.Item, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	// 既存のアイテムを取得
	existingItem, err := u.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing item: %w", err)
	}

	// 更新可能なフィールドのみ処理
	allowedFields := map[string]bool{
		"name":           true,
		"brand":          true,
		"purchase_price": true,
	}

	hasUpdates := false
	for field := range updates {
		if !allowedFields[field] {
			return nil, fmt.Errorf("field '%s' is not allowed for update", field)
		}
		hasUpdates = true
	}

	if !hasUpdates {
		return nil, fmt.Errorf("no valid fields provided for update")
	}

	// フィールドを選択的に更新
	newName := existingItem.Name
	newBrand := existingItem.Brand
	newPurchasePrice := existingItem.PurchasePrice

	if name, ok := updates["name"]; ok {
		if nameStr, ok := name.(string); ok {
			newName = nameStr
		} else {
			return nil, fmt.Errorf("invalid type for name field")
		}
	}

	if brand, ok := updates["brand"]; ok {
		if brandStr, ok := brand.(string); ok {
			newBrand = brandStr
		} else {
			return nil, fmt.Errorf("invalid type for brand field")
		}
	}

	if price, ok := updates["purchase_price"]; ok {
		if priceFloat, ok := price.(float64); ok {
			newPurchasePrice = int(priceFloat)
		} else if priceInt, ok := price.(int); ok {
			newPurchasePrice = priceInt
		} else {
			return nil, fmt.Errorf("invalid type for purchase_price field")
		}
	}

	// バリデーション付きで更新（変更されていないフィールドは既存の値を使用）
	if err := existingItem.Update(newName, existingItem.Category, newBrand, newPurchasePrice, existingItem.PurchaseDate); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// データベースに保存
	if err := u.itemRepo.Update(ctx, existingItem); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return existingItem, nil
}

// DeleteItem は指定されたIDのアイテムを削除します
func (u *ItemUsecase) DeleteItem(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	if err := u.itemRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

// GetItemsSummary はカテゴリー別のアイテム数とトータル数を取得します
func (u *ItemUsecase) GetItemsSummary(ctx context.Context) (map[string]interface{}, error) {
	summary, err := u.itemRepo.GetSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get items summary: %w", err)
	}

	total := 0
	for _, count := range summary {
		total += count
	}

	result := map[string]interface{}{
		"categories": summary,
		"total":      total,
	}

	return result, nil
} 