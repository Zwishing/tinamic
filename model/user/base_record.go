package user

import "time"

type BaseRecord struct {
	Created time.Time `json:"created"`
	Edited  time.Time `json:"edited"`
	Deleted bool      `json:"deleted,omitempty"`
}

func NewBaseRecord() *BaseRecord {
	return &BaseRecord{
		Created: time.Now(),
		Edited:  time.Now(),
		Deleted: false,
	}
}

type BaseRecorder struct {
	Creator string `json:"creator,omitempty"`
	Editor  string `json:"editor,omitempty"`
}

func NewBaseRecorder(creator string, editor string) *BaseRecorder {
	return &BaseRecorder{
		creator,
		editor,
	}
}
