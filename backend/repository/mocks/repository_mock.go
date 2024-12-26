// Code generated by MockGen. DO NOT EDIT.
// Source: repository/interfaces.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	models "smart-store-admin/backend/models"
	time "time"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockProductRepositoryMockRecorder) Create(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductRepository)(nil).Create), ctx, product)
}

// Delete mocks base method.
func (m *MockProductRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProductRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProductRepository)(nil).Delete), ctx, id)
}

// GetByCategory mocks base method.
func (m *MockProductRepository) GetByCategory(ctx context.Context, category string) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCategory", ctx, category)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCategory indicates an expected call of GetByCategory.
func (mr *MockProductRepositoryMockRecorder) GetByCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCategory", reflect.TypeOf((*MockProductRepository)(nil).GetByCategory), ctx, category)
}

// GetByID mocks base method.
func (m *MockProductRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockProductRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockProductRepository)(nil).GetByID), ctx, id)
}

// GetLowStock mocks base method.
func (m *MockProductRepository) GetLowStock(ctx context.Context) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLowStock", ctx)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLowStock indicates an expected call of GetLowStock.
func (mr *MockProductRepositoryMockRecorder) GetLowStock(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLowStock", reflect.TypeOf((*MockProductRepository)(nil).GetLowStock), ctx)
}

// List mocks base method.
func (m *MockProductRepository) List(ctx context.Context, skip, limit int64) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, skip, limit)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockProductRepositoryMockRecorder) List(ctx, skip, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockProductRepository)(nil).List), ctx, skip, limit)
}

// Update mocks base method.
func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockProductRepositoryMockRecorder) Update(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProductRepository)(nil).Update), ctx, product)
}

// MockDeliveryRepository is a mock of DeliveryRepository interface.
type MockDeliveryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDeliveryRepositoryMockRecorder
}

// MockDeliveryRepositoryMockRecorder is the mock recorder for MockDeliveryRepository.
type MockDeliveryRepositoryMockRecorder struct {
	mock *MockDeliveryRepository
}

// NewMockDeliveryRepository creates a new mock instance.
func NewMockDeliveryRepository(ctrl *gomock.Controller) *MockDeliveryRepository {
	mock := &MockDeliveryRepository{ctrl: ctrl}
	mock.recorder = &MockDeliveryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeliveryRepository) EXPECT() *MockDeliveryRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDeliveryRepository) Create(ctx context.Context, delivery *models.Delivery) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, delivery)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDeliveryRepositoryMockRecorder) Create(ctx, delivery interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDeliveryRepository)(nil).Create), ctx, delivery)
}

// GetActiveDeliveries mocks base method.
func (m *MockDeliveryRepository) GetActiveDeliveries(ctx context.Context) ([]*models.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveDeliveries", ctx)
	ret0, _ := ret[0].([]*models.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveDeliveries indicates an expected call of GetActiveDeliveries.
func (mr *MockDeliveryRepositoryMockRecorder) GetActiveDeliveries(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveDeliveries", reflect.TypeOf((*MockDeliveryRepository)(nil).GetActiveDeliveries), ctx)
}

// GetByID mocks base method.
func (m *MockDeliveryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDeliveryRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDeliveryRepository)(nil).GetByID), ctx, id)
}

// GetDeliveriesByRobot mocks base method.
func (m *MockDeliveryRepository) GetDeliveriesByRobot(ctx context.Context, robotID string) ([]*models.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeliveriesByRobot", ctx, robotID)
	ret0, _ := ret[0].([]*models.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeliveriesByRobot indicates an expected call of GetDeliveriesByRobot.
func (mr *MockDeliveryRepositoryMockRecorder) GetDeliveriesByRobot(ctx, robotID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeliveriesByRobot", reflect.TypeOf((*MockDeliveryRepository)(nil).GetDeliveriesByRobot), ctx, robotID)
}

// List mocks base method.
func (m *MockDeliveryRepository) List(ctx context.Context, skip, limit int64) ([]*models.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, skip, limit)
	ret0, _ := ret[0].([]*models.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDeliveryRepositoryMockRecorder) List(ctx, skip, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDeliveryRepository)(nil).List), ctx, skip, limit)
}

