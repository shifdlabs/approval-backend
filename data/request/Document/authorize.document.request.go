package document

type Authorize struct {
	DocumentID string `validate:"required" json:"documentId"`
	State      int    `validate:"required" json:"state"` // 1: approve, 2: reject,3: cancelled
	Comment    string `json:"comment"`
}
