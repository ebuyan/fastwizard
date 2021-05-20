package repository

import "wizard/pkg/db"

type Phones map[int]bool

func FindBlackListPhones(db *db.DB) (Phones, error) {
	rows, err := db.Query(`
	SELECT
		gbl.phone :: BIGINT 
	FROM
		global_black_lists gbl 
	WHERE
		gbl.rs_id = 0 
		AND ( gbl.user_id = $1 OR gbl.user_id = 0 ) 
		AND NOT EXISTS (
			SELECT 1 FROM global_white_list gwl WHERE gwl.rs_id = 0 AND ( gwl.user_id = $1 OR gwl.user_id = 0 ) AND gwl.phone = gbl.phone 
		)
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
		return nil, err
	}
	phones := Phones{}
	for rows.Next() {
		var phone int
		rows.Scan(&phone)
		phones[phone] = true
	}
	return phones, nil
}
