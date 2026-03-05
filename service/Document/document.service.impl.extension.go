package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
	"fmt"
)

func (t DocumentServiceImpl) mapDocumentsToDocumentResponse(documents []model.Document) []response.DocumentResponse {
	responseDocuments := make([]response.DocumentResponse, len(documents))
	for i, document := range documents {
		responseDocuments[i] = t.convertDocumentToDocumentResponse(document)
	}
	return responseDocuments
}

func (t DocumentServiceImpl) convertDocumentToDocumentResponse(document model.Document) response.DocumentResponse {

	var currentApproverTitle *string
	var lastRejector *response.RejectorResponse

	if len(document.DocumentSequence) > 0 {
		currentApprover, _ := t.DocumentSequenceRepository.GetCurrentApprover(document.ID.String())

		user, _ := t.UserRepository.Get(currentApprover.UserID.String(), true)

		title := fmt.Sprintf("%s %s - %s",
			user.FirstName,
			user.LastName,
			user.Position.Name,
		)

		currentApproverTitle = &title
	}

	lastRejectorResponse, _ := t.DocumentHistoryRepository.GetLastRejection(document.ID.String())
	if lastRejectorResponse != nil {
		rejectorData, _ := t.UserRepository.Get(string(lastRejectorResponse.UserID.String()), true)
		rejectorName := fmt.Sprintf("%s %s",
			rejectorData.FirstName,
			rejectorData.LastName,
		)

		lastRejector = &response.RejectorResponse{
			Name:   &rejectorName,
			Reason: &lastRejectorResponse.Description,
		}
	}

	// Perform necessary conversion logic here, potentially selecting specific fields
	responseDocument := response.DocumentResponse{
		Id:                  &document.ID,
		Subject:             document.Subject,
		Body:                document.Body,
		Type:                document.Type,
		Step:                document.Step,
		Status:              document.Status,
		Priority:            document.Priority,
		Author:              document.Author,
		DocumentSequence:    document.DocumentSequence,
		DocumentHistory:     document.DocumentHistory,
		DocumentAttachment:  document.DocumentAttachment,
		CreatedAt:           *document.CreatedAt,
		UpdatedAt:           *document.UpdatedAt,
		CurrentApprovalName: currentApproverTitle,
		LastRejector:        lastRejector,
	}

	return responseDocument
}

func (t DocumentServiceImpl) convertToDocumentDetailResponse(document model.Document, userId string) response.DocumentDetailResponse {
	inProgressOverview, _ := t.GetInProgressOverviewByDocId(document.ID.String())

	documentSequence, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
	currentApprover := t.getCurrentApprover(documentSequence, document)
	documentHistories := t.getDocumentHistory(document)
	documentAttachment := t.getDocumentAttachment(document)
	recipients := t.getInternalRecipients(document.ID.String())
	referencesResult, _ := t.DocumentReferenceRepository.GetAll(document.ID)
	documentReferences := make([]response.DocumentReference, len(referencesResult))

	if len(referencesResult) > 0 {
		for i, reference := range referencesResult {
			helper.PrintValue(reference.DocumentID, "Reference ID")
			document, err := t.DocumentRepository.Get(reference.DocumentID.String())
			if err != nil || document == nil {
				// Skip this reference, move on to the next
				continue
			}

			documentReferences[i] = response.DocumentReference{
				Id:      document.ID.String(),
				Subject: document.Subject,
			}
		}
	}

	var publicationValue string
	switch document.PublicationNumberType {
	case 1, 2:
		response, _ := t.DocumentNumbersRepository.GetByDocumentID(document.ID)
		if response != nil {
			publicationValue = response.Value
		}
	case 3:
		publicationValue = *document.CustomPublicationNumber

	case 4:
		publicationValue = ""
	}

	response := response.DocumentDetailResponse{
		Id:                 &document.ID,
		PublicationValue:   publicationValue,
		ExternalRecipient:  document.ExternalRecipient,
		Subject:            document.Subject,
		Body:               document.Body,
		Type:               document.Type,
		Step:               document.Step,
		Status:             document.Status,
		Priority:           document.Priority,
		Author:             *document.Author,
		DocumentSequence:   *inProgressOverview,
		DocumentHistory:    &documentHistories,
		DocumentAttachment: &documentAttachment,
		DocumentReferences: &documentReferences,
		InternalRecipients: &recipients,
		CreatedAt:          *document.CreatedAt,
		UpdatedAt:          *document.UpdatedAt,
		IsApprover:         currentApprover.UserID.String() == userId,
		IsAllowToUpdate:    document.Author.ID.String() == userId && document.Status == 99,
	}

	return response
}

func (t DocumentServiceImpl) getInternalRecipients(documentId string) []response.InternalRecipient {
	recipientsResponse, _ := t.RecipientRepository.GetRecipientsByDocId(documentId)
	recipients := make([]response.InternalRecipient, len(recipientsResponse))
	for i, recipient := range recipientsResponse {
		user, _ := t.UserRepository.Get(recipient.UserID.String(), true)

		recipients[i] = response.InternalRecipient{
			Name:  user.FirstName + " " + user.LastName,
			Title: user.Position.Name,
		}
	}

	return recipients
}

