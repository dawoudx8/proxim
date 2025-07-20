# Signaling Protocol – Proxim

This document defines the signaling protocol used by the Proxim system to establish peer-to-peer WebRTC sessions between the mobile application and the desktop application, via the embedded signaling server.

---

## Protocol Overview

- Transport: **WebSocket**
- Message Format: **JSON**
- Communication is routed through the signaling server
- Messages are exchanged between:
  - **Mobile App ↔ Signaling Server ↔ Desktop App**

---

## Connection Setup Flow

1. Desktop app launches and connects to signaling server
2. User creates a session → generates `sessionID` + `UUID`
3. Mobile app scans QR and connects to signaling server using `IP:port`
4. Mobile app sends `pair_request` with `sessionID`
5. Server forwards the request to matching desktop
6. Desktop app shows approval prompt (with timeout)
7. If approved, signaling messages (SDP, ICE) are exchanged
8. WebRTC connection is established directly between peers

---

## Message Format

All messages are sent over WebSocket in JSON format.
Each message must include a `type` field.

```json
{
  "type": "pair_request",
  "sessionID": "abc123",
  "deviceInfo": {
    "name": "Galaxy S23",
    "id": "device-xyz"
  }
}
```

---

## Message Types

### 1. `pair_request` (from mobile → server)
Sent after scanning QR and establishing WebSocket connection.
```json
{
  "type": "pair_request",
  "sessionID": "abc123",
  "deviceInfo": { "name": "Galaxy S23", "id": "device-xyz" }
}
```

### 2. `pair_prompt` (from server → desktop)
Sent by server to notify desktop of incoming request.
```json
{
  "type": "pair_prompt",
  "sessionID": "abc123",
  "deviceInfo": { "name": "Galaxy S23", "id": "device-xyz" }
}
```

### 3. `pair_response` (from desktop → server)
Sent after manual approval or rejection.
```json
{
  "type": "pair_response",
  "sessionID": "abc123",
  "approved": true
}
```

### 4. `pair_result` (from server → mobile)
Server informs mobile if request was accepted or rejected.
```json
{
  "type": "pair_result",
  "sessionID": "abc123",
  "status": "approved" // or "rejected"
}
```

### 5. `signal` (bidirectional)
Used for SDP and ICE exchange.
```json
{
  "type": "signal",
  "sessionID": "abc123",
  "payload": { "sdp": "..." } // or { "candidate": "..." }
}
```

---

## Timeouts & Cleanup

- Pairing requests expire after **40 seconds** if no response.
- Sessions expire after **15 minutes** or as per selected session type.
- Expired sessions are removed from server memory.

---

## Notes

- All connections are initiated by the client side (no push over raw TCP).
- The signaling server does not store any data permanently.
- The `sessionID` acts as the join key and must be unique.

