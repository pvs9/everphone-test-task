package pgx

import (
	"context"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/uptrace/bun"
)

type TagPGX struct {
	db *bun.DB
}

func NewTagPGX(db *bun.DB) *TagPGX {
	return &TagPGX{db: db}
}

func (r *TagPGX) AttachManyEmployees(employeeTags []christmas.EmployeeTag) ([]christmas.EmployeeTag, error) {
	_, err := r.db.NewInsert().
		Model(&employeeTags).
		Exec(context.TODO())

	if err != nil {
		return employeeTags, err
	}

	return employeeTags, nil
}

func (r *TagPGX) AttachManyGifts(giftTags []christmas.GiftTag) ([]christmas.GiftTag, error) {
	_, err := r.db.NewInsert().
		Model(&giftTags).
		Exec(context.TODO())

	if err != nil {
		return giftTags, err
	}

	return giftTags, nil
}

func (r *TagPGX) CreateByNames(names []string) ([]christmas.Tag, error) {
	var tags []christmas.Tag

	for _, tagName := range names {
		tags = append(tags, christmas.Tag{Name: tagName})
	}

	_, err := r.db.NewInsert().
		Model(&tags).
		Exec(context.TODO())

	if err != nil {
		return tags, err
	}

	return tags, nil
}

func (r *TagPGX) GetByNames(names []string, createMissing bool) ([]christmas.Tag, error) {
	var tags []christmas.Tag
	err := r.db.NewSelect().Model(&tags).Where("name IN (?)", bun.In(names)).Scan(context.TODO())

	if err != nil {
		return tags, err
	}

	if createMissing && (len(tags) != len(names)) {
		var missingTags []string

		for _, tagName := range names {
			found := false

			for i := range tags {
				if tags[i].Name == tagName {
					found = true
					break
				}
			}

			if !found {
				missingTags = append(missingTags, tagName)
			}
		}

		if len(missingTags) > 0 {
			newTags, err := r.CreateByNames(missingTags)

			if err != nil {
				return tags, err
			}

			tags = append(tags, newTags...)
		}
	}

	return tags, nil
}
