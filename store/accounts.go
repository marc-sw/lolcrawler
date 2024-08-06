package store

import (
	"github.com/KnutZuidema/golio/riot/account"
)

func (store *Store) AddOrIgnoreAccount(account account.Account) error {
	_, err := store.db.Exec("INSERT OR IGNORE INTO accounts (puuid, game_name, tag_line) VALUES (?,?,?)",
		account.Puuid,
		account.GameName,
		account.TagLine)
	return err
}

func (store *Store) GetAccount(puuid string) (account.Account, error) {
	acc := account.Account{}
	err := store.db.QueryRow("SELECT puuid, game_name, tag_line FROM accounts WHERE puuid = ?", puuid).
		Scan(&acc.Puuid, &acc.GameName, &acc.TagLine)
	return acc, err
}

func (store *Store) GetAccountAtRow(row int) (account.Account, error) {
	acc := account.Account{}
	err := store.db.QueryRow("SELECT puuid, game_name, tag_line FROM accounts LIMIT 1 OFFSET ?", row).
		Scan(&acc.Puuid, &acc.GameName, &acc.TagLine)
	return acc, err
}

func (store *Store) GetAccountsCount() (int, error) {
	var count int
	err := store.db.QueryRow("SELECT COUNT() FROM accounts").Scan(&count)
	return count, err
}

func (store *Store) GetEmptyAccountsPuuid() ([]string, error) {
	rows, err := store.db.Query("SELECT puuid FROM accounts WHERE game_name = '' OR tag_line = ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var puuids []string
	var puuid string
	for rows.Next() {
		err = rows.Scan(&puuid)
		if err != nil {
			return nil, err
		}
		puuids = append(puuids, puuid)
	}
	return puuids, rows.Err()
}

func (store *Store) UpdateAccount(account account.Account) error {
	_, err := store.db.Exec("UPDATE accounts SET game_name = ?, tag_line = ? WHERE puuid = ?",
		account.GameName, account.TagLine, account.Puuid)
	return err
}

func (store *Store) DeleteAccount(puuid string) error {
	_, err := store.db.Exec("DELETE FROM accounts WHERE puuid = ?", puuid)
	return err
}
