package document

import (
	request "Microservice/data/request/Document"
	response "Microservice/data/response/Document"
	"Microservice/helper"
	"Microservice/model"
	carbonCopyReposiory "Microservice/repository/CarbonCopy"
	repository "Microservice/repository/Document"
	documentAttachmentRepository "Microservice/repository/DocumentAttachment"
	documentHistoryReposiory "Microservice/repository/DocumentHistory"
	documentNumbersRepository "Microservice/repository/DocumentNumbers"
	documentReferenceRepository "Microservice/repository/DocumentReference"
	documentSequenceReposiory "Microservice/repository/DocumentSequence"
	recipientReposiory "Microservice/repository/Recipient"
	signatureRepository "Microservice/repository/Signature"
	userRepository "Microservice/repository/User"
	userLogRepository "Microservice/repository/UserLog"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type DocumentServiceImpl struct {
	DocumentRepository           repository.DocumentRepository
	UserRepository               userRepository.UserRepository
	DocumentSequenceRepository   documentSequenceReposiory.DocumentSequenceRepository
	DocumentAttachmentRepository documentAttachmentRepository.DocumentAttachmentRepository
	DocumentHistoryRepository    documentHistoryReposiory.DocumentHistoryRepository
	RecipientRepository          recipientReposiory.RecipientRepository
	CarbonCopyRepository         carbonCopyReposiory.CarbonCopyRepository
	UserLogRepository            userLogRepository.UserLogRepository
	DocumentNumbersRepository    documentNumbersRepository.DocumentNumbersRepository
	DocumentReferenceRepository  documentReferenceRepository.DocumentReferenceRepository
	SignatureRepository          signatureRepository.SignatureRepository
	Validate                     *validator.Validate
	Db                           *gorm.DB
}

func NewDocumentServiceImpl(
	documentRepository repository.DocumentRepository,
	userRepository userRepository.UserRepository,
	documentSequenceRepository documentSequenceReposiory.DocumentSequenceRepository,
	documentAttachmentRepository documentAttachmentRepository.DocumentAttachmentRepository,
	documentHistoryRepository documentHistoryReposiory.DocumentHistoryRepository,
	recipientRepository recipientReposiory.RecipientRepository,
	carbonCopyRepository carbonCopyReposiory.CarbonCopyRepository,
	userLogRepository userLogRepository.UserLogRepository,
	documentNumbersRepository documentNumbersRepository.DocumentNumbersRepository,
	documentReferenceRepository documentReferenceRepository.DocumentReferenceRepository,
	signatureRepository signatureRepository.SignatureRepository,
	Db *gorm.DB,
	validate *validator.Validate) DocumentService {
	return &DocumentServiceImpl{
		DocumentRepository:           documentRepository,
		UserRepository:               userRepository,
		DocumentSequenceRepository:   documentSequenceRepository,
		DocumentAttachmentRepository: documentAttachmentRepository,
		DocumentHistoryRepository:    documentHistoryRepository,
		RecipientRepository:          recipientRepository,
		CarbonCopyRepository:         carbonCopyRepository,
		UserLogRepository:            userLogRepository,
		DocumentNumbersRepository:    documentNumbersRepository,
		DocumentReferenceRepository:  documentReferenceRepository,
		SignatureRepository:          signatureRepository,
		Db:                           Db,
		Validate:                     validate,
	}
}

func (t DocumentServiceImpl) Create(request request.CreateDocumentRequest) (*model.Document, *helper.ErrorModel) {
	errStructure := t.Validate.Struct(request)
	if errStructure != nil {
		msg := "Structure Error"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Get User Data
	user, errUser := t.UserRepository.Get(request.AuthorID, true)
	if errUser != nil {
		return nil, nil
	}

	// Start Transaction
	trx := t.Db.Begin()

	newDocument, errConvert := t.convertRequestToCreateModel(request, user)
	if errConvert != nil {
		return nil, errConvert
	}

	t.DocumentRepository.Create(*trx, newDocument)

	if newDocument.ID == uuid.Nil {
		trx.Rollback()
		msg := "Document ID is nil after creation"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	// Store Internal Recipients
	if request.Recipients != nil {
		var recipients []model.Recipient
		for _, userId := range request.Recipients {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				trx.Rollback()
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errStructure, 500, &msg)
			}

			recipients = append(recipients, model.Recipient{
				Document: newDocument,
				UserID:   userUuid,
			})
		}

		t.RecipientRepository.Create(
			*trx,
			recipients,
		)
	}

	// Store Carbon Copies
	if request.CarbonCopies != nil {
		var carbonCopies []model.CarbonCopy
		for _, userId := range request.CarbonCopies {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				trx.Rollback()
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errStructure, 500, &msg)
			}

			carbonCopies = append(carbonCopies, model.CarbonCopy{
				Document: newDocument,
				UserID:   userUuid,
			})
		}

		t.CarbonCopyRepository.Create(
			*trx,
			carbonCopies,
		)
	}

	// Store Document Approvers
	for index, value := range request.Sequences {
		userId, _ := uuid.FromString(value.UserID)
		t.DocumentSequenceRepository.Create(
			trx,
			model.DocumentSequence{
				DocumentID: &newDocument.ID,
				UserID:     userId,
				Step:       (index + 1),
				Signature:  value.Signature,
			},
		)
	}

	// Store Document Sequences
	for _, value := range request.Attachments {
		t.DocumentAttachmentRepository.Create(
			trx,
			model.DocumentAttachment{
				Document:     newDocument,
				OriginalName: value.OriginalName,
				FileName:     value.FileName,
				Path:         value.Path,
				Size:         value.Size,
				Type:         value.Type,
			},
		)
	}

	if request.References != nil {
		for _, referenceID := range request.References {
			referenceUUID, _ := uuid.FromString(referenceID)
			t.DocumentReferenceRepository.Create(trx, model.DocumentReference{
				ReferenceID: referenceUUID,
				DocumentID:  newDocument.ID,
			})
		}
	}

	// End Transaction
	trx.Commit()
	return newDocument, nil
}

