package props

import (
	"errors"
)

type GetCard struct {
	UUID string `json:"uuid"`
}

func (c *GetCard) Validate() error {
	if c.UUID == "" {
		return errors.New("uuid is empty")
	}
	return nil
}
