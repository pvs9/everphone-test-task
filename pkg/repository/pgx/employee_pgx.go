package pgx

import (
	"context"
	"database/sql"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/uptrace/bun"
)

type EmployeePGX struct {
	db *bun.DB
}

func NewEmployeePGX(db *bun.DB) *EmployeePGX {
	return &EmployeePGX{db: db}
}

func (r *EmployeePGX) CreateByNames(names []string) ([]christmas.Employee, error) {
	var employees []christmas.Employee

	for _, employeeName := range names {
		employees = append(employees, christmas.Employee{Name: employeeName})
	}

	_, err := r.db.NewInsert().
		Model(&employees).
		Exec(context.TODO())

	if err != nil {
		return employees, err
	}

	return employees, nil
}

func (r *EmployeePGX) DetachTagsFromManyByIds(ids []int64) error {
	_, err := r.db.NewDelete().
		Model((*christmas.EmployeeTag)(nil)).
		Where("employee_id IN (?)", bun.In(ids)).
		Exec(context.TODO())

	if err != nil {
		return err
	}

	return nil
}

func (r *EmployeePGX) GetAll() ([]christmas.Employee, error) {
	var employees []christmas.Employee
	err := r.db.NewSelect().Model(&employees).Relation("Tags").Scan(context.TODO())

	if err != nil {
		return employees, err
	}

	return employees, nil
}

func (r *EmployeePGX) GetById(id int64) (*christmas.Employee, error) {
	var employee christmas.Employee
	err := r.db.NewSelect().
		Model(&employee).
		Where("id = ?", id).
		Relation("Tags").
		Scan(context.TODO())

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &employee, nil
}

func (r *EmployeePGX) GetByNames(names []string, createMissing bool) ([]christmas.Employee, error) {
	var employees []christmas.Employee
	err := r.db.NewSelect().Model(&employees).Where("name IN (?)", bun.In(names)).Scan(context.TODO())

	if err != nil {
		return employees, err
	}

	if createMissing && (len(employees) != len(names)) {
		var missingEmployees []string

		for _, employeeName := range names {
			found := false

			for i := range employees {
				if employees[i].Name == employeeName {
					found = true
					break
				}
			}

			if !found {
				missingEmployees = append(missingEmployees, employeeName)
			}
		}

		if len(missingEmployees) > 0 {
			newEmployees, err := r.CreateByNames(missingEmployees)

			if err != nil {
				return employees, err
			}

			employees = append(employees, newEmployees...)
		}
	}

	return employees, nil
}

func (r *EmployeePGX) Update(employee christmas.Employee) (christmas.Employee, error) {
	_, err := r.db.NewUpdate().
		Model(&employee).
		WherePK().
		Exec(context.TODO())

	if err != nil {
		return employee, err
	}

	return employee, nil
}

func (r *EmployeePGX) SyncTags(employee christmas.Employee, tags []christmas.Tag) (christmas.Employee, error) {
	_, err := r.db.NewDelete().
		Model((*christmas.EmployeeTag)(nil)).
		Where("employee_id = ?", employee.ID).
		Exec(context.TODO())
	var employeeTags []christmas.EmployeeTag

	for _, tag := range tags {
		employeeTags = append(employeeTags, christmas.EmployeeTag{EmployeeID: employee.ID, TagID: tag.ID})
	}

	_, err = r.db.NewInsert().
		Model(&employeeTags).
		Exec(context.TODO())

	if err != nil {
		return employee, err
	}

	return employee, nil
}