func (t DocumentServiceImpl) GetDocument(id string) (*response.DocumentResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentRepository.Get(id)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertDocumentToDocumentResponse(*document)
	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetDetailDocument(id string, currentUserId string) (*response.DocumentDetailResponse, *helper.ErrorModel) {
	helper.PrintValue("Rezz", id)
	document, fetchError := t.DocumentRepository.Get(id)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertToDocumentDetailResponse(*document, currentUserId)

	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetDetailForEdit(id string) (*response.EditDocumentResponse, *helper.ErrorModel) {
	document, fetchError := t.DocumentRepository.Get(id)

	if fetchError != nil {
		return nil, fetchError
	}

	documentResponse := t.convertDocumentToEditDocumentResponse(*document)

	return &documentResponse, fetchError
}

func (t DocumentServiceImpl) GetAllDocument() ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAll()

	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetAllReferences(query string) ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAllReferences(query)

	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetAllAuthorization(userId string) ([]response.DocumentResponse, *helper.ErrorModel) {

	result, fetchError := t.DocumentRepository.GetAllAuthorization(userId)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetAllInProgress(userId string) ([]response.DocumentResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetAllInProgress(userId)
	if fetchError != nil {
		return nil, fetchError
	} else {
		return t.mapDocumentsToDocumentResponse(result), nil
	}
}

func (t DocumentServiceImpl) GetDocumentStatistics(userId string) (*response.DocumentStatisticResponse, *helper.ErrorModel) {
	result, fetchError := t.DocumentRepository.GetDocumentStatistics(userId)

	if len(result) == 4 {
		return &response.DocumentStatisticResponse{
			Authorization: result[0],
			InProgress:    result[1],
			Rejected:      result[2],
			Completed:     result[3],
		}, nil
	} else {
		return nil, fetchError
	}
}

func (t DocumentServiceImpl) GetInProgressOverviewByDocId(documentId string) (*response.DocumentInProgressResponse, *helper.ErrorModel) {
	// InProgress Overview
	document, fetchDocErr := t.DocumentRepository.Get(documentId)
	approvers := []response.ApproverForOverview{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		// Get all sequences first (including approved ones)
		allSequences, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
		helper.PrintValue(len(allSequences), "Total sequences fetched")
		sequenceMap := make(map[string]*model.DocumentSequence)
		if allSequences != nil {
			for i, seq := range allSequences {
				helper.PrintValue(seq.UserID.String(), "Sequence UserID")
				helper.PrintValue(seq.Signature, "Sequence Signature value")
				sequenceMap[seq.UserID.String()] = &allSequences[i]
			}
		}

		history, fetchHistoryErr := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		addedUserIds := make(map[string]bool)

		if history != nil {
			for _, approver := range history {
				user, err := t.UserRepository.Get(approver.UserID.String(), true)
				if err != nil {
					return nil, err
				}

				addedUserIds[user.ID.String()] = true

				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String())
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				hasSigned := false
				if seq, exists := sequenceMap[user.ID.String()]; exists {
					hasSigned = seq.Signature
				}

				dateStr := approver.CreatedAt.String()
				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        user.Position.Name,
					Approved:     &approver.IsApproved,
					Date:         &dateStr,
					Signature:    hasSigned,
					SignatureUrl: signatureUrl,
				})
			}
		}

		if allSequences != nil {
			for _, sequence := range allSequences {

				if addedUserIds[sequence.UserID.String()] {
					continue
				}

				user, err := t.UserRepository.Get(sequence.UserID.String(), true)
				if err != nil {
					return nil, err
				}

				// Fetch signature if exists
				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String())
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        user.Position.Name,
					Approved:     nil,
					Date:         nil,
					Signature:    sequence.Signature,
					SignatureUrl: signatureUrl,
				})
			}
		}
	}

	if document != nil {
		return &response.DocumentInProgressResponse{
			Subject:   document.Subject,
			Approvers: approvers,
		}, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetInProgressOverview(userId string) (*response.DocumentInProgressResponse, *helper.ErrorModel) {
	// InProgress Overview
	document, fetchDocErr := t.DocumentRepository.GetOneLatestInprogress(userId)

	approvers := []response.ApproverForOverview{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		allSequences, _ := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
		sequenceMap := make(map[string]*model.DocumentSequence)
		if allSequences != nil {
			for i, seq := range allSequences {
				sequenceMap[seq.UserID.String()] = &allSequences[i]
			}
		}

		history, fetchHistoryErr := t.DocumentHistoryRepository.GetAllHistoryByDocumentId(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		// Track which users are already added from history
		addedUserIds := make(map[string]bool)

		if history != nil {
			for _, approver := range history {
				user, err := t.UserRepository.Get(approver.UserID.String(), true)
				if err != nil {
					return nil, err
				}

				// Mark this user as added
				addedUserIds[user.ID.String()] = true

				// Fetch signature if exists
				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String())
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				// Get sequence to check if they signed
				hasSigned := false
				if seq, exists := sequenceMap[user.ID.String()]; exists {
					hasSigned = seq.Signature
				}

				dateStr := approver.CreatedAt.String()
				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        user.Position.Name,
					Approved:     &approver.IsApproved,
					Date:         &dateStr,
					Signature:    hasSigned,
					SignatureUrl: signatureUrl,
				})
			}
		}

		if allSequences != nil {
			for _, sequence := range allSequences {

				if addedUserIds[sequence.UserID.String()] {
					continue
				}

				user, err := t.UserRepository.Get(sequence.UserID.String(), true)
				if err != nil {
					return nil, err
				}

				var signatureUrl *string
				signature, _ := t.SignatureRepository.GetByUserId(user.ID.String())
				if signature != nil {
					signatureUrl = &signature.ImageURL
				}

				approvers = append(approvers, response.ApproverForOverview{
					Name:         user.FirstName + " " + user.LastName,
					Title:        user.Position.Name,
					Approved:     nil,
					Date:         nil,
					Signature:    sequence.Signature,
					SignatureUrl: signatureUrl,
				})
			}
		}
	}

	if document != nil {
		return &response.DocumentInProgressResponse{
			Subject:   document.Subject,
			Approvers: approvers,
		}, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetRejectedOverview(userId string) (*response.RejectedOverviewResponse, *helper.ErrorModel) {
	document, fetchDocErr := t.DocumentRepository.GetLastestRejected(userId)

	rejected := response.RejectedOverviewResponse{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		rejectedBy, fetchHistoryErr := t.DocumentHistoryRepository.GetLastRejection(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		if rejectedBy != nil {
			user, err := t.UserRepository.Get(rejectedBy.UserID.String(), true)
			if err != nil {
				return nil, err
			}

			rejected = response.RejectedOverviewResponse{
				Name:    user.FirstName + " " + user.LastName,
				Title:   user.Position.Name,
				Subject: document.Subject,
				Reason:  rejectedBy.Description,
				Date:    rejectedBy.CreatedAt.String(),
			}
		}
	}

	if document != nil {
		return &rejected, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetCompletedOverview(userId string) (*response.CompletedOverviewResponse, *helper.ErrorModel) {
	document, fetchDocErr := t.DocumentRepository.GetLastestCompleted(userId)

	completed := response.CompletedOverviewResponse{}

	if fetchDocErr != nil {
		return nil, fetchDocErr
	}

	if document != nil {
		history, fetchHistoryErr := t.DocumentHistoryRepository.GetLastApprover(document.ID.String())
		if fetchHistoryErr != nil {
			return nil, fetchDocErr
		}

		user, err := t.UserRepository.Get(history.UserID.String(), true)
		if err != nil {
			return nil, err
		}

		recipients, fetchRecipientErr := t.RecipientRepository.GetRecipientsByDocId(document.ID.String())
		if fetchRecipientErr != nil {
			return nil, fetchDocErr
		}

		internalRecipients := []response.InternalRecipientForOverview{}
		if recipients != nil {
			for _, recipient := range recipients {
				user, err := t.UserRepository.Get(recipient.UserID.String(), true)
				if err != nil {
					return nil, err
				}
				internalRecipients = append(internalRecipients, response.InternalRecipientForOverview{
					Name:  user.FirstName + " " + user.LastName,
					Title: user.Position.Name,
				})
			}
		}

		completed = response.CompletedOverviewResponse{
			IsFinished:        document.Status == 2,
			Name:              user.FirstName + " " + user.LastName,
			Title:             user.Position.Name,
			Subject:           document.Subject,
			Date:              history.CreatedAt.String(),
			InternalRecipient: internalRecipients,
			ExternalRecipient: &document.ExternalRecipient,
		}
	}

	if document != nil {
		return &completed, nil
	} else {
		return nil, nil
	}
}

func (t DocumentServiceImpl) GetAllInbox(userId string) ([]response.DocumentResponse, *helper.ErrorModel) {
	response, fetchError := t.DocumentRepository.GetAllInbox(userId)

	if fetchError != nil {
		return nil, fetchError
	}

	return t.mapDocumentsToDocumentResponse(response), nil
}

func (t DocumentServiceImpl) Update(request request.UpdateDocumentRequest) (*model.Document, *helper.ErrorModel) {
	errStructure := t.Validate.Struct(request)

	if errStructure != nil {
		msg := "Structure Error"
		return nil, helper.ErrorCatcher(errStructure, 500, &msg)
	}

	trx := t.Db.Begin()

	document, err := t.DocumentRepository.Get(request.Id)
	if err != nil {
		msg := "Document Not Found"
		return nil, helper.ErrorCatcher(err, 500, &msg)
	}

	document.Type = request.Type
	document.Priority = request.Priority
	document.Subject = request.Subject
	document.Body = request.Body
	document.ExternalRecipient = request.ExternalRecipient
	document.LetterHead = request.LetterHead

	if document.PublicationNumberType == request.PublicationNumberType {
		if request.PublicationNumberType == 3 {
			document.CustomPublicationNumber = request.PublicationValue
		}
	} else {
		document.PublicationNumberType = request.PublicationNumberType
		switch request.PublicationNumberType {
		case 3:
			document.CustomPublicationNumber = request.PublicationValue
		case 4:
			document.CustomPublicationNumber = nil
		}
	}

	if request.IsDraft {
		document.Status = 0
	} else {
		document.Status = 1
	}

	t.DocumentRepository.Update(*document)
	t.DocumentReferenceRepository.Update(request.References, document.ID)

	// Update Internal Recipients
	if len(request.Recipients) > 0 {
		var recipients []model.Recipient
		for _, userId := range request.Recipients {
			userUuid, _ := uuid.FromString(userId)

			recipients = append(recipients, model.Recipient{
				Document:   document,
				UserID:     userUuid,
				DocumentID: document.ID,
			})
		}

		err := t.RecipientRepository.Update(
			*document,
			recipients,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	if len(request.CarbonCopies) > 0 {
		var carbonCopies []model.CarbonCopy
		for _, userId := range request.CarbonCopies {
			userUuid, errorParse := uuid.FromString(userId)
			if errorParse != nil {
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errorParse, 500, &msg)
			}

			carbonCopies = append(carbonCopies, model.CarbonCopy{
				Document: document,
				UserID:   userUuid,
			})
		}

		err := t.CarbonCopyRepository.Update(
			*document,
			carbonCopies,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	if len(request.Sequences) > 0 {
		var sequences []model.DocumentSequence
		for index, sequence := range request.Sequences {
			userUuid, errorParse := uuid.FromString(sequence.UserID)
			if errorParse != nil {
				msg := "Structure Error"
				return nil, helper.ErrorCatcher(errorParse, 500, &msg)
			}

			sequences = append(sequences, model.DocumentSequence{
				DocumentID: &document.ID,
				UserID:     userUuid,
				Step:       index + 1,
				Signature:  sequence.Signature,
			})
		}

		err := t.DocumentSequenceRepository.Update(
			*document,
			sequences,
		)

		if err != nil {
			msg := "Structure Error"
			return nil, helper.ErrorCatcher(err, 500, &msg)
		}
	}

	for _, value := range request.NewAttachments {
		t.DocumentAttachmentRepository.Create(
			trx,
			model.DocumentAttachment{
				Document:     document,
				OriginalName: value.OriginalName,
				FileName:     value.FileName,
				Path:         value.Path,
				Size:         value.Size,
				Type:         value.Type,
			},
		)
	}

	trx.Commit()
	return document, nil
}

func (t DocumentServiceImpl) Authorize(request request.Authorize, userId string) *helper.ErrorModel {
	errStructure := t.Validate.Struct(request)

	if errStructure != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	document, err := t.DocumentRepository.Get(request.DocumentID)
	if err != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	sequences, fetchSequenceErr := t.DocumentSequenceRepository.GetSequencesByDocumentId(document.ID.String())
	if fetchSequenceErr != nil {
		msg := "Structure Error"
		return helper.ErrorCatcher(errStructure, 500, &msg)
	}

	if document != nil && len(sequences) > 0 {
		userIdUUID, _ := uuid.FromString(userId)
		if request.State == 1 { // Approved
			// Check if user has signature
			hasSignature := false
			_, errSignature := t.SignatureRepository.GetByUserId(userId)
			if errSignature == nil {
				hasSignature = true
			}

			for i, seq := range sequences {
				if seq.UserID.String() == userId && seq.Step == document.Step {
					sequences[i].Signature = hasSignature
					errUpdateSeq := t.Db.Save(&sequences[i]).Error
					if errUpdateSeq != nil {
						msg := "Failed to update signature status"
						return helper.ErrorCatcher(errUpdateSeq, 500, &msg)
					}
					break
				}
			}

			if (document.Step + 1) <= len(sequences) {
				document.Status = 1
				document.Step = (document.Step + 1)
			} else {
				document.Status = 2
			}
		} else if request.State == 2 { // Rejected
			document.Status = 99
		} else if request.State == 3 { // Cancelled
			document.Status = 3
		}

		isApproved := request.State == 1

		errResponse := t.DocumentHistoryRepository.Create(
			model.DocumentHistory{
				Document:    document,
				Description: request.Comment,
				UserID:      userIdUUID,
				IsApproved:  isApproved,
			},
		)

		if errResponse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}

		errDocumentResoponse := t.DocumentRepository.Update(*document) // There's problem, the step not updated!
		if errDocumentResoponse != nil {
			msg := "Structure Error"
			return helper.ErrorCatcher(errStructure, 500, &msg)
		}
	}

	return nil
}

// // User Log
// t.UserLogRepository.Create(
// 	model.UserLog{
// 		UserID: userId,
// 		Action: string(enums.Approve),
// 		Module: string(enums.Document),
// 		Log:    helper.ToJSON(request),
// 	},
// )

func (t DocumentServiceImpl) GetCompleteByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen berdasarkan AuthorID dari repository

	//fmt.Println("Executing query for AuthorID:", authorID)
	documents, fetchError := t.DocumentRepository.GetCompleteByAuthorID(authorID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetDraftByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetDraftByAuthorID(authorID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetRejectedByAuthorID(authorID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetRejectedByAuthorID(authorID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}

func (t DocumentServiceImpl) GetAllAuthorDocuments(authorID string) ([]response.DocumentResponse, *helper.ErrorModel) {
	// Ambil data dokumen draft berdasarkan AuthorID dari repository
	documents, fetchError := t.DocumentRepository.GetAllAuthorDocuments(authorID)
	if fetchError != nil {
		return nil, fetchError
	}

	documentResponses := t.mapDocumentsToDocumentResponse(documents)

	return documentResponses, nil
}
