package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anmolbabu/contact-book/cb_errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/anmolbabu/contact-book/dao"
	"github.com/anmolbabu/contact-book/dao/mocks"
	"github.com/anmolbabu/contact-book/e2e"
	"github.com/anmolbabu/contact-book/models"
)

func TestFetch(t *testing.T) {
	tests := []struct {
		name       string
		emailId    string
		wantStatus int
		wantErr    error
		mockRet    models.Contact
	}{
		{
			name:       "Case 1: Successful fetch of contact",
			wantStatus: http.StatusOK,
			emailId:    "abc@email.com",
			mockRet: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			wantErr: nil,
		},
		{
			name:       "Case 2: Contact does not exist",
			emailId:    "pqr@email.com",
			mockRet:    *(models.GetDefaultContact()),
			wantStatus: http.StatusNotFound,
			wantErr:    cb_errors.CONTACT_NOT_FOUND,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			w := httptest.NewRecorder()

			fkDB := new(mocks.DataAccess)
			dao.SetDAO(fkDB)

			r := gin.Default()
			chi := Init(fkDB)
			fkDB.On("Get", tt.emailId).Return(tt.mockRet, tt.wantErr)

			reqURL := fmt.Sprintf("/contacts/%s", tt.emailId)
			req, _ := http.NewRequest("GET", reqURL, nil)
			req.Header.Add("Authorization", "Basic "+e2e.BasicAuth(e2e.UserName, e2e.Password))
			r.GET("/contacts/:emailid", chi.Get)

			r.ServeHTTP(w, req)

			assert.Equal(tt.wantStatus, w.Code)
		})
	}
}
