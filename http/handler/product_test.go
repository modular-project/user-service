package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-service/mocks"
	"users-service/pkg"

	"github.com/labstack/echo"
	"github.com/modular-project/protobuffers/information/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		give     product.Product
		wantCode int
		wantID   uint
		wantErr  bool
	}{
		{
			give: product.Product{
				Name:        "Product OK",
				Url:         "img",
				Description: "nuevo producto",
				Price:       15.8,
			},
			wantCode: http.StatusCreated,
			wantID:   1,
			wantErr:  false,
		}, {
			give:     product.Product{Name: "Empty"},
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
	}
	assert := assert.New(t)
	ps := mocks.NewProductServicer(t)
	ps.On("Create", mock.Anything, mock.Anything).Return(uint64(1), nil).Once()
	ps.On("Create", mock.Anything, mock.Anything).Return(uint64(0), pkg.NewAppError("Fail at create", nil, http.StatusBadRequest)).Once()
	puc := NewProductUC(ps)
	for i, tt := range tests {

		data, err := json.Marshal(tt.give)
		if err != nil {
			t.Fatalf("Error at marshal: %s", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(data))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		ctx := e.NewContext(r, w)
		t.Log("CASE: ", i)
		gotErr := puc.Create(ctx)
		if !tt.wantErr {
			assert.Equal(tt.wantCode, w.Code)
			assert.NoError(gotErr)
		} else {
			//var err pkg.BadErr
			assert.Error(gotErr)
		}

	}
}

func TestGetProduct(t *testing.T) {
	tests := []struct {
		give        uint64
		wantCode    int
		wantPorudct product.Product
		wantErr     bool
	}{
		{
			give:     1,
			wantCode: http.StatusOK,
			wantPorudct: product.Product{
				Id:          1,
				Name:        "Product Ok",
				Description: "First product",
				Url:         "img",
				Price:       18.7,
			},
		}, {
			give:     2,
			wantCode: http.StatusBadRequest,
			wantErr:  true,
		},
	}
	assert := assert.New(t)
	ps := mocks.NewProductServicer(t)
	ps.On("Get", mock.Anything, uint64(1)).Return(tests[0].wantPorudct, nil)
	ps.On("Get", mock.Anything, uint64(2)).Return(product.Product{}, errors.New("not found"))
	h := NewProductUC(ps)
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/v1/product/:id")
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprint(tt.give))
		gotErr := h.Get(c)
		if tt.wantErr {
			assert.Error(gotErr)
		} else {
			assert.NoError(gotErr)
			assert.Equal(tt.wantCode, rec.Code)
			if tt.wantCode == http.StatusOK {
				p, err := json.Marshal(tt.wantPorudct)
				if assert.NoError(err) {
					assert.Equal(p, rec.Body.Bytes()[:len(rec.Body.Bytes())-1])
				}
			}

		}

	}
}

func TestGetInBatch(t *testing.T) {
	type IDs struct {
		IDs []uint64 `json:"ids"`
	}
	tests := []struct {
		give         IDs
		wantCode     int
		wantProducts []product.Product
	}{
		{
			give:     IDs{IDs: []uint64{1, 2}},
			wantCode: http.StatusOK,
			wantProducts: []product.Product{
				{
					Id:          1,
					Name:        "Product Ok",
					Description: "First product",
					Url:         "img",
					Price:       18.7,
				}, {
					Id:          2,
					Name:        "Product Ok 2",
					Description: "Second product",
					Url:         "img-2",
					Price:       28.14,
				},
			},
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		ps := mocks.NewProductServicer(t)
		//ps.GetInBatch()
		products := make([]*product.Product, len(tt.wantProducts))
		for i := range tt.wantProducts {
			products[i] = &tt.wantProducts[i]
		}
		for _, p := range products {
			t.Logf("%+v", *p)
		}
		ps.On("GetInBatch", mock.Anything, mock.Anything).Return(products, nil)
		ids, err := json.Marshal(tt.give)
		if !assert.NoError(err) {
			return
		}
		h := NewProductUC(ps)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(ids))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/v1/product/batch/")
		if assert.NoError(h.GetInBatch(c)) {
			assert.Equal(tt.wantCode, rec.Code)
			if tt.wantCode == http.StatusOK {
				var gotProducts []product.Product
				err := json.Unmarshal(rec.Body.Bytes(), &gotProducts)
				if assert.NoError(err) {
					assert.Equal(tt.wantProducts, gotProducts)
				}
			}
		}
	}
}
