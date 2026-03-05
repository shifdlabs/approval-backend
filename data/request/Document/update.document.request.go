package document

type UpdateDocumentRequest struct {
	Id                    string                      `json:"id"`
	AuthorID              string                      `validate:"required" json:"authorID"`
	PublicationNumberType int                         `json:"publicationNumberType"` // 1: Auto-Generated, 2: Booking Number, 3: Custom, 4: N/A (No Number)
	PublicationValue      *string                     `json:"publicationValue"`
	Type                  int                         `validate:"required" json:"type"`
	Priority              int                         `validate:"required" json:"priority"`
	Subject               string                      `validate:"required" json:"subject"`
	Body                  string                      `validate:"required" json:"body"`
	ExternalRecipient     string                      `json:"externalRecipient"`
	LetterHead            bool                        `json:"letterHead"`
	Recipients            []string                    `json:"recipients"`
	CarbonCopies          []string                    `json:"carbonCopies"`
	Sequences             []DocumentSequence          `validate:"required" json:"sequences"`
	NewAttachments        []DocumentAttachmentRequest `json:"newAttachments"`
	IsDraft               bool                        `json:"isDraft"`
	References            []string                    `json:"references"`
}
