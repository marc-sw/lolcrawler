package store

const (
	CountAccountRow = "account_row"
	CountMatchIdRow = "match_id_row"
)

func (store *Store) GetCount(name string) (int, error) {
	var count int
	err := store.db.QueryRow("SELECT value FROM counts WHERE name = ?", name).Scan(&count)
	return count, err
}

func (store *Store) SetCount(name string, value int) error {
	_, err := store.db.Exec("UPDATE counts SET value = ? WHERE name = ?", value, name)
	return err
}
