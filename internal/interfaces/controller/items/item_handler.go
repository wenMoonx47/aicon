package items

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"Aicon-assignment/internal/domain/errors"
	"Aicon-assignment/internal/usecase"
)

type ItemHandler struct {
	itemUsecase *usecase.ItemUsecase
}

func NewItemHandler(itemUsecase *usecase.ItemUsecase) *ItemHandler {
	return &ItemHandler{
		itemUsecase: itemUsecase,
	}
}

// CreateItemRequest はアイテム作成リクエストの構造体です
type CreateItemRequest struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Brand         string `json:"brand"`
	PurchasePrice int    `json:"purchase_price"`
	PurchaseDate  string `json:"purchase_date"`
}

// PatchItemRequest はアイテム部分更新リクエストの構造体です
type PatchItemRequest struct {
	Name          *string `json:"name,omitempty"`
	Brand         *string `json:"brand,omitempty"`
	PurchasePrice *int    `json:"purchase_price,omitempty"`
}

// ErrorResponse はエラーレスポンスの構造体です
type ErrorResponse struct {
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

// GetItems は全てのアイテムを取得します
func (h *ItemHandler) GetItems(c echo.Context) error {
	ctx := c.Request().Context()

	items, err := h.itemUsecase.GetAllItems(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get items",
		})
	}

	return c.JSON(http.StatusOK, items)
}

// GetItem は指定されたIDのアイテムを取得します
func (h *ItemHandler) GetItem(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid id format",
		})
	}

	item, err := h.itemUsecase.GetItemByID(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get item",
		})
	}

	return c.JSON(http.StatusOK, item)
}

// CreateItem は新しいアイテムを作成します
func (h *ItemHandler) CreateItem(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
	}

	item, err := h.itemUsecase.CreateItem(ctx, req.Name, req.Category, req.Brand, req.PurchasePrice, req.PurchaseDate)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			details := h.parseValidationErrors(err.Error())
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation failed",
				Details: details,
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to create item",
		})
	}

	return c.JSON(http.StatusCreated, item)
}

// PatchItem は既存のアイテムを部分更新します
func (h *ItemHandler) PatchItem(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid id format",
		})
	}

	var req PatchItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid request body",
		})
	}

	// リクエストボディを map に変換
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Brand != nil {
		updates["brand"] = *req.Brand
	}
	if req.PurchasePrice != nil {
		updates["purchase_price"] = *req.PurchasePrice
	}

	if len(updates) == 0 {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "no fields provided for update",
		})
	}

	item, err := h.itemUsecase.PatchItem(ctx, id, updates)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		if strings.Contains(err.Error(), "validation failed") {
			details := h.parseValidationErrors(err.Error())
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "validation failed",
				Details: details,
			})
		}
		if strings.Contains(err.Error(), "not allowed for update") {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to update item",
		})
	}

	return c.JSON(http.StatusOK, item)
}

// DeleteItem は指定されたIDのアイテムを削除します
func (h *ItemHandler) DeleteItem(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid id format",
		})
	}

	err = h.itemUsecase.DeleteItem(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "item not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to delete item",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetSummary はカテゴリー別のアイテム数とトータル数を取得します
func (h *ItemHandler) GetSummary(c echo.Context) error {
	ctx := c.Request().Context()

	summary, err := h.itemUsecase.GetItemsSummary(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get summary",
		})
	}

	return c.JSON(http.StatusOK, summary)
}

// parseValidationErrors はバリデーションエラーメッセージを分割します
func (h *ItemHandler) parseValidationErrors(errMsg string) []string {
	if strings.Contains(errMsg, "validation failed: ") {
		errMsg = strings.TrimPrefix(errMsg, "validation failed: ")
	}
	
	// カンマで分割してトリム
	details := strings.Split(errMsg, ", ")
	for i, detail := range details {
		details[i] = strings.TrimSpace(detail)
	}
	
	return details
} 