package document

type CreateDocumentRequest struct {
	AuthorID              string                      `validate:"required" json:"authorID"`
	PublicationNumberType int                         `validate:"required" json:"publicationNumberType"` // 1: Auto-Generated, 2: Booking Number, 3: Custom, 4: N/A (No Number)
	PublicationValue      *string                     `json:"publicationValue"`                          // it could be Booked Number, Format ID or Custom Number
	Type                  int                         `validate:"required" json:"type"`
	Priority              int                         `validate:"required" json:"priority"`
	Subject               string                      `validate:"required" json:"subject"`
	Body                  string                      `validate:"required" json:"body"`
	ExternalRecipient     string                      `json:"externalRecipient"`
	Step                  int                         `validate:"required" json:"step"`
	LetterHead            bool                        `json:"letterHead"`
	Status                int                         `json:"status"`
	Recipients            []string                    `json:"recipients"`
	CarbonCopies          []string                    `json:"carbonCopies"`
	Sequences             []DocumentSequence          `validate:"required" json:"sequences"`
	Attachments           []DocumentAttachmentRequest `json:"attachments"`
	References            []string                    `json:"references"`
}

type DocumentSequence struct {
	UserID    string `validate:"required" json:"userID"`
	Signature bool   `json:"signature"`
}

type DocumentAttachmentRequest struct {
	OriginalName string `validate:"required" json:"originalName"`
	FileName     string `validate:"required" json:"fileName"`
	Path         string `validate:"required" json:"path"`
	Size         string `validate:"required" json:"size"`
	Type         string `validate:"required" json:"type"`
}
