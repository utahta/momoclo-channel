package entity

type (
	// UstreamStatus represents ustream status
	UstreamStatus struct {
		ID     string `datastore:"-" goon:"id" validate:"required"`
		IsLive bool
	}
)

const (
	UstreamStatusID = "ustream_status"
)

// NewUstreamStatus returns UstreamStatus
func NewUstreamStatus() *UstreamStatus {
	return &UstreamStatus{
		ID:     UstreamStatusID,
		IsLive: false,
	}
}
