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
)

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		emailid    string
		wantStatus int
		mockRet    interface{}
	}{
		{
			name:       "Case 1: Successful deletion of contact",
			emailid:    "abc@email.com",
			wantStatus: http.StatusOK,
			mockRet:    nil,
		},
		{
			name:       "Case 2: failed creation request due to invalid emailid",
			emailid:    "abc",
			wantStatus: http.StatusBadRequest,
			mockRet:    nil,
		},
		{
			name:       "Case 3: failed createion request because contact already exists",
			emailid:    "def@email.com",
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
			chi := Init(fkDB)
			fkDB.On("Delete", tt.emailid).Return(tt.mockRet)
			reqUrl := fmt.Sprintf("/contacts/%s", tt.emailid)

			req, err := http.NewRequest("DELETE", reqUrl, nil)
			if err != nil {
				t.Errorf("Invalid http request %+v. Error: %+v", tt, err)
				return
			}

			req.Header.Add("Authorization", "Basic "+e2e.BasicAuth(e2e.UserName, e2e.Password))
			r.DELETE("/contacts/:emailid", chi.Delete)

			r.ServeHTTP(w, req)

			assert.Equal(tt.wantStatus, w.Code)
		})
	}
}
