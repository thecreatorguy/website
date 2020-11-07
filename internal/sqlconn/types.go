package sqlconn

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
    if !nt.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(nt.Time)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Valid = true
	return json.Unmarshal(data, &nt.Time)
}