package user

import "time"

type BaseRecord struct {
	Created time.Time `json:"created"`
	Edited  time.Time `json:"edited"`
	Deleted bool      `json:"deleted"`
}

func NewBaseRecord() *BaseRecord {
	return &BaseRecord{
		Created: time.Now(),
		Edited:  time.Now(),
		Deleted: false,
	}
}

type BaseRecorder struct {
	Creator string `json:"creator"`
	Editor  string `json:"editor"`
}

func NewBaseRecorder(creator string, editor string) *BaseRecorder {
	return &BaseRecorder{
		creator,
		editor,
	}
}
