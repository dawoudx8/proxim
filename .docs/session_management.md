# Session Management – Proxim

This document describes how sessions are managed by the Proxim signaling server. Sessions define the scope, duration, and approval status of each mobile-to-desktop connection attempt.

---

## Session Structure

Each session is stored in memory and represented by the following structure:

```json
{
  "sessionID": "abc123",
  "uuid": "desktop-uuid-456",
  "createdAt": "2025-07-17T19:00:00Z",
  "expiresAt": "2025-07-17T19:15:00Z",
  "status": "pending", // or "approved", "rejected", "expired"
  "deviceInfo": null
}
```

- `sessionID`: Unique per session, generated on QR creation
- `uuid`: Static identifier of the desktop device
- `createdAt`: Time session was created
- `expiresAt`: Calculated based on session type
- `status`: Session state
- `deviceInfo`: Populated upon pairing request

---

## Session Types

| Type        | Description                                   | Expiration Logic               |
|-------------|-----------------------------------------------|---------------------------------|
| One-time    | Expires after first disconnect                | Marked expired on disconnect   |
| Persistent  | Stays active until user deletes it manually   | No automatic expiration        |
| Temporary   | User defines a custom time limit              | Expires at `expiresAt`         |

Session type is selected by the user before generating the QR code.

---

## Lifecycle Flow

1. Rust app generates `sessionID` + `uuid` and registers session in server memory
2. Session stored with `status = pending`, `expiresAt` computed by type
3. When mobile sends `pair_request`, `deviceInfo` is attached to the session
4. Server forwards `pair_prompt` to the desktop
5. If user approves within 40 seconds, session status becomes `approved`
6. If rejected or timeout occurs, session status becomes `rejected` or `expired`
7. Approved sessions allow WebRTC signaling
8. Session is cleaned up:
   - On manual close
   - On disconnect (for one-time)
   - On expiration (temporary)

---

## Cleanup Mechanism

- The signaling server runs a background cleanup loop (e.g. every 30 seconds)
- It checks current time against `expiresAt`
- Expired or inactive sessions are removed from memory
- One-time sessions are removed after disconnect

---

## Security Notes

- All sessions are volatile and stored in RAM only
- No user data or session logs are persisted
- QR codes are single-use tokens scoped to one sessionID

---

## Future Considerations

- Optional support for persistent sessions with manual listing/deletion
- Optional in-memory session caching for reconnect use cases
- Potential integration with external database (if needed)

