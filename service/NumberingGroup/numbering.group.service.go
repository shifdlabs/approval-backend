package numberinggroup

import (
	request "Microservice/data/request/NumberingGroup"   // Untuk NumberingGroupRequest
	response "Microservice/data/response/NumberingGroup" // Untuk NumberingGroupResponse
	"Microservice/model"

	// Untuk UserResponse
	"Microservice/helper"
	repository "Microservice/repository/NumberingGroup" // Untuk NumberingGroupRepository

	"github.com/go-playground/validator/v10"
)

type NumberingGroupServiceImpl struct {
	NumberingGroupRepository repository.NumberingGroupRepository
	Validate                 *validator.Validate
}

func NewNumberingGroupServiceImpl(
	documentRepository repository.NumberingGroupRepository,
	validate *validator.Validate) NumberingGroupService {
	return &NumberingGroupServiceImpl{
		NumberingGroupRepository: documentRepository,
		Validate:                 validate,
	}
}

func (t NumberingGroupServiceImpl) Create(request request.NumberingGroupRequest) *helper.ErrorModel {
	data := model.NumberingGroup{
		Name:        request.Name,
		Description: request.Description,
	}

	fetchError := t.NumberingGroupRepository.Create(data)
	if fetchError != nil {
		return fetchError
	}

	return nil
}

func (t NumberingGroupServiceImpl) Get(id string) (*response.NumberingGroupResponse, *helper.ErrorModel) {
	document, fetchError := t.NumberingGroupRepository.Get(id)
	if fetchError != nil {
		return nil, fetchError
	}

	if document == nil {
		return nil, nil
	}

	documentResponse := t.convertNumberingGroupToNumberingGroupResponse(*document)

	return &documentResponse, fetchError
}

func (t NumberingGroupServiceImpl) GetAll() ([]response.NumberingGroupResponse, *helper.ErrorModel) {
	result, fetchError := t.NumberingGroupRepository.GetAll()
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapNumberingGroupToNumberingGroupResponse(result), nil
	}
}

func (t NumberingGroupServiceImpl) Delete(id string) *helper.ErrorModel {
	errResponse := t.NumberingGroupRepository.Delete(id)
	if errResponse != nil {
		return errResponse
	}

	return nil
}
