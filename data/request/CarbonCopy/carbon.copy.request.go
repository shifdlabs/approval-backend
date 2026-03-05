package position

type CarbonCopyRequest struct {
	DocumentId string   `validate:"required,min=1,max=200" json:"documentId"`
	UserIds    []string `validate:"required" json:"userIds"`
}
