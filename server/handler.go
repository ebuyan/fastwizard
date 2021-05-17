package server

import (
	"encoding/json"
	"net/http"
	phonestore "wizard/phone_store"
	"wizard/repository"
)

type Handler struct{}

func (h Handler) report(w http.ResponseWriter, r *http.Request) {
	req := ReportRequest{}
	json.NewDecoder(r.Body).Decode(&req)
	if len(req.Key) == 0 {
		http.Error(w, "Empty key", http.StatusInternalServerError)
		return
	}
	phoneStorage, err := phonestore.NewPhoneStorage(req.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blackListPhones, err := repository.BlackListRepository{}.FindBlackListPhones()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	report := phonestore.NewPhoneReport(phoneStorage, blackListPhones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal(report.GetReport())
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type ReportRequest struct {
	Key              string `json:"dispatch_key"`
	RemoveDuplicates bool   `json:"remove_duplicates"`
	ResetStorage     bool   `json:"reset_storage"`
}
