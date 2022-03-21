package pgx

import (
	"context"
	"database/sql"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/uptrace/bun"
)

type GiftPGX struct {
	db *bun.DB
}

func NewGiftPGX(db *bun.DB) *GiftPGX {
	return &GiftPGX{db: db}
}

func (r *GiftPGX) CreateByNames(names []string) ([]christmas.Gift, error) {
	var gifts []christmas.Gift

	for _, employeeName := range names {
		gifts = append(gifts, christmas.Gift{Name: employeeName})
	}

	_, err := r.db.NewInsert().
		Model(&gifts).
		Exec(context.TODO())

	if err != nil {
		return gifts, err
	}

	return gifts, nil
}

func (r *GiftPGX) DetachTagsFromManyByIds(ids []int64) error {
	_, err := r.db.NewDelete().
		Model((*christmas.GiftTag)(nil)).
		Where("gift_id IN (?)", bun.In(ids)).
		Exec(context.TODO())

	if err != nil {
		return err
	}

	return nil
}

func (r *GiftPGX) GetAll() ([]christmas.Gift, error) {
	var gifts []christmas.Gift
	err := r.db.NewSelect().Model(&gifts).Relation("Tags").Scan(context.TODO())

	if err != nil {
		return gifts, err
	}

	return gifts, nil
}

func (r *GiftPGX) GetById(id int64) (*christmas.Gift, error) {
	var gift christmas.Gift
	err := r.db.NewSelect().
		Model(&gift).
		Where("id = ?", id).
		Relation("Tags").
		Scan(context.TODO())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &gift, nil
}

func (r *GiftPGX) GetByNames(names []string, createMissing bool) ([]christmas.Gift, error) {
	var gifts []christmas.Gift
	err := r.db.NewSelect().Model(&gifts).Where("name IN (?)", bun.In(names)).Scan(context.TODO())

	if err != nil {
		return gifts, err
	}

	if createMissing && (len(gifts) != len(names)) {
		var missingGifts []string

		for _, giftName := range names {
			found := false

			for i := range gifts {
				if gifts[i].Name == giftName {
					found = true
					break
				}
			}

			if !found {
				missingGifts = append(missingGifts, giftName)
			}
		}

		if len(missingGifts) > 0 {
			newEmployees, err := r.CreateByNames(missingGifts)

			if err != nil {
				return gifts, err
			}

			gifts = append(gifts, newEmployees...)
		}
	}

	return gifts, nil
}

func (r *GiftPGX) Update(gift christmas.Gift) (christmas.Gift, error) {
	_, err := r.db.NewUpdate().
		Model(&gift).
		WherePK().
		Exec(context.TODO())

	if err != nil {
		return gift, err
	}

	return gift, nil
}

func (r *GiftPGX) SyncTags(gift christmas.Gift, tags []christmas.Tag) (christmas.Gift, error) {
	_, err := r.db.NewDelete().
		Model((*christmas.GiftTag)(nil)).
		Where("gift_id = ?", gift.ID).
		Exec(context.TODO())
	var giftTags []christmas.GiftTag

	for _, tag := range tags {
		giftTags = append(giftTags, christmas.GiftTag{GiftID: gift.ID, TagID: tag.ID})
	}

	_, err = r.db.NewInsert().
		Model(&giftTags).
		Exec(context.TODO())

	if err != nil {
		return gift, err
	}

	return gift, nil
}
