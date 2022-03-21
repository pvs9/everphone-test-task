package christmas

type EmployeeTag struct {
	EmployeeID int64     `bun:",pk"`
	Employee   *Employee `bun:"rel:belongs-to,join:employee_id=id"`
	TagID      int64     `bun:",pk"`
	Tag        *Tag      `bun:"rel:belongs-to,join:tag_id=id"`
}
