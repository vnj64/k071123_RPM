package props

import "github.com/google/uuid"

type UpdateSessionPaid struct {
	SessionUUID uuid.UUID `json:"session_uuid"`
}
