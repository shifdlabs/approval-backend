package documentnumbers

type DocumentNumbersRequest struct {
	NumberingFormatID string `validate:"required" json:"numbering_format_id"`
}
