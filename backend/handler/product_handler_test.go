package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"smart-store-admin/backend/models"
	mock_service "smart-store-admin/backend/service/mocks"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックサービスを作成
	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	// テストケースの定義
	testCases := []struct {
		name           string
		input          models.Product
		mockBehavior   func(product *models.Product)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系: 商品が作成される",
			input: models.Product{
				Name:  "Test Product",
				Price: 100,
			},
			mockBehavior: func(product *models.Product) {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), product).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"ID":"","Name":"Test Product","Price":100,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:  "異常系: 無効なリクエストボディ",
			input: models.Product{},
			mockBehavior: func(product *models.Product) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"無効なリクエストボディです"}`,
		},
		{
			name: "異常系: 商品の作成に失敗",
			input: models.Product{
				Name:  "Test Product",
				Price: 100,
			},
			mockBehavior: func(product *models.Product) {
				mockService.EXPECT().
					CreateProduct(gomock.Any(), product).
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"商品の作成に失敗しました"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエストの作成
			e := echo.New()
			reqBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックのセットアップ
			tc.mockBehavior(&tc.input)

			// ハンドラーの実行
			err := handler.CreateProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックサービスを作成
	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	// テストケースの定義
	testCases := []struct {
		name           string
		productID      string
		mockBehavior   func(id primitive.ObjectID)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "正常系: 商品が取得される",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			mockBehavior: func(id primitive.ObjectID) {
				mockService.EXPECT().
					GetProductByID(gomock.Any(), id).
					Return(&models.Product{
						ID:    id,
						Name:  "Test Product",
						Price: 100,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ID":"60a7b0b0f0f0f0f0f0f0f0f0","Name":"Test Product","Price":100,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "異常系: 無効なID",
			productID: "invalid-id",
			mockBehavior: func(id primitive.ObjectID) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"無効なIDです"}`,
		},
		{
			name:      "異常系: 商品が見つからない",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			mockBehavior: func(id primitive.ObjectID) {
				mockService.EXPECT().
					GetProductByID(gomock.Any(), id).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"商品が見つかりません"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエストの作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products/"+tc.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.productID)

			// モックのセットアップ
			id, _ := primitive.ObjectIDFromHex(tc.productID)
			tc.mockBehavior(id)

			// ��ンドラーの実行
			err := handler.GetProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックサービスを作成
	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	// テストケースの定義
	testCases := []struct {
		name           string
		queryParams    string
		mockBehavior   func(skip int64, limit int64)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系: 商品リストが取得される (ページネーションなし)",
			queryParams: "",
			mockBehavior: func(skip int64, limit int64) {
				mockService.EXPECT().
					List(gomock.Any(), skip, limit).
					Return([]*models.Product{
						{
							ID:    primitive.NewObjectID(),
							Name:  "Test Product 1",
							Price: 100,
						},
						{
							ID:    primitive.NewObjectID(),
							Name:  "Test Product 2",
							Price: 200,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":"000000000000000000000000","Name":"Test Product 1","Price":100,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"},{"ID":"000000000000000000000000","Name":"Test Product 2","Price":200,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:        "正常系: 商品リストが取得される (ページネーションあり)",
			queryParams: "?page=2&limit=5",
			mockBehavior: func(skip int64, limit int64) {
				mockService.EXPECT().
					List(gomock.Any(), skip, limit).
					Return([]*models.Product{
						{
							ID:    primitive.NewObjectID(),
							Name:  "Test Product 3",
							Price: 300,
						},
						{
							ID:    primitive.NewObjectID(),
							Name:  "Test Product 4",
							Price: 400,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":"000000000000000000000000","Name":"Test Product 3","Price":300,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"},{"ID":"000000000000000000000000","Name":"Test Product 4","Price":400,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name:        "異常系: 商品リストの取得に失敗",
			queryParams: "",
			mockBehavior: func(skip int64, limit int64) {
				mockService.EXPECT().
					List(gomock.Any(), skip, limit).
					Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"商品リストの取得に失敗しました"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエストの作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products"+tc.queryParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックのセットアップ
			page := 1
			limit := 10
			if pageStr := c.QueryParam("page"); pageStr != "" {
				if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
					page = p
				}
			}
			if limitStr := c.QueryParam("limit"); limitStr != "" {
				if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
					limit = l
				}
			}
			skip := int64((page - 1) * limit)
			tc.mockBehavior(skip, int64(limit))

			// ハンドラーの実行
			err := handler.ListProducts(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatus, rec.Code)
			// IDが毎回変わるので、一旦無視する
			var expectedBody []map[string]interface{}
			json.Unmarshal([]byte(tc.expectedBody), &expectedBody)
			var actualBody []map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualBody)
			for i := range expectedBody {
				expectedBody[i]["ID"] = actualBody[i]["ID"]
			}
			expectedBodyBytes, _ := json.Marshal(expectedBody)
			assert.JSONEq(t, string(expectedBodyBytes), rec.Body.String())
		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックサービスを作成
	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	// テストケースの定義
	testCases := []struct {
		name           string
		productID      string
		input          models.Product
		mockBehavior   func(id primitive.ObjectID, product *models.Product)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "正常系: 商品が更新される",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			input: models.Product{
				Name:  "Updated Product",
				Price: 200,
			},
			mockBehavior: func(id primitive.ObjectID, product *models.Product) {
				product.ID = id
				mockService.EXPECT().
					Update(gomock.Any(), product).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ID":"60a7b0b0f0f0f0f0f0f0f0f0","Name":"Updated Product","Price":200,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:      "異常系: 無効なID",
			productID: "invalid-id",
			input: models.Product{
				Name:  "Updated Product",
				Price: 200,
			},
			mockBehavior: func(id primitive.ObjectID, product *models.Product) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"無効なIDです"}`,
		},
		{
			name:      "異常系: 無効なリクエストボディ",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			input:     models.Product{},
			mockBehavior: func(id primitive.ObjectID, product *models.Product) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"無効なリクエストボディです"}`,
		},
		{
			name:      "異常系: 商品の更新に失敗",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			input: models.Product{
				Name:  "Updated Product",
				Price: 200,
			},
			mockBehavior: func(id primitive.ObjectID, product *models.Product) {
				product.ID = id
				mockService.EXPECT().
					Update(gomock.Any(), product).
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"商品の更新に失敗しました"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエストの作成
			e := echo.New()
			reqBody, _ := json.Marshal(tc.input)
			req := httptest.NewRequest(http.MethodPut, "/products/"+tc.productID, bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.productID)

			// モックのセットアップ
			id, _ := primitive.ObjectIDFromHex(tc.productID)
			tc.mockBehavior(id, &tc.input)

			// ハンドラーの実行
			err := handler.UpdateProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// モックサービスを作成
	mockService := mock_service.NewMockProductServiceInterface(ctrl)
	handler := NewProductHandler(mockService)

	// テストケースの定義
	testCases := []struct {
		name           string
		productID      string
		mockBehavior   func(id primitive.ObjectID)
		expectedStatus int
	}{
		{
			name:      "正常系: 商品が削除される",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			mockBehavior: func(id primitive.ObjectID) {
				mockService.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "異常系: 無効なID",
			productID: "invalid-id",
			mockBehavior: func(id primitive.ObjectID) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "異常系: 商品の削除に失敗",
			productID: "60a7b0b0f0f0f0f0f0f0f0f0",
			mockBehavior: func(id primitive.ObjectID) {
				mockService.EXPECT().
					Delete(gomock.Any(), id).
					Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// リクエストの作成
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/products/"+tc.productID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.productID)

			// モックのセットアップ
			id, _ := primitive.ObjectIDFromHex(tc.productID)
			tc.mockBehavior(id)

			// ハンドラーの実行
			err := handler.DeleteProduct(c)
			assert.NoError(t, err)

			// レスポンスの検証
			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
