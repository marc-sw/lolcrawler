package store

import (
	"strings"
)

func (store *Store) AddOrIgnoreMatchIds(matchIds []string) error {
	builder := strings.Builder{}
	values := make([]interface{}, len(matchIds))

	builder.WriteString("INSERT OR IGNORE INTO match_ids (match_id) VALUES ")
	for i, matchId := range matchIds {
		if i > 0 {
			builder.WriteRune(',')
		}
		builder.WriteString("(?) ")

		values[i] = any(matchId)
	}
	_, err := store.db.Exec(builder.String(), values...)
	return err
}

func (store *Store) GetMatchIdAtRow(row int) (string, error) {
	var matchId string
	err := store.db.QueryRow("SELECT match_id FROM match_ids LIMIT 1 OFFSET ?", row).Scan(&matchId)
	return matchId, err
}

func (store *Store) DeleteMatchId(matchId string) error {
	_, err := store.db.Exec("DELETE FROM match_ids WHERE match_id = ?", matchId)
	return err
}

func (store *Store) GetMatchIdsCount() (int, error) {
	var count int
	err := store.db.QueryRow("SELECT COUNT() FROM match_ids").Scan(&count)
	return count, err
}