func (t DocumentServiceImpl) getDocumentAttachment(document model.Document) []response.DocumentAttachment {
	documentAttachment := make([]response.DocumentAttachment, len(document.DocumentAttachment))
	for i, attachment := range document.DocumentAttachment {
		documentAttachment[i] = response.DocumentAttachment{
			Id:           attachment.ID.String(),
			OriginalName: attachment.OriginalName,
			FileName:     attachment.FileName,
			Path:         attachment.Path,
			Size:         attachment.Size,
			Type:         attachment.Type,
		}
	}

	return documentAttachment
}

func (t DocumentServiceImpl) getDocumentHistory(document model.Document) []response.DocumentHistory {
	documentHistories := make([]response.DocumentHistory, len(document.DocumentHistory))
	for i, history := range document.DocumentHistory {
		user, _ := t.UserRepository.Get(history.UserID.String(), true)

		documentHistories[i] = response.DocumentHistory{
			Name:       user.FirstName + " " + user.LastName,
			Title:      user.Position.Name,
			IsApproved: history.IsApproved,
			Reason:     history.Description,
			UpdatedAt:  history.CreatedAt.String(),
		}
	}

	return documentHistories
}

func (t DocumentServiceImpl) getApprovers(document model.Document) []string {
	approverIds := make([]string, len(document.DocumentSequence))
	for i, history := range document.DocumentSequence {
		user, _ := t.UserRepository.Get(history.UserID.String(), true)

		approverIds[i] = user.ID.String()
	}

	return approverIds
}

func (t DocumentServiceImpl) getCurrentApprover(sequences []model.DocumentSequence, document model.Document) model.DocumentSequence {
	currentApprover := model.DocumentSequence{}
	if len(sequences) > 0 {
		currentApprover = sequences[document.Step-1]
	}

	return currentApprover
}

func (t DocumentServiceImpl) convertRequestToCreateModel(documentRequest request.CreateDocumentRequest, user *model.User) (*model.Document, *helper.ErrorModel) {
	var customPublicationCode *string = nil
	if documentRequest.PublicationNumberType == 3 {
		customPublicationCode = documentRequest.PublicationValue
	}

	// Store to DB
	document := model.Document{
		Author:                  user,
		PublicationNumberType:   documentRequest.PublicationNumberType,
		CustomPublicationNumber: customPublicationCode,
		Type:                    documentRequest.Type,
		Priority:                documentRequest.Priority,
		Subject:                 documentRequest.Subject,
		Body:                    documentRequest.Body,
		ExternalRecipient:       documentRequest.ExternalRecipient,
		Step:                    documentRequest.Step,
		LetterHead:              documentRequest.LetterHead,
		Status:                  documentRequest.Status,
	}

	return &document, nil
}

func (t DocumentServiceImpl) convertDocumentToEditDocumentResponse(document model.Document) response.EditDocumentResponse {

	internalRecipient, _ := t.RecipientRepository.GetRecipientsByDocId(document.ID.String())
	carbonCopy, _ := t.CarbonCopyRepository.GetCarbonCopysByDocId(document.ID.String())
	documentAttachment := t.getDocumentAttachment(document)
	approvers := t.getApprovers(document)
	recipientIds := make([]string, 0, len(internalRecipient))
	for _, r := range internalRecipient {
		recipientIds = append(recipientIds, string(r.UserID.String()))
	}

	carbonCopiesIds := make([]string, 0, len(carbonCopy))
	for _, r := range internalRecipient {
		carbonCopiesIds = append(carbonCopiesIds, string(r.UserID.String()))
	}

	referencesResult, _ := t.DocumentReferenceRepository.GetAll(document.ID)
	documentReferences := make([]response.DocumentReference, len(referencesResult))

	if len(referencesResult) > 0 {
		for i, reference := range referencesResult {
			document, err := t.DocumentRepository.Get(reference.DocumentID.String())
			if err != nil || document == nil {
				// Skip this reference, move on to the next
				continue
			}

			documentReferences[i] = response.DocumentReference{
				Id:      document.ID.String(),
				Subject: document.Subject,
			}
		}
	}

	var publicationNumber string

	switch document.PublicationNumberType {
	case 1, 2:
		documentNumber, _ := t.DocumentNumbersRepository.GetByDocumentID(document.ID)
		publicationNumber = documentNumber.Value
	case 3:
		publicationNumber = *document.CustomPublicationNumber
	default:
		publicationNumber = ""
	}

	response := response.EditDocumentResponse{
		Id:                    &document.ID,
		PublicationNumberType: document.PublicationNumberType,
		PublicationValue:      &publicationNumber,
		DocumentReferences:    &documentReferences,
		Subject:               document.Subject,
		Body:                  document.Body,
		Type:                  document.Type,
		Step:                  document.Step,
		Status:                document.Status,
		Priority:              document.Priority,
		Author:                *document.Author,
		DocumentAttachment:    &documentAttachment,
		ExternalRecipient:     &document.ExternalRecipient,
		InternalRecipients:    &recipientIds,
		CarbonCopy:            &carbonCopiesIds,
		Approvers:             &approvers,
	}

	return response
}
