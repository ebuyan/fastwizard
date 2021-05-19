package repository

import (
	"wizard/pkg/db"
)

type PhoneRecords []PhoneRecord
type PhoneRecord struct {
	Phone string
	Field string
}

func FindPhoneRecords() (phones PhoneRecords, err error) {
	rows, err := db.Conn.Query(`
	SELECT
		phone_records.phone,
		field1 || ',' ||	field2 || ',' ||	field3 || ',' ||	field4 as field
	FROM
		phone_records 
	WHERE
		book_id = $1
	`, 10754)
	if err != nil {
		return
	}
	phones = PhoneRecords{}
	for rows.Next() {
		var phone string
		var field string
		err = rows.Scan(&phone, &field)
		phones = append(phones, PhoneRecord{phone, field})
	}
	return
}
