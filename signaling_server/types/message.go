package types

type IdentifyMessage struct {
	Type      string `json:"type"`       // must be "identify"
	Role      string `json:"role"`       // "mobile" or "desktop"
	SessionID string `json:"session_id"` // UUID shared by both sides
}
