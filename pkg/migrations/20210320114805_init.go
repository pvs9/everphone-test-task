package migrations

import (
	"context"
	christmas "github.com/pvs9/everphone-test-task"
	log "github.com/sirupsen/logrus"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		m2mModels := []interface{}{
			(*christmas.EmployeeTag)(nil),
			(*christmas.GiftTag)(nil),
		}

		models := []interface{}{
			(*christmas.Employee)(nil),
			(*christmas.Gift)(nil),
			(*christmas.Tag)(nil),
			(*christmas.EmployeeTag)(nil),
			(*christmas.GiftTag)(nil),
		}

		for _, m2mModel := range m2mModels {
			db.RegisterModel(m2mModel)
		}

		for _, model := range models {
			_, err := db.NewCreateTable().Model(model).IfNotExists().WithForeignKeys().Exec(ctx)

			if err != nil {
				log.Errorf("[MIGRATIONS] Error occured while migrating: %s", err.Error())
			}
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		m2mModels := []interface{}{
			(*christmas.EmployeeTag)(nil),
			(*christmas.GiftTag)(nil),
		}

		models := []interface{}{
			(*christmas.Employee)(nil),
			(*christmas.Gift)(nil),
			(*christmas.Tag)(nil),
		}

		for _, m2mModel := range m2mModels {
			db.RegisterModel(m2mModel)
		}

		for _, model := range models {
			_, err := db.NewDropTable().Model(model).IfExists().Exec(ctx)

			if err != nil {
				log.Errorf("[MIGRATIONS] Error occured while rolling back: %s", err.Error())
			}
		}

		return nil
	})
}
