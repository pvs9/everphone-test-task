package service

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/queue"
	"github.com/pvs9/everphone-test-task/pkg/repository"
	"github.com/pvs9/everphone-test-task/pkg/request"
	log "github.com/sirupsen/logrus"
)

type Employee interface {
	GetAll() ([]christmas.Employee, error)
	GetById(id int64) (*christmas.Employee, error)
	Update(employee christmas.Employee, employeeData request.EmployeeUpdateRequest) (christmas.Employee, error)
	UploadDataset(fileName string) (*string, error)
	ProcessDataset(fileName string) error
}

type Gift interface {
	GetAll() ([]christmas.Gift, error)
	GetById(id int64) (*christmas.Gift, error)
	Update(gift christmas.Gift, giftData request.GiftUpdateRequest) (christmas.Gift, error)
	UploadDataset(fileName string) (*string, error)
	ProcessDataset(fileName string) error
}

type Service struct {
	Employee
	Gift
}

type DatasetMessage struct {
	Type     string
	Filename string
}

func NewService(queues *queue.Queue, repositories *repository.Repository) *Service {
	return &Service{
		Employee: NewEmployeeService(queues.Publisher, repositories.Employee, repositories.Tag),
		Gift:     NewGiftService(queues.Publisher, repositories.Gift, repositories.Tag),
	}
}

func (s *Service) DatasetMessageHandler(m *sqs.Message) error {
	var messageBody DatasetMessage
	err := json.Unmarshal([]byte(*m.Body), &messageBody)

	if err != nil {
		return err
	}

	switch messageBody.Type {
	case "employee":
		err = s.Employee.ProcessDataset(messageBody.Filename)
	case "gift":
		err = s.Gift.ProcessDataset(messageBody.Filename)
	}

	if err != nil {
		return err
	}

	log.Infof("Dataset with type %s from message %s imported successfuly", messageBody.Type, *(m.MessageId))

	return nil
}
