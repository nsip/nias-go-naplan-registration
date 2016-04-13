// shared types for passing validation errors
// between services
// and for process management/status messages

package lib

type ValidationError struct {
	TxID         string `json:"txID"`         // transaction this record is part of - assigned on first read by csvserver
	Field        string `json:"errField"`     // the field that has an error
	Description  string `json:"description"`  // error description
	OriginalLine string `json:"originalLine"` // input file record line that has the error
	Vtype        string `json:"validationType"`
}

type ProcessingNotification struct {
	TxID  string
	Vtype string
}

type TransactionSummary struct {
	TxID        string
	RecordCount int
}