// UpdateLocation mocks base method.
func (m *MockDeliveryRepository) UpdateLocation(ctx context.Context, id primitive.ObjectID, location models.Location) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLocation", ctx, id, location)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLocation indicates an expected call of UpdateLocation.
func (mr *MockDeliveryRepositoryMockRecorder) UpdateLocation(ctx, id, location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLocation", reflect.TypeOf((*MockDeliveryRepository)(nil).UpdateLocation), ctx, id, location)
}

// UpdateStatus mocks base method.
func (m *MockDeliveryRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.DeliveryStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockDeliveryRepositoryMockRecorder) UpdateStatus(ctx, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockDeliveryRepository)(nil).UpdateStatus), ctx, id, status)
}

// MockSaleRepository is a mock of SaleRepository interface.
type MockSaleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSaleRepositoryMockRecorder
}

// MockSaleRepositoryMockRecorder is the mock recorder for MockSaleRepository.
type MockSaleRepositoryMockRecorder struct {
	mock *MockSaleRepository
}

// NewMockSaleRepository creates a new mock instance.
func NewMockSaleRepository(ctrl *gomock.Controller) *MockSaleRepository {
	mock := &MockSaleRepository{ctrl: ctrl}
	mock.recorder = &MockSaleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSaleRepository) EXPECT() *MockSaleRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSaleRepository) Create(ctx context.Context, sale *models.Sale) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, sale)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSaleRepositoryMockRecorder) Create(ctx, sale interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSaleRepository)(nil).Create), ctx, sale)
}

// GetByID mocks base method.
func (m *MockSaleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*models.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockSaleRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockSaleRepository)(nil).GetByID), ctx, id)
}

