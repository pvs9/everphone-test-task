package service

import (
	"encoding/json"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/queue"
	"github.com/pvs9/everphone-test-task/pkg/repository"
	"github.com/pvs9/everphone-test-task/pkg/request"
	"io/ioutil"
)

type EmployeeService struct {
	publisher     queue.Publisher
	repository    repository.Employee
	tagRepository repository.Tag
}

func NewEmployeeService(publisher queue.Publisher, repository repository.Employee, tagRepository repository.Tag) *EmployeeService {
	return &EmployeeService{publisher: publisher, repository: repository, tagRepository: tagRepository}
}

func (s *EmployeeService) GetAll() ([]christmas.Employee, error) {
	employees, err := s.repository.GetAll()

	if err != nil {
		return employees, err
	}

	return employees, nil
}

func (s *EmployeeService) GetById(id int64) (*christmas.Employee, error) {
	employee, err := s.repository.GetById(id)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) Update(employee christmas.Employee, employeeData request.EmployeeUpdateRequest) (christmas.Employee, error) {
	tags, err := s.tagRepository.GetByNames(employeeData.Tags, true)

	if err != nil {
		return employee, err
	}

	employee.Name = employeeData.Name
	employee.Tags = tags

	_, err = s.repository.Update(employee)

	if err != nil {
		return employee, err
	}

	employee, err = s.repository.SyncTags(employee, tags)

	return employee, nil
}

func (s *EmployeeService) UploadDataset(fileName string) (*string, error) {
	message := DatasetMessage{Type: "employee", Filename: fileName}
	messageBody, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	messageId, err := s.publisher.Publish(string(messageBody))

	if err != nil {
		return nil, err
	}

	return messageId, nil
}

func (s *EmployeeService) ProcessDataset(fileName string) error {
	dataset, err := ioutil.ReadFile("/go/src/everphone-test-task.io/pkg/storage/" + fileName)

	if err != nil {
		return err
	}

	var data []request.EmployeeDatasetEntity
	err = json.Unmarshal(dataset, &data)
	var employeeNames []string
	var tagNames []string
	employeeTagsMap := make(map[string][]string)

	for index := range data {
		employeeNames = append(employeeNames, data[index].Name)
		employeeTagsMap[data[index].Name] = data[index].Tags
		tagNames = append(tagNames, data[index].Tags...)
	}

	var uniqueTagNames []string
	present := make(map[string]bool)

	for _, tagName := range tagNames {

		if _, ok := present[tagName]; !ok {
			present[tagName] = true
			uniqueTagNames = append(uniqueTagNames, tagName)
		}
	}

	employees, err := s.repository.GetByNames(employeeNames, true)

	if err != nil {
		return err
	}

	tags, err := s.tagRepository.GetByNames(uniqueTagNames, true)

	if err != nil {
		return err
	}

	var employeeTagModels []christmas.EmployeeTag
	var employeeIds []int64

	for _, employee := range employees {
		employeeTagNames, ok := employeeTagsMap[employee.Name]
		employeeIds = append(employeeIds, employee.ID)

		if ok && employeeTagNames != nil && len(employeeTagNames) > 0 {
			for i := range employeeTagNames {
				for j := range tags {
					if tags[j].Name == employeeTagNames[i] {
						employeeTagModels = append(employeeTagModels, christmas.EmployeeTag{EmployeeID: employee.ID, TagID: tags[j].ID})
					}
				}
			}
		}
	}

	err = s.repository.DetachTagsFromManyByIds(employeeIds)

	if err != nil {
		return err
	}

	if employeeTagModels != nil && len(employeeTagModels) > 0 {
		_, err := s.tagRepository.AttachManyEmployees(employeeTagModels)

		if err != nil {
			return err
		}
	}

	return nil
}
