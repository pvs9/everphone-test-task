package christmas

type GiftTag struct {
	GiftID int64 `bun:",pk"`
	Gift   *Gift `bun:"rel:belongs-to,join:gift_id=id"`
	TagID  int64 `bun:",pk"`
	Tag    *Tag  `bun:"rel:belongs-to,join:tag_id=id"`
}
