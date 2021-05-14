package repository

import (
	"wizard/psql"
)

type BlackListRepository struct{}

type Phones map[string]bool

func (r BlackListRepository) FindBlackListPhones() (phones Phones, err error) {
	rows, err := psql.DB.Query(`
		SELECT
			phone :: BIGINT 
		FROM
			global_black_lists 
		WHERE
			( user_id = $1 AND rs_id = 0 ) 
			OR ( user_id = 0 AND rs_id = 0 ) 
			AND phone NOT IN ( SELECT phone FROM global_white_list WHERE ( user_id = $1 AND rs_id = 0 ) OR ( user_id = 0 AND rs_id = 0 ) )
		UNION
		SELECT
			phone 
		FROM
			phone_records
			LEFT JOIN address_book ON address_book.ID = phone_records.book_id 
		WHERE
			address_book.black_lst = 1 
			AND address_book.user_id = $1
	`, 2607)
	if err != nil {
		return
	}
	phones = Phones{}
	for rows.Next() {
		var phone string
		err = rows.Scan(&phone)
		phones[phone] = true
	}
	return
}