// GetDailySales mocks base method.
func (m *MockSaleRepository) GetDailySales(ctx context.Context, date time.Time) ([]*models.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDailySales", ctx, date)
	ret0, _ := ret[0].([]*models.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDailySales indicates an expected call of GetDailySales.
func (mr *MockSaleRepositoryMockRecorder) GetDailySales(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDailySales", reflect.TypeOf((*MockSaleRepository)(nil).GetDailySales), ctx, date)
}

// GetEnvironmentalImpactAnalytics mocks base method.
func (m *MockSaleRepository) GetEnvironmentalImpactAnalytics(ctx context.Context, start, end time.Time) (*models.EnvironmentalImpact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnvironmentalImpactAnalytics", ctx, start, end)
	ret0, _ := ret[0].(*models.EnvironmentalImpact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEnvironmentalImpactAnalytics indicates an expected call of GetEnvironmentalImpactAnalytics.
func (mr *MockSaleRepositoryMockRecorder) GetEnvironmentalImpactAnalytics(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnvironmentalImpactAnalytics", reflect.TypeOf((*MockSaleRepository)(nil).GetEnvironmentalImpactAnalytics), ctx, start, end)
}

// GetSalesByCategory mocks base method.
func (m *MockSaleRepository) GetSalesByCategory(ctx context.Context, start, end time.Time) (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSalesByCategory", ctx, start, end)
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSalesByCategory indicates an expected call of GetSalesByCategory.
func (mr *MockSaleRepositoryMockRecorder) GetSalesByCategory(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSalesByCategory", reflect.TypeOf((*MockSaleRepository)(nil).GetSalesByCategory), ctx, start, end)
}

// GetSalesByDateRange mocks base method.
func (m *MockSaleRepository) GetSalesByDateRange(ctx context.Context, start, end time.Time) ([]*models.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSalesByDateRange", ctx, start, end)
	ret0, _ := ret[0].([]*models.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSalesByDateRange indicates an expected call of GetSalesByDateRange.
func (mr *MockSaleRepositoryMockRecorder) GetSalesByDateRange(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSalesByDateRange", reflect.TypeOf((*MockSaleRepository)(nil).GetSalesByDateRange), ctx, start, end)
}

// GetSalesByTimeOfDay mocks base method.
func (m *MockSaleRepository) GetSalesByTimeOfDay(ctx context.Context, timeOfDay string) ([]*models.Sale, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSalesByTimeOfDay", ctx, timeOfDay)
	ret0, _ := ret[0].([]*models.Sale)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSalesByTimeOfDay indicates an expected call of GetSalesByTimeOfDay.
func (mr *MockSaleRepositoryMockRecorder) GetSalesByTimeOfDay(ctx, timeOfDay interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSalesByTimeOfDay", reflect.TypeOf((*MockSaleRepository)(nil).GetSalesByTimeOfDay), ctx, timeOfDay)
}

// GetTotalSalesAmount mocks base method.
func (m *MockSaleRepository) GetTotalSalesAmount(ctx context.Context, start, end time.Time) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalSalesAmount", ctx, start, end)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalSalesAmount indicates an expected call of GetTotalSalesAmount.
func (mr *MockSaleRepositoryMockRecorder) GetTotalSalesAmount(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalSalesAmount", reflect.TypeOf((*MockSaleRepository)(nil).GetTotalSalesAmount), ctx, start, end)
}

// MockStoreOperationRepository is a mock of StoreOperationRepository interface.
type MockStoreOperationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStoreOperationRepositoryMockRecorder
}

// MockStoreOperationRepositoryMockRecorder is the mock recorder for MockStoreOperationRepository.
type MockStoreOperationRepositoryMockRecorder struct {
	mock *MockStoreOperationRepository
}

// NewMockStoreOperationRepository creates a new mock instance.
func NewMockStoreOperationRepository(ctrl *gomock.Controller) *MockStoreOperationRepository {
	mock := &MockStoreOperationRepository{ctrl: ctrl}
	mock.recorder = &MockStoreOperationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreOperationRepository) EXPECT() *MockStoreOperationRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStoreOperationRepository) Create(ctx context.Context, op *models.StoreOperation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, op)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockStoreOperationRepositoryMockRecorder) Create(ctx, op interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStoreOperationRepository)(nil).Create), ctx, op)
}

// GetAverageEnergyUsage mocks base method.
func (m *MockStoreOperationRepository) GetAverageEnergyUsage(ctx context.Context, start, end time.Time) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAverageEnergyUsage", ctx, start, end)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAverageEnergyUsage indicates an expected call of GetAverageEnergyUsage.
func (mr *MockStoreOperationRepositoryMockRecorder) GetAverageEnergyUsage(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAverageEnergyUsage", reflect.TypeOf((*MockStoreOperationRepository)(nil).GetAverageEnergyUsage), ctx, start, end)
}

// GetByTimeRange mocks base method.
func (m *MockStoreOperationRepository) GetByTimeRange(ctx context.Context, start, end time.Time) ([]*models.StoreOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByTimeRange", ctx, start, end)
	ret0, _ := ret[0].([]*models.StoreOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByTimeRange indicates an expected call of GetByTimeRange.
func (mr *MockStoreOperationRepositoryMockRecorder) GetByTimeRange(ctx, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByTimeRange", reflect.TypeOf((*MockStoreOperationRepository)(nil).GetByTimeRange), ctx, start, end)
}

// GetLatest mocks base method.
func (m *MockStoreOperationRepository) GetLatest(ctx context.Context) (*models.StoreOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatest", ctx)
	ret0, _ := ret[0].(*models.StoreOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatest indicates an expected call of GetLatest.
func (mr *MockStoreOperationRepositoryMockRecorder) GetLatest(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatest", reflect.TypeOf((*MockStoreOperationRepository)(nil).GetLatest), ctx)
}

// UpdateCheckoutStatus mocks base method.
func (m *MockStoreOperationRepository) UpdateCheckoutStatus(ctx context.Context, opID primitive.ObjectID, status models.CheckoutStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCheckoutStatus", ctx, opID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCheckoutStatus indicates an expected call of UpdateCheckoutStatus.
func (mr *MockStoreOperationRepositoryMockRecorder) UpdateCheckoutStatus(ctx, opID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCheckoutStatus", reflect.TypeOf((*MockStoreOperationRepository)(nil).UpdateCheckoutStatus), ctx, opID, status)
}

// UpdateShelfStatus mocks base method.
func (m *MockStoreOperationRepository) UpdateShelfStatus(ctx context.Context, opID primitive.ObjectID, status models.ShelfStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShelfStatus", ctx, opID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShelfStatus indicates an expected call of UpdateShelfStatus.
func (mr *MockStoreOperationRepositoryMockRecorder) UpdateShelfStatus(ctx, opID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShelfStatus", reflect.TypeOf((*MockStoreOperationRepository)(nil).UpdateShelfStatus), ctx, opID, status)
}
