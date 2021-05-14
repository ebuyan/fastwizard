package repository

import "wizard/db"

type BlackListRepository struct{ *db.DB }

type Phones map[string]bool

func NewBlackListRepository(db *db.DB) BlackListRepository {
	return BlackListRepository{db}
}

func (r BlackListRepository) FindBlackListPhones() (phones Phones, err error) {
	rows, err := r.Query(`
	select phone
	from global_black_lists
	where (user_id = $1 and rs_id = 0) or  (user_id = 0 and rs_id = 0)
	and phone::text not in (
		select phone
		from global_white_list
		where (user_id = $1 and rs_id = 0) or  (user_id = 0 and rs_id = 0)
	)
	union
	select phone::text
	from phone_records
	left join address_book on address_book.id = phone_records.book_id
	where address_book.black_lst = 1 and address_book.user_id = $1
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
