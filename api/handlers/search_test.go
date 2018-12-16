package handlers

import (
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

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		queryStr       string
		queryContact   models.Contact
		queryPageNo    int
		queryPageLimit int
		wantStatus     int
		wantErr        error
		mockRet        []models.Contact
	}{
		{
			name:       "Case 1: Successful fetch of contact filtered by name and emailid",
			queryStr:   "?name=abc&emailid=abc@email.com",
			wantStatus: http.StatusOK,
			queryContact: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			queryPageNo:    -1,
			queryPageLimit: -1,
			mockRet: []models.Contact{
				models.Contact{
					Name:    "abc",
					EmailID: "abc@email.com",
				},
			},
			wantErr: nil,
		},
		{
			name:     "Case 2: Successful fetch of contacts filtered by name and emailid in page 1",
			queryStr: "?name=def&emailid=def@email.com&page=1&pagelimit=3",
			mockRet: []models.Contact{
				models.Contact{
					Name:    "def",
					EmailID: "def@email.com",
				},
			},
			queryContact: models.Contact{
				Name:    "def",
				EmailID: "def@email.com",
			},
			queryPageNo:    1,
			queryPageLimit: 3,
			wantStatus:     http.StatusOK,
			wantErr:        nil,
		},

		{
			name:       "Case 3: failed fetch request of 0 contacts fetched by name in page 2",
			queryStr:   "?name=abc&page=2&pagelimit=3",
			wantStatus: http.StatusNotFound,
			queryContact: models.Contact{
				Name:    "abc",
				EmailID: utils.INVALID_STRING,
			},
			queryPageNo:    2,
			queryPageLimit: 3,
			mockRet:        nil,
			wantErr:        cb_errors.CONTACT_NOT_FOUND,
		},
		{
			name:       "Case 4: invalid pageno and pagelimit",
			queryStr:   "?name=abc&page=0&pagelimit=0",
			wantStatus: http.StatusBadRequest,
			queryContact: models.Contact{
				Name: "abc",
			},
			queryPageNo:    0,
			queryPageLimit: 0,
			mockRet:        nil,
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
			fkDB.On("GetAll", &tt.queryContact, tt.queryPageNo, tt.queryPageLimit).Return(tt.mockRet, tt.wantErr)

			reqURL := fmt.Sprintf("/contacts%s", tt.queryStr)
			req, _ := http.NewRequest("GET", reqURL, nil)
			req.Header.Add("Authorization", "Basic "+e2e.BasicAuth(e2e.UserName, e2e.Password))
			r.GET("/contacts", chi.GetAll)

			r.ServeHTTP(w, req)

			assert.Equal(tt.wantStatus, w.Code)
		})
	}
}
