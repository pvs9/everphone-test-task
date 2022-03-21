package christmas

import "time"

type Employee struct {
	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Name      string    `bun:",unique,nullzero,notnull" json:"name"`
	Tags      []Tag     `bun:"m2m:employee_tags,join:Employee=Tag" json:"interests,omitempty"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deletedAt"`
}
