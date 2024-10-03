package products

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/kobietka/product-service/pkg/api/v1"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandleGetProduct(t *testing.T) {
	tests := []struct {
		Name         string
		MockValue    v1.Product
		MockError    error
		ExpectedCode int
		ExpectedBody *v1.Product
	}{
		{
			Name:         "returns product correctly",
			MockValue:    v1.Product{Ean: "12345678", Name: "Product name"},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: &v1.Product{Ean: "12345678", Name: "Product name"},
		},
		{
			Name:         "store returns not found error",
			MockValue:    v1.Product{},
			MockError:    v1.ErrorDataNotFound,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: nil,
		},
		{
			Name:         "store returns unknown error",
			MockValue:    v1.Product{},
			MockError:    errors.New("error"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On("GetProduct", mock.Anything, mock.Anything).Return(test.MockValue, test.MockError)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleGetProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj v1.Product
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, *test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandleSearchProduct(t *testing.T) {
	tests := []struct {
		Name         string
		Query        string
		Limit        string
		MockValue    []v1.Product
		MockError    error
		ExpectedCode int
		ExpectedBody []v1.Product
	}{
		{
			Name:         "returns product correctly",
			Query:        "prod",
			Limit:        "5",
			MockValue:    []v1.Product{{Ean: "12345678", Name: "Product name"}},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []v1.Product{{Ean: "12345678", Name: "Product name"}},
		},
		{
			Name:         "store returns an unknown error",
			Query:        "prod",
			Limit:        "5",
			MockValue:    []v1.Product{},
			MockError:    errors.New("error"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
		{
			Name:         "empty query",
			Query:        "",
			Limit:        "5",
			MockValue:    []v1.Product{},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
		{
			Name:         "limit negative",
			Query:        "prod",
			Limit:        "-5",
			MockValue:    []v1.Product{},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
		{
			Name:         "limit zero",
			Query:        "prod",
			Limit:        "0",
			MockValue:    []v1.Product{},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
		{
			Name:         "limit higher than predefined limit",
			Query:        "prod",
			Limit:        fmt.Sprintf("%d", searchLimit+10),
			MockValue:    []v1.Product{{Ean: "12345678", Name: "Product name"}},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []v1.Product{{Ean: "12345678", Name: "Product name"}},
		},
		{
			Name:         "limit higher than int8 can store",
			Query:        "prod",
			Limit:        "200",
			MockValue:    []v1.Product{},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On("SearchProducts", mock.Anything, mock.Anything, mock.Anything).Return(test.MockValue, test.MockError).Run(func(args mock.Arguments) {
				limit := args[2].(int8)
				testLimitInt, err := strconv.Atoi(test.Limit)
				assert.NoError(t, err)

				assert.Equal(t, min(int8(testLimitInt), searchLimit), limit)
			})

			request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products?query=%s&limit=%s", test.Query, test.Limit), nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleSearchProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj []v1.Product
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandlePostProduct(t *testing.T) {
	correctProduct := v1.Product{
		Ean:  "12345678",
		Name: "Product name",
		Packaging: v1.Quantity{
			Value: 12,
			Unit:  "g",
		},
		Nutrition: v1.Nutrition{
			Per: v1.Quantity{
				Value: 1,
				Unit:  "g",
			},
			Kcal: 123,
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
			},
		},
	}

	tests := []struct {
		Name         string
		RequestBody  *v1.Product
		MockError    error
		ExpectedCode int
		ExpectedBody *v1.ErrorResponse
	}{
		{
			Name:         "no request body",
			RequestBody:  nil,
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: &v1.ErrorResponse{Code: ErrorProductEanMissing.Error()},
		},
		{
			Name:         "creates product",
			RequestBody:  &correctProduct,
			MockError:    nil,
			ExpectedCode: http.StatusCreated,
			ExpectedBody: nil,
		},
		{
			Name: "validation returns an error",
			RequestBody: &v1.Product{
				Ean: "123",
			},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: &v1.ErrorResponse{Code: ErrorProductEanInvalid.Error()},
		},
		{
			Name:         "store returns invalid data error",
			RequestBody:  &correctProduct,
			MockError:    v1.ErrorInvalidData,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
		{
			Name:         "store returns an unknown error",
			RequestBody:  &correctProduct,
			MockError:    errors.New("error"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On("CreateProduct", mock.Anything, mock.Anything).Return(test.MockError)

			var request *http.Request
			if test.RequestBody != nil {
				jsonBytes, err := json.Marshal(*test.RequestBody)
				assert.NoError(t, err)
				request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBytes))
				request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				request = httptest.NewRequest(http.MethodPost, "/", nil)
			}

			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handlePostProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj v1.ErrorResponse
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, *test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandlePutProduct(t *testing.T) {
	correctProduct := v1.Product{
		Ean:  "12345678",
		Name: "Product name",
		Packaging: v1.Quantity{
			Value: 12,
			Unit:  "g",
		},
		Nutrition: v1.Nutrition{
			Per: v1.Quantity{
				Value: 1,
				Unit:  "g",
			},
			Kcal: 123,
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 1,
						Unit:  "g",
					},
				},
			},
		},
	}

	tests := []struct {
		Name         string
		RequestBody  *v1.Product
		MockError    error
		ExpectedCode int
		ExpectedBody *v1.ErrorResponse
	}{
		{
			Name:         "no request body",
			RequestBody:  nil,
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: &v1.ErrorResponse{Code: ErrorProductEanMissing.Error()},
		},
		{
			Name:         "updates product",
			RequestBody:  &correctProduct,
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: nil,
		},
		{
			Name: "validation returns an error",
			RequestBody: &v1.Product{
				Ean: "123",
			},
			MockError:    nil,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: &v1.ErrorResponse{Code: ErrorProductEanInvalid.Error()},
		},
		{
			Name:         "store returns invalid data error",
			RequestBody:  &correctProduct,
			MockError:    v1.ErrorInvalidData,
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: nil,
		},
		{
			Name:         "store returns product does not exist error",
			RequestBody:  &correctProduct,
			MockError:    v1.ErrorProductDoesNotExist,
			ExpectedCode: http.StatusNotFound,
			ExpectedBody: nil,
		},
		{
			Name:         "store returns an unknown error",
			RequestBody:  &correctProduct,
			MockError:    errors.New("error"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On("UpdateProduct", mock.Anything, mock.Anything).Return(test.MockError)

			var request *http.Request
			if test.RequestBody != nil {
				jsonBytes, err := json.Marshal(*test.RequestBody)
				assert.NoError(t, err)
				request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBytes))
				request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				request = httptest.NewRequest(http.MethodPost, "/", nil)
			}

			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handlePutProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj v1.ErrorResponse
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, *test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandleDeleteProduct(t *testing.T) {
	tests := []struct {
		Name         string
		MockError    error
		ExpectedCode int
	}{
		{
			Name:         "deletes product",
			MockError:    nil,
			ExpectedCode: http.StatusNoContent,
		},
		{
			Name:         "store returns product does not exist error",
			MockError:    v1.ErrorProductDoesNotExist,
			ExpectedCode: http.StatusNotFound,
		},
		{
			Name:         "store returns unknown error",
			MockError:    errors.New("error"),
			ExpectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On("DeleteProduct", mock.Anything, mock.Anything).Return(test.MockError)

			request := httptest.NewRequest(http.MethodDelete, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleDeleteProduct(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
		})
	}
}

type MockStore struct {
	mock.Mock
}

func (s *MockStore) GetProduct(ctx context.Context, ean string) (v1.Product, error) {
	args := s.Called(ctx, ean)
	return args.Get(0).(v1.Product), args.Error(1)
}

func (s *MockStore) SearchProducts(ctx context.Context, query string, limit int8) ([]v1.Product, error) {
	args := s.Called(ctx, query, limit)
	return args.Get(0).([]v1.Product), args.Error(1)
}

func (s *MockStore) CreateProduct(ctx context.Context, product v1.Product) error {
	args := s.Called(ctx, product)
	return args.Error(0)
}

func (s *MockStore) UpdateProduct(ctx context.Context, product v1.Product) error {
	args := s.Called(ctx, product)
	return args.Error(0)
}

func (s *MockStore) DeleteProduct(ctx context.Context, ean string) error {
	args := s.Called(ctx, ean)
	return args.Error(0)
}
