package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// uuidBytes converts uuid.UUID → [16]byte for pgx encoding.
func uuidBytes(id uuid.UUID) [16]byte { return [16]byte(id) }

// nullableUUIDParam wraps *uuid.UUID into pgtype.UUID for nullable encoding.
func nullableUUIDParam(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: [16]byte(*id), Valid: true}
}

// pgUUIDPtr converts a scanned pgtype.UUID into *uuid.UUID.
func pgUUIDPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}

// toTextArray converts []string → pgtype.Array[pgtype.Text] for pgx encoding.
func toTextArray(ss []string) pgtype.Array[pgtype.Text] {
	elems := make([]pgtype.Text, len(ss))
	for i, s := range ss {
		elems[i] = pgtype.Text{String: s, Valid: true}
	}
	dims := []pgtype.ArrayDimension{{Length: int32(len(ss)), LowerBound: 1}}
	return pgtype.Array[pgtype.Text]{Elements: elems, Dims: dims, Valid: true}
}

// fromTextArray converts a scanned pgtype.Array[pgtype.Text] → []string.
func fromTextArray(a pgtype.Array[pgtype.Text]) []string {
	out := make([]string, 0, len(a.Elements))
	for _, e := range a.Elements {
		if e.Valid {
			out = append(out, e.String)
		}
	}
	return out
}

// numericPtr converts a scanned pgtype.Numeric → *float64 (nil when NULL).
func numericPtr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return nil
	}
	v := f.Float64
	return &v
}

// nullableTime converts *time.Time → pgtype.Timestamptz for encoding.
func nullableTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

// pgTimestamptzPtr converts a scanned pgtype.Timestamptz → *time.Time.
func pgTimestamptzPtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}

// marshalJSON marshals v → json.RawMessage; returns "{}" on error.
func marshalJSON(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage("{}")
	}
	return b
}
