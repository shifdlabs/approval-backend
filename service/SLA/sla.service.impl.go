package sla

import (
	"Microservice/model"
	appSettingsRepository "Microservice/repository/AppSettings"
	documentRepository "Microservice/repository/Document"
	documentHistoryRepository "Microservice/repository/DocumentHistory"
	"fmt"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SLAServiceImpl struct {
	AppSettingsRepository    appSettingsRepository.AppSettingsRepository
	DocumentRepository       documentRepository.DocumentRepository
	DocumentHistoryRepository documentHistoryRepository.DocumentHistoryRepository
}

func NewSLAServiceImpl(
	appSettingsRepo appSettingsRepository.AppSettingsRepository,
	documentRepo documentRepository.DocumentRepository,
	documentHistoryRepo documentHistoryRepository.DocumentHistoryRepository,
) SLAService {
	return &SLAServiceImpl{
		AppSettingsRepository:    appSettingsRepo,
		DocumentRepository:       documentRepo,
		DocumentHistoryRepository: documentHistoryRepo,
	}
}

func (s *SLAServiceImpl) RunAutoApprove() {
	fmt.Println("[SLA] Running auto-approve check...")

	maxDays := s.getSLAMaxDays()
	cutoff := time.Now().Add(-time.Duration(maxDays) * 24 * time.Hour)

	documents, err := s.DocumentRepository.GetAllInProgressForSLA()
	if err != nil {
		fmt.Printf("[SLA] Failed to get in-progress documents: %v\n", err.Message)
		return
	}

	count := 0
	for _, doc := range documents {
		if doc.UpdatedAt == nil || !doc.UpdatedAt.Before(cutoff) {
			continue
		}
		if s.autoApproveDocument(doc) {
			count++
		}
	}

	fmt.Printf("[SLA] Auto-approved %d document(s)\n", count)
}

func (s *SLAServiceImpl) getSLAMaxDays() int {
	setting, err := s.AppSettingsRepository.GetByKey("sla_max_days")
	if err != nil || setting == nil {
		return 7 // default fallback
	}
	days, parseErr := strconv.Atoi(setting.Value)
	if parseErr != nil || days <= 0 {
		return 7
	}
	return days
}

func (s *SLAServiceImpl) autoApproveDocument(doc model.Document) bool {
	// Find the current step's approver from preloaded sequences
	var approverUserID uuid.UUID
	found := false
	for _, seq := range doc.DocumentSequence {
		if seq.Step == doc.Step {
			approverUserID = seq.UserID
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("[SLA] Doc %s: no sequence found for step %d, skipping\n", doc.ID, doc.Step)
		return false
	}

	// Create history entry — approver is the actor, marked as system auto-approve
	history := model.DocumentHistory{
		Document:    &doc,
		UserID:      approverUserID,
		IsApproved:  true,
		Description: "Auto-approved by system (SLA exceeded)",
	}
	if errHistory := s.DocumentHistoryRepository.Create(history); errHistory != nil {
		fmt.Printf("[SLA] Doc %s: failed to create history: %v\n", doc.ID, errHistory.Message)
		return false
	}

	// Advance document
	totalSteps := len(doc.DocumentSequence)
	if (doc.Step + 1) <= totalSteps {
		doc.Step = doc.Step + 1
		doc.Status = 1
	} else {
		doc.Status = 2 // Finished
	}

	if errUpdate := s.DocumentRepository.Update(doc); errUpdate != nil {
		fmt.Printf("[SLA] Doc %s: failed to update document: %v\n", doc.ID, errUpdate.Message)
		return false
	}

	fmt.Printf("[SLA] Doc %s auto-approved at step %d (new step=%d, status=%d)\n",
		doc.ID, doc.Step-1, doc.Step, doc.Status)
	return true
}
