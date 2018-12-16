package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anmolbabu/contact-book/cb_errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/anmolbabu/contact-book/dao"
	"github.com/anmolbabu/contact-book/dao/mocks"
	"github.com/anmolbabu/contact-book/e2e"
	"github.com/anmolbabu/contact-book/models"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    models.Contact
		wantStatus int
		mockRet    interface{}
	}{
		{
			name: "Case 1: Successful creation of contact",
			reqBody: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			wantStatus: http.StatusCreated,
			mockRet:    nil,
		},
		{
			name: "Case 2: failed creation request due to invalid emailid",
			reqBody: models.Contact{
				Name:    "def",
				EmailID: "def@",
			},
			wantStatus: http.StatusBadRequest,
			mockRet:    nil,
		},
		{
			name: "Case 3: failed createion request because contact already exists",
			reqBody: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			wantStatus: http.StatusConflict,
			mockRet:    cb_errors.DUPLICATE_CONTACT,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			w := httptest.NewRecorder()

			fkDB := new(mocks.DataAccess)
			dao.SetDAO(fkDB)

			r := gin.Default()
			payload, _ := json.Marshal(tt.reqBody)
			chi := Init(fkDB)
			fkDB.On("Add", mock.Anything).Return(tt.mockRet)
			req, _ := http.NewRequest("POST", "/contacts", bytes.NewReader(payload))
			req.Header.Add("Authorization", "Basic "+e2e.BasicAuth(e2e.UserName, e2e.Password))
			r.POST("/contacts", chi.Add)

			r.ServeHTTP(w, req)

			assert.Equal(tt.wantStatus, w.Code)
		})
	}
}
