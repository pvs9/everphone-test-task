package service

import (
	"encoding/json"
	christmas "github.com/pvs9/everphone-test-task"
	"github.com/pvs9/everphone-test-task/pkg/queue"
	"github.com/pvs9/everphone-test-task/pkg/repository"
	"github.com/pvs9/everphone-test-task/pkg/request"
	"io/ioutil"
)

type GiftService struct {
	publisher     queue.Publisher
	repository    repository.Gift
	tagRepository repository.Tag
}

func NewGiftService(publisher queue.Publisher, repository repository.Gift, tagRepository repository.Tag) *GiftService {
	return &GiftService{publisher: publisher, repository: repository, tagRepository: tagRepository}
}

func (s *GiftService) GetAll() ([]christmas.Gift, error) {
	gifts, err := s.repository.GetAll()

	if err != nil {
		return gifts, err
	}

	return gifts, nil
}

func (s *GiftService) GetById(id int64) (*christmas.Gift, error) {
	gift, err := s.repository.GetById(id)

	if err != nil {
		return nil, err
	}

	return gift, nil
}

func (s *GiftService) Update(gift christmas.Gift, giftData request.GiftUpdateRequest) (christmas.Gift, error) {
	tags, err := s.tagRepository.GetByNames(giftData.Tags, true)

	if err != nil {
		return gift, err
	}

	gift.Name = giftData.Name
	gift.Tags = tags

	_, err = s.repository.Update(gift)

	if err != nil {
		return gift, err
	}

	gift, err = s.repository.SyncTags(gift, tags)

	return gift, nil
}

func (s *GiftService) UploadDataset(fileName string) (*string, error) {
	message := DatasetMessage{Type: "gift", Filename: fileName}
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

func (s *GiftService) ProcessDataset(fileName string) error {
	dataset, err := ioutil.ReadFile("/go/src/everphone-test-task.io/pkg/storage/" + fileName)

	if err != nil {
		return err
	}

	var data []request.GiftDatasetEntity
	err = json.Unmarshal(dataset, &data)
	var giftNames []string
	var tagNames []string
	giftTagsMap := make(map[string][]string)

	for index := range data {
		giftNames = append(giftNames, data[index].Name)
		giftTagsMap[data[index].Name] = data[index].Tags
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

	gifts, err := s.repository.GetByNames(giftNames, true)

	if err != nil {
		return err
	}

	tags, err := s.tagRepository.GetByNames(uniqueTagNames, true)

	if err != nil {
		return err
	}

	var giftTagModels []christmas.GiftTag
	var giftIds []int64

	for _, gift := range gifts {
		giftTagNames, ok := giftTagsMap[gift.Name]
		giftIds = append(giftIds, gift.ID)

		if ok && giftTagNames != nil && len(giftTagNames) > 0 {
			for i := range giftTagNames {
				for j := range tags {
					if tags[j].Name == giftTagNames[i] {
						giftTagModels = append(giftTagModels, christmas.GiftTag{GiftID: gift.ID, TagID: tags[j].ID})
					}
				}
			}
		}
	}

	err = s.repository.DetachTagsFromManyByIds(giftIds)

	if err != nil {
		return err
	}

	if giftTagModels != nil && len(giftTagModels) > 0 {
		_, err := s.tagRepository.AttachManyGifts(giftTagModels)

		if err != nil {
			return err
		}
	}

	return nil
}
