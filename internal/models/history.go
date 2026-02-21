package models

import "time"

// History represents a change log for warehouse items
type History struct {
	ID        int       `json:"id" db:"id"`
	ItemID    int       `json:"item_id" db:"item_id"`
	Action    string    `json:"action" db:"action"` // CREATE, UPDATE, DELETE
	OldName   *string   `json:"old_name,omitempty" db:"old_name"`
	NewName   *string   `json:"new_name,omitempty" db:"new_name"`
	OldQty    *int      `json:"old_quantity,omitempty" db:"old_quantity"`
	NewQty    *int      `json:"new_quantity,omitempty" db:"new_quantity"`
	UserID    string    `json:"user_id" db:"user_id"`
	UserName  string    `json:"user_name" db:"user_name"`
	ChangedAt time.Time `json:"changed_at" db:"changed_at"`
}
