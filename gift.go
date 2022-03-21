package christmas

import "time"

type Gift struct {
	ID        int64     `bun:",pk,autoincrement" json:"id"`
	Name      string    `bun:",unique,nullzero,notnull" json:"name"`
	Tags      []Tag     `bun:"m2m:gift_tags,join:Gift=Tag" json:"categories,omitempty"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updatedAt"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deletedAt"`
}
