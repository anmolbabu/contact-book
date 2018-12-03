package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"net/http"

	"github.com/anmolbabu/contact-book/models"

	"time"
)

func TestCreateContact(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		contact models.Contact
	}{
		{
			name: "Case 1: Valid contact create request",
			contact: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
	}
	for _, tt := range tests {
		cmd, err := PrepTestSuite()
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		defer cmd.Process.Kill()

		requestByte, err := json.Marshal(tt.contact)
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		reqBody := bytes.NewReader(requestByte)

		res, err := HttpRequest(http.MethodPost, fmt.Sprintf("%s/contacts", BaseURL), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		resp := make(map[string]string)
		err = json.Unmarshal(res, &resp)
		if err != nil {
			t.Errorf("Failed to unmarshal %s", string(res))
			return
		}
		if status, ok := resp["Status"]; !ok || status != "Success" {
			t.Errorf("Failed to create resource %+v with response %+v", tt.contact, resp)
			return
		}

		var got models.Contact

		getRes, err := HttpRequest(http.MethodGet, fmt.Sprintf("%s/contacts/%s", BaseURL, tt.contact.EmailID), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		err = json.Unmarshal(getRes, &got)
		if err != nil {
			t.Errorf("failed to unmarshal %s", string(getRes))
			return
		}
		if !got.IsSame(tt.contact) {
			t.Errorf("expected %+v, received %+v", tt.contact, got)
			return
		}
	}
}

func TestUpdateContact(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		update  map[string]string
		contact models.Contact
		prep    func(models.Contact) error
		want    models.Contact
	}{
		{
			name: "Case 1: Valid contact update request",
			contact: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			update: map[string]string{"Name": "def"},
			want: models.Contact{
				Name:    "def",
				EmailID: "abc@email.com",
			},
			prep: func(contact models.Contact) error {
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				requestByte, err := json.Marshal(contact)
				if err != nil {
					return err
				}
				reqBody := bytes.NewReader(requestByte)
				res, err := HttpRequest(http.MethodPost, fmt.Sprintf("%s/contacts", BaseURL), headers, reqBody)
				if err != nil {
					return err
				}
				resp := make(map[string]string)
				err = json.Unmarshal(res, &resp)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		cmd, err := PrepTestSuite()
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		defer cmd.Process.Kill()
		requestByte, _ := json.Marshal(tt.update)
		reqBody := bytes.NewReader(requestByte)

		err = tt.prep(tt.contact)
		if err != nil {
			t.Errorf("Error in prep func: %+v", err)
			return
		}

		res, err := HttpRequest(http.MethodPut, fmt.Sprintf("%s/contacts/%s", BaseURL, tt.contact.EmailID), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		resp := make(map[string]string)
		err = json.Unmarshal(res, &resp)
		if err != nil {
			t.Errorf("Failed to unmarshal %s", string(res))
			return
		}
		if status, ok := resp["Status"]; !ok || status != "Success" {
			t.Errorf("Failed to create resource %+v with response %+v", tt.contact, resp)
			return
		}

		var got models.Contact

		getRes, err := HttpRequest(http.MethodGet, fmt.Sprintf("%s/contacts/%s", BaseURL, tt.contact.EmailID), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		err = json.Unmarshal(getRes, &got)
		if err != nil {
			t.Errorf("failed to unmarshal %s", string(getRes))
			return
		}
		if !got.IsSame(tt.want) {
			t.Errorf("expected %+v, received %+v", tt.contact, got)
			return
		}
	}
}

func TestListContact(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		prep    func([]models.Contact) error
		want    []models.Contact
	}{
		{
			name: "Case 1: Valid contact list request",
			want: []models.Contact{
				models.Contact{
					Name:    "abc",
					EmailID: "abc@email.com",
				},
				models.Contact{
					Name:    "def",
					EmailID: "def@email.com",
				},
			},
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			prep: func(contacts []models.Contact) error {
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				for _, contact := range contacts {
					requestByte, err := json.Marshal(contact)
					if err != nil {
						return err
					}
					reqBody := bytes.NewReader(requestByte)
					res, err := HttpRequest(http.MethodPost, fmt.Sprintf("%s/contacts", BaseURL), headers, reqBody)
					if err != nil {
						return err
					}
					resp := make(map[string]string)
					err = json.Unmarshal(res, &resp)
					if err != nil {
						return err
					}
					if status, ok := resp["Status"]; !ok || status != "Success" {
						return fmt.Errorf("Failed to create resource %+v with response %+v", contact, resp)
					}
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		cmd, err := PrepTestSuite()
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		defer cmd.Process.Kill()

		err = tt.prep(tt.want)
		if err != nil {
			t.Errorf("Error in prep func: %+v", err)
			return
		}

		var got []models.Contact
		requestByte, _ := json.Marshal(tt.want)
		reqBody := bytes.NewReader(requestByte)

		getRes, err := HttpRequest(http.MethodGet, fmt.Sprintf("%s/contacts", BaseURL), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		err = json.Unmarshal(getRes, &got)
		if err != nil {
			t.Errorf("failed to unmarshal %s", string(getRes))
			return
		}
		sort.Slice(got[:], func(i, j int) bool {
			return got[i].EmailID < got[j].EmailID
		})
		sort.Slice(tt.want[:], func(i, j int) bool {
			return tt.want[i].EmailID < tt.want[j].EmailID
		})
		if len(tt.want) != len(got) {
			t.Errorf("expected %+v, received %+v", tt.want, got)
			return
		}
		for ind := range got {
			if !got[ind].IsSame(tt.want[ind]) {
				t.Errorf("expected %+v, received %+v", tt.want, got)
				return
			}
		}
	}
}

func TestSearchContact(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		prep     func([]models.Contact) error
		want     models.Contact
		contacts []models.Contact
		wantName string
	}{
		{
			name: "Case 1: Valid contact search request",
			contacts: []models.Contact{
				models.Contact{
					Name:    "abc",
					EmailID: "abc@email.com",
				},
				models.Contact{
					Name:    "def",
					EmailID: "def@email.com",
				},
			},
			want: models.Contact{
				Name:    "def",
				EmailID: "def@email.com",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			wantName: "def",
			prep: func(contacts []models.Contact) error {
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				for _, contact := range contacts {
					requestByte, err := json.Marshal(contact)
					if err != nil {
						return err
					}
					reqBody := bytes.NewReader(requestByte)
					res, err := HttpRequest(http.MethodPost, fmt.Sprintf("%s/contacts", BaseURL), headers, reqBody)
					if err != nil {
						return err
					}
					resp := make(map[string]string)
					err = json.Unmarshal(res, &resp)
					if err != nil {
						return err
					}
					if status, ok := resp["Status"]; !ok || status != "Success" {
						return fmt.Errorf("Failed to create resource %+v with response %+v", contact, resp)
					}
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		cmd, err := PrepTestSuite()
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		defer cmd.Process.Kill()

		err = tt.prep(tt.contacts)
		if err != nil {
			t.Errorf("Error in prep func: %+v", err)
			return
		}

		var got []models.Contact
		requestByte, _ := json.Marshal(tt.want)
		reqBody := bytes.NewReader(requestByte)

		getRes, err := HttpRequest(http.MethodGet, fmt.Sprintf("%s/contacts?Name=%s", BaseURL, tt.wantName), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		err = json.Unmarshal(getRes, &got)
		if err != nil {
			t.Errorf("failed to unmarshal %s", string(getRes))
			return
		}
		if len(got) != 1 {
			t.Errorf("expected %+v, received %+v", tt.want, got)
			return
		}
		if !got[0].IsSame(tt.want) {
			t.Errorf("expected %+v, received %+v", tt.want, got)
			return
		}
	}
}

func TestDeleteContact(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		update  map[string]string
		contact models.Contact
		prep    func(models.Contact) error
		want    models.Contact
	}{
		{
			name: "Case 1: Valid contact delete request",
			contact: models.Contact{
				Name:    "abc",
				EmailID: "abc@email.com",
			},
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			prep: func(contact models.Contact) error {
				headers := map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				}
				requestByte, err := json.Marshal(contact)
				if err != nil {
					return err
				}
				reqBody := bytes.NewReader(requestByte)
				res, err := HttpRequest(http.MethodPost, fmt.Sprintf("%s/contacts", BaseURL), headers, reqBody)
				if err != nil {
					return err
				}
				resp := make(map[string]string)
				err = json.Unmarshal(res, &resp)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
	for _, tt := range tests {
		cmd, err := PrepTestSuite()
		if err != nil {
			t.Errorf("Error : %+v", err)
			return
		}
		defer cmd.Process.Kill()
		requestByte, _ := json.Marshal(tt.update)
		reqBody := bytes.NewReader(requestByte)

		err = tt.prep(tt.contact)
		if err != nil {
			t.Errorf("Error in prep func: %+v", err)
			return
		}

		res, err := HttpRequest(http.MethodDelete, fmt.Sprintf("%s/contacts/%s", BaseURL, tt.contact.EmailID), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		resp := make(map[string]string)
		err = json.Unmarshal(res, &resp)
		if err != nil {
			t.Errorf("Failed to unmarshal %s", string(res))
			return
		}
		if status, ok := resp["Status"]; !ok || status != "Success" {
			t.Errorf("Failed to create resource %+v with response %+v", tt.contact, resp)
			return
		}

		var got models.Contact

		getRes, err := HttpRequest(http.MethodGet, fmt.Sprintf("%s/contacts/%s", BaseURL, tt.contact.EmailID), tt.headers, reqBody)
		if err != nil {
			t.Errorf("Error %+v", err)
			return
		}
		err = json.Unmarshal(getRes, &got)
		if err != nil {
			t.Errorf("failed to unmarshal %s", string(getRes))
			return
		}
		if !got.IsSame(models.Contact{}) {
			t.Errorf("expected %+v, received %+v", tt.contact, got)
			return
		}
	}
}

func PrepTestSuite() (*exec.Cmd, error) {
	goPath := os.Getenv("GOPATH")
	cmdStr := filepath.Join(goPath, "bin", "contact-book")

	cmd := exec.Command(cmdStr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		cmd.Wait()
	}()

	time.Sleep(5 * time.Second)
	return cmd, nil
}
