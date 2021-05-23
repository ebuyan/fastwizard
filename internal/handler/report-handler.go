package handler

import (
	"encoding/json"
	"wizard/internal/app"
	"wizard/internal/phonestore"
	"wizard/internal/repository"
	"wizard/pkg/memory"
)

func GetReport(app app.App) (js []byte, err error) {
	phoneStorage, err := phonestore.NewPhoneStorage(app.Redis, app.Key)
	if err != nil {
		return
	}
	blackListPhones, err := repository.FindBlackListPhones(app.DB)
	if err != nil {
		return
	}
	report := phonestore.NewPhoneReport(phoneStorage, blackListPhones)
	if err != nil {
		return
	}
	js, _ = json.Marshal(report.GetReport())
	memory.Usage()
	return
}
