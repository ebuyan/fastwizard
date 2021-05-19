package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"wizard/internal/phonestore"
	"wizard/internal/repository"
)

type Handler struct{}

func (h Handler) report(w http.ResponseWriter, r *http.Request) (err error) {
	req := ReportRequest{}
	json.NewDecoder(r.Body).Decode(&req)
	if len(req.Key) == 0 {
		err = errors.New("Empty key")
		return
	}
	phoneStorage, err := phonestore.NewPhoneStorage(req.Key)
	if err != nil {
		return
	}
	blackListPhones, err := repository.FindBlackListPhones()
	if err != nil {
		return
	}
	report := phonestore.NewPhoneReport(phoneStorage, blackListPhones)
	if err != nil {
		return
	}
	js, _ := json.Marshal(report.GetReport())
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	report = nil
	return
}

func (h Handler) test(w http.ResponseWriter, r *http.Request) (err error) {
	phones, _ := repository.FindPhoneRecords()
	for key, phone := range phones {
		if key == 100 {
			fmt.Println(phone)
		}
	}
	phones = nil
	return
}

type ReportRequest struct {
	Key              string `json:"dispatch_key"`
	RemoveDuplicates bool   `json:"remove_duplicates"`
	ResetStorage     bool   `json:"reset_storage"`
}
