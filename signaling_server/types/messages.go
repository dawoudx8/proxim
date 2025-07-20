package types

type InitMessage struct {
	Type string `json:"type"`
	Role string `json:"role"`
}

type PairRequest struct {
	Type       string     `json:"type"`
	SessionID  string     `json:"sessionID"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}
type DeviceInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type PairPrompt struct {
	Type       string     `json:"type"`
	SessionID  string     `json:"sessionID"`
	DeviceInfo DeviceInfo `json:"deviceInfo"`
}

type PairResponse struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionID"`
	Approved  bool   `json:"approved"`
}

type PairResult struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionID"`
	Status    string `json:"status"` // "approved" or "rejected"
}
