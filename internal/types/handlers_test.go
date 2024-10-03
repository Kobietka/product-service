package types

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetUnits(t *testing.T) {
	tests := []struct {
		Name         string
		MockMethod   string
		MockValue    []string
		MockError    error
		ExpectedCode int
		ExpectedBody []string
	}{
		{
			Name:         "returns units correctly",
			MockMethod:   "GetUnits",
			MockValue:    []string{"g", "kg"},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{"g", "kg"},
		},
		{
			Name:         "returns empty units correctly",
			MockMethod:   "GetUnits",
			MockValue:    []string{},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{},
		},
		{
			Name:         "units store returns an error",
			MockMethod:   "GetUnits",
			MockValue:    nil,
			MockError:    errors.New("err"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On(test.MockMethod, mock.Anything).Return(test.MockValue, test.MockError)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleGetUnits(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj []string
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandleGetNutrients(t *testing.T) {
	tests := []struct {
		Name         string
		MockMethod   string
		MockValue    []string
		MockError    error
		ExpectedCode int
		ExpectedBody []string
	}{
		{
			Name:         "returns nutrient types correctly",
			MockMethod:   "GetNutrientTypes",
			MockValue:    []string{"PROTEIN", "FAT"},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{"PROTEIN", "FAT"},
		},
		{
			Name:         "returns empty nutrient types correctly",
			MockMethod:   "GetNutrientTypes",
			MockValue:    []string{},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{},
		},
		{
			Name:         "nutrient type store returns an error",
			MockMethod:   "GetNutrientTypes",
			MockValue:    nil,
			MockError:    errors.New("err"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On(test.MockMethod, mock.Anything).Return(test.MockValue, test.MockError)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleGetNutrients(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj []string
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandleGetVitamins(t *testing.T) {
	tests := []struct {
		Name         string
		MockMethod   string
		MockValue    []string
		MockError    error
		ExpectedCode int
		ExpectedBody []string
	}{
		{
			Name:         "returns vitamin types correctly",
			MockMethod:   "GetVitaminTypes",
			MockValue:    []string{"PROTEIN", "FAT"},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{"PROTEIN", "FAT"},
		},
		{
			Name:         "returns empty vitamin types correctly",
			MockMethod:   "GetVitaminTypes",
			MockValue:    []string{},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{},
		},
		{
			Name:         "vitamin type store returns an error",
			MockMethod:   "GetVitaminTypes",
			MockValue:    nil,
			MockError:    errors.New("err"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On(test.MockMethod, mock.Anything).Return(test.MockValue, test.MockError)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleGetVitamins(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj []string
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

func TestHandleGetMinerals(t *testing.T) {
	tests := []struct {
		Name         string
		MockMethod   string
		MockValue    []string
		MockError    error
		ExpectedCode int
		ExpectedBody []string
	}{
		{
			Name:         "returns mineral types correctly",
			MockMethod:   "GetMineralTypes",
			MockValue:    []string{"IRON", "CALCIUM"},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{"IRON", "CALCIUM"},
		},
		{
			Name:         "returns empty mineral types correctly",
			MockMethod:   "GetMineralTypes",
			MockValue:    []string{},
			MockError:    nil,
			ExpectedCode: http.StatusOK,
			ExpectedBody: []string{},
		},
		{
			Name:         "mineral type store returns an error",
			MockMethod:   "GetMineralTypes",
			MockValue:    nil,
			MockError:    errors.New("err"),
			ExpectedCode: http.StatusInternalServerError,
			ExpectedBody: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			store := new(MockStore)
			server := NewServer(store)

			store.On(test.MockMethod, mock.Anything).Return(test.MockValue, test.MockError)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			response := httptest.NewRecorder()
			c := echo.New().NewContext(request, response)

			err := server.handleGetMinerals(c)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedCode, response.Code)
			if test.ExpectedBody != nil {
				assert.Equal(t, echo.MIMEApplicationJSON, response.Header().Get(echo.HeaderContentType))
				var obj []string
				err = json.NewDecoder(response.Body).Decode(&obj)
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedBody, obj)
			} else {
				assert.Equal(t, 0, len(response.Body.Bytes()))
			}
		})
	}
}

type MockStore struct {
	mock.Mock
}

func (m *MockStore) GetUnits(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStore) GetNutrientTypes(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStore) GetVitaminTypes(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStore) GetMineralTypes(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}
