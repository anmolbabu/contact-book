package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anmolbabu/contact-book/cb_errors"
	"github.com/anmolbabu/contact-book/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/anmolbabu/contact-book/dao"
	"github.com/anmolbabu/contact-book/dao/mocks"
	"github.com/anmolbabu/contact-book/e2e"
	"github.com/anmolbabu/contact-book/models"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    models.Contact
		wantStatus int
		emailId    string
		mockRet    interface{}
	}{
		{
			name: "Case 1: Successful updation of contact",
			reqBody: models.Contact{
				Name: "pqr",
			},
			wantStatus: http.StatusOK,
			emailId:    "abc@email.com",
			mockRet:    nil,
		},
		{
			name: "Case 2: failed updation request due to invalid emailid",
			reqBody: models.Contact{
				Name:    utils.INVALID_STRING,
				EmailID: "ghi@",
			},
			emailId:    "def@",
			wantStatus: http.StatusBadRequest,
			mockRet:    nil,
		},
		{
			name: "Case 3: failed updation request because contact does not exist",
			reqBody: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			wantStatus: http.StatusNotFound,
			mockRet:    cb_errors.CONTACT_NOT_FOUND,
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
			fkDB.On("Update", tt.emailId, tt.reqBody.Name, tt.reqBody.EmailID).Return(tt.mockRet)

			reqURL := fmt.Sprintf("/contacts/%s", tt.emailId)
			req, _ := http.NewRequest("PUT", reqURL, bytes.NewReader(payload))
			req.Header.Add("Authorization", "Basic "+e2e.BasicAuth(e2e.UserName, e2e.Password))
			r.PUT("/contacts/:emailid", chi.Update)

			r.ServeHTTP(w, req)

			assert.Equal(tt.wantStatus, w.Code)
		})
	}
}
