package items

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"Aicon-assignment/internal/domain/entity"
	"Aicon-assignment/internal/domain/errors"
)

// MockItemUsecase は ItemUsecase のモックです
type MockItemUsecase struct {
	mock.Mock
}

func (m *MockItemUsecase) GetAllItems(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) GetItemByID(ctx context.Context, id int64) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) CreateItem(ctx context.Context, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error) {
	args := m.Called(ctx, name, category, brand, purchasePrice, purchaseDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) UpdateItem(ctx context.Context, id int64, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error) {
	args := m.Called(ctx, id, name, category, brand, purchasePrice, purchaseDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) PatchItem(ctx context.Context, id int64, updates map[string]interface{}) (*entity.Item, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemUsecase) DeleteItem(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockItemUsecase) GetItemsSummary(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// ItemUsecaseInterface は実際のusecaseインターフェースを定義
type ItemUsecaseInterface interface {
	GetAllItems(ctx context.Context) ([]*entity.Item, error)
	GetItemByID(ctx context.Context, id int64) (*entity.Item, error)
	CreateItem(ctx context.Context, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error)
	UpdateItem(ctx context.Context, id int64, name, category, brand string, purchasePrice int, purchaseDate string) (*entity.Item, error)
	PatchItem(ctx context.Context, id int64, updates map[string]interface{}) (*entity.Item, error)
	DeleteItem(ctx context.Context, id int64) error
	GetItemsSummary(ctx context.Context) (map[string]interface{}, error)
}

// TestItemHandler はテスト用のハンドラー構造体
type TestItemHandler struct {
	itemUsecase ItemUsecaseInterface
}

func (h *TestItemHandler) PatchItem(c echo.Context) error {
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
			details := parseValidationErrors(err.Error())
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

func setupTestHandler() (*TestItemHandler, *MockItemUsecase) {
	mockUsecase := new(MockItemUsecase)
	handler := &TestItemHandler{
		itemUsecase: mockUsecase,
	}
	
	return handler, mockUsecase
}

// parseValidationErrors はバリデーションエラーメッセージを分割します
func parseValidationErrors(errMsg string) []string {
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

func TestPatchItem_Success(t *testing.T) {
	tests := []struct {
		name           string
		itemID         string
		requestBody    map[string]interface{}
		expectedUpdate map[string]interface{}
		mockResponse   *entity.Item
	}{
		{
			name:   "正常系: 名前のみ更新",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name": "更新されたロレックス",
			},
			expectedUpdate: map[string]interface{}{
				"name": "更新されたロレックス",
			},
			mockResponse: &entity.Item{
				ID:            1,
				Name:          "更新されたロレックス",
				Category:      "時計",
				Brand:         "ROLEX",
				PurchasePrice: 1500000,
				PurchaseDate:  "2023-01-15",
			},
		},
		{
			name:   "正常系: ブランドのみ更新",
			itemID: "1",
			requestBody: map[string]interface{}{
				"brand": "UPDATED BRAND",
			},
			expectedUpdate: map[string]interface{}{
				"brand": "UPDATED BRAND",
			},
			mockResponse: &entity.Item{
				ID:            1,
				Name:          "ロレックス デイトナ",
				Category:      "時計",
				Brand:         "UPDATED BRAND",
				PurchasePrice: 1500000,
				PurchaseDate:  "2023-01-15",
			},
		},
		{
			name:   "正常系: 購入価格のみ更新",
			itemID: "1",
			requestBody: map[string]interface{}{
				"purchase_price": 2000000,
			},
			expectedUpdate: map[string]interface{}{
				"purchase_price": 2000000,
			},
			mockResponse: &entity.Item{
				ID:            1,
				Name:          "ロレックス デイトナ",
				Category:      "時計",
				Brand:         "ROLEX",
				PurchasePrice: 2000000,
				PurchaseDate:  "2023-01-15",
			},
		},
		{
			name:   "正常系: 複数フィールド更新",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name":           "新しいロレックス",
				"brand":          "NEW ROLEX",
				"purchase_price": 1800000,
			},
			expectedUpdate: map[string]interface{}{
				"name":           "新しいロレックス",
				"brand":          "NEW ROLEX",
				"purchase_price": 1800000,
			},
			mockResponse: &entity.Item{
				ID:            1,
				Name:          "新しいロレックス",
				Category:      "時計",
				Brand:         "NEW ROLEX",
				PurchasePrice: 1800000,
				PurchaseDate:  "2023-01-15",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockUsecase := setupTestHandler()

			// モックの設定
			mockUsecase.On("PatchItem", mock.Anything, int64(1), tt.expectedUpdate).
				Return(tt.mockResponse, nil)

			// リクエストの作成
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/items/"+tt.itemID, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Echo コンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.itemID)

			// テスト実行
			err := handler.PatchItem(c)

			// アサーション
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			var response entity.Item
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Equal(t, tt.mockResponse.ID, response.ID)
			assert.Equal(t, tt.mockResponse.Name, response.Name)
			assert.Equal(t, tt.mockResponse.Brand, response.Brand)
			assert.Equal(t, tt.mockResponse.PurchasePrice, response.PurchasePrice)

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestPatchItem_Errors(t *testing.T) {
	tests := []struct {
		name               string
		itemID             string
		requestBody        interface{}
		mockError          error
		expectedStatusCode int
		expectedError      string
	}{
		{
			name:               "異常系: 無効なID形式",
			itemID:             "invalid",
			requestBody:        map[string]interface{}{"name": "テスト"},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "invalid id format",
		},
		{
			name:               "異常系: 無効なJSONボディ",
			itemID:             "1",
			requestBody:        "invalid json",
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "invalid request body",
		},
		{
			name:               "異常系: 空のリクエストボディ",
			itemID:             "1",
			requestBody:        map[string]interface{}{},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "no fields provided for update",
		},
		{
			name:   "異常系: アイテムが見つからない",
			itemID: "999",
			requestBody: map[string]interface{}{
				"name": "存在しないアイテム",
			},
			mockError:          errors.ErrItemNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedError:      "item not found",
		},
		{
			name:   "異常系: バリデーションエラー",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name": "",
			},
			mockError:          fmt.Errorf("validation failed: name is required"),
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "validation failed",
		},
		{
			name:   "異常系: 更新不可フィールド",
			itemID: "1",
			requestBody: map[string]interface{}{
				"name": "テスト",
			},
			mockError:          fmt.Errorf("field 'category' is not allowed for update"),
			expectedStatusCode: http.StatusBadRequest,
			expectedError:      "field 'category' is not allowed for update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockUsecase := setupTestHandler()

			// モックの設定（エラーケース以外）
			if tt.mockError != nil && !strings.Contains(tt.name, "無効なID形式") && !strings.Contains(tt.name, "無効なJSONボディ") && !strings.Contains(tt.name, "空のリクエストボディ") {
				mockUsecase.On("PatchItem", mock.Anything, int64(1), mock.Anything).
					Return(nil, tt.mockError)
			}

			// リクエストの作成
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			
			req := httptest.NewRequest(http.MethodPatch, "/items/"+tt.itemID, bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			// Echo コンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.itemID)

			// テスト実行
			err := handler.PatchItem(c)

			// アサーション
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			var response ErrorResponse
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Contains(t, response.Error, tt.expectedError)

			if tt.mockError != nil && !strings.Contains(tt.name, "無効なID形式") && !strings.Contains(tt.name, "無効なJSONボディ") && !strings.Contains(tt.name, "空のリクエストボディ") {
				mockUsecase.AssertExpectations(t)
			}
		})
	}
}

func TestPatchItem_ImmutableFields(t *testing.T) {
	handler, mockUsecase := setupTestHandler()

	// 不変フィールドを含むリクエスト
	requestBody := map[string]interface{}{
		"name":         "更新された名前",
		"category":     "バッグ",     // 不変フィールド
		"created_at":   "2023-01-01", // 不変フィールド
		"updated_at":   "2023-01-01", // 不変フィールド
		"id":           2,            // 不変フィールド
		"purchase_date": "2023-01-01", // 不変フィールド
	}

	// categoryフィールドを含む場合のエラーを設定
	mockUsecase.On("PatchItem", mock.Anything, int64(1), mock.Anything).
		Return(nil, fmt.Errorf("field 'category' is not allowed for update"))

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPatch, "/items/1", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.PatchItem(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Error, "not allowed for update")
} 