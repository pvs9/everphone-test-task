package christmas

type Tag struct {
	ID   int64  `bun:",pk,autoincrement" json:"id,omitempty"`
	Name string `bun:",unique,nullzero,notnull" json:"name,omitempty"`
}
