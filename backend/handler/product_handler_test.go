package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
)

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) CreateProduct(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductService) GetProductByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *mockProductService) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	args := m.Called(ctx, skip, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *mockProductService) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *mockProductService) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockProductService) GetProductsByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *mockProductService) UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(s *mockProductService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常系: 商品が作成される",
			requestBody: models.Product{
				Name:  "Test Product",
				Price: 100,
			},
			mockBehavior: func(s *mockProductService) {
				s.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"name":  "Test Product",
				"price": float64(100),
			},
		},
		{
			name:        "異常系: 無効なリクエストボディ",
			requestBody: "invalid",
			mockBehavior: func(s *mockProductService) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "無効なリクエストボディです",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockService := new(mockProductService)
			tt.mockBehavior(mockService)
			handler := NewProductHandler(mockService)

			// HTTPリクエストの準備
			e := echo.New()
			jsonData, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// ハンドラーの実行
			err := handler.CreateProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// エラーレスポンスの場合
			if tt.expectedStatus != http.StatusCreated {
				assert.Equal(t, tt.expectedBody, response)
				return
			}

			// 正常レスポンスの場合、期待する値が含まれているか確認
			for key, value := range tt.expectedBody {
				assert.Equal(t, value, response[key])
			}
		})
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	validID := primitive.NewObjectID()
	product := &models.Product{
		ID:    validID,
		Name:  "Test Product",
		Price: 100,
	}

	tests := []struct {
		name           string
		productID      string
		mockBehavior   func(s *mockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "正常系: 商品が取得される",
			productID: validID.Hex(),
			mockBehavior: func(s *mockProductService) {
				s.On("GetProductByID", mock.Anything, validID).Return(product, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   product,
		},
		{
			name:      "異常系: 無効なID",
			productID: "invalid",
			mockBehavior: func(s *mockProductService) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "無効なIDです",
			},
		},
		{
			name:      "異常系: 商品が見つからない",
			productID: primitive.NewObjectID().Hex(),
			mockBehavior: func(s *mockProductService) {
				s.On("GetProductByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "商品が見つかりません",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockService := new(mockProductService)
			tt.mockBehavior(mockService)
			handler := NewProductHandler(mockService)

			// HTTPリクエストの準備
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			// ハンドラーの実行
			err := handler.GetProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// レスポンスの内容を検証
			if tt.expectedStatus == http.StatusOK {
				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	products := []*models.Product{
		{
			ID:    primitive.NewObjectID(),
			Name:  "Product 1",
			Price: 100,
		},
		{
			ID:    primitive.NewObjectID(),
			Name:  "Product 2",
			Price: 200,
		},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockBehavior   func(s *mockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:        "正常系: 商品リストが取得される (ページネーションなし)",
			queryParams: map[string]string{},
			mockBehavior: func(s *mockProductService) {
				s.On("List", mock.Anything, int64(0), int64(10)).Return(products, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   products,
		},
		{
			name:        "異常系: 商品リストの取得に失敗",
			queryParams: map[string]string{},
			mockBehavior: func(s *mockProductService) {
				s.On("List", mock.Anything, int64(0), int64(10)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "商品リストの取得に失敗しました",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockService := new(mockProductService)
			tt.mockBehavior(mockService)
			handler := NewProductHandler(mockService)

			// HTTPリクエストの準備
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// クエリパラメータの設定
			q := make(map[string][]string)
			for k, v := range tt.queryParams {
				q[k] = []string{v}
			}
			req.URL.RawQuery = ""

			// ハンドラーの実行
			err := handler.ListProducts(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// レスポンスの内容を検証
			if tt.expectedStatus == http.StatusOK {
				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	validID := primitive.NewObjectID()
	updatedProduct := &models.Product{
		ID:    validID,
		Name:  "Updated Product",
		Price: 200,
	}

	tests := []struct {
		name           string
		productID      string
		requestBody    interface{}
		mockBehavior   func(s *mockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "正常系: 商品が更新される",
			productID: validID.Hex(),
			requestBody: models.Product{
				Name:  "Updated Product",
				Price: 200,
			},
			mockBehavior: func(s *mockProductService) {
				s.On("Update", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   updatedProduct,
		},
		{
			name:        "異常系: 無効なID",
			productID:   "invalid",
			requestBody: models.Product{},
			mockBehavior: func(s *mockProductService) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "無効な商品IDです",
			},
		},
		{
			name:        "異常系: 無効なリクエストボディ",
			productID:   validID.Hex(),
			requestBody: "invalid",
			mockBehavior: func(s *mockProductService) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "無効なリクエストボディです",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockService := new(mockProductService)
			tt.mockBehavior(mockService)
			handler := NewProductHandler(mockService)

			// HTTPリクエストの準備
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			// リクエストボディの準備
			jsonData, _ := json.Marshal(tt.requestBody)
			req.Body = io.NopCloser(bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			// ハンドラーの実行
			err := handler.UpdateProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// レスポンスの内容を検証
			if tt.expectedStatus == http.StatusOK {
				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			} else {
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	validID := primitive.NewObjectID()

	tests := []struct {
		name           string
		productID      string
		mockBehavior   func(s *mockProductService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "正常系: 商品が削除される",
			productID: validID.Hex(),
			mockBehavior: func(s *mockProductService) {
				s.On("Delete", mock.Anything, validID).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "商品を削除しました",
			},
		},
		{
			name:      "異常系: 無効なID",
			productID: "invalid",
			mockBehavior: func(s *mockProductService) {
				// モックは呼ばれないはず
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "無効な商品IDです",
			},
		},
		{
			name:      "異常系: 商品の削除に失敗",
			productID: validID.Hex(),
			mockBehavior: func(s *mockProductService) {
				s.On("Delete", mock.Anything, validID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "商品の削除に失敗しました",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの準備
			mockService := new(mockProductService)
			tt.mockBehavior(mockService)
			handler := NewProductHandler(mockService)

			// HTTPリクエストの準備
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.productID)

			// ハンドラーの実行
			err := handler.DeleteProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// レスポンスの内容を検証
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
