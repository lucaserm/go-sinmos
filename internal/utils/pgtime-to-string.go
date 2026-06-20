package utils

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func PgTimeToString(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}

	// microseconds -> segundos do dia
	seconds := t.Microseconds / 1_000_000

	h := seconds / 3600
	m := (seconds % 3600) / 60

	return fmt.Sprintf("%02d:%02d", h, m)
}
