package model

type (
	// UstreamStatus represents ustream status
	UstreamStatus struct {
		ID     string `datastore:"-" goon:"id" validate:"required"`
		IsLive bool
	}

	// UstreamStatusRepository interface
	UstreamStatusRepository interface {
		Find(string) (*UstreamStatus, error)
		Save(*UstreamStatus) error
	}

	// UstreamStatusChecker interface
	UstreamStatusChecker interface {
		IsLive() (bool, error)
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
