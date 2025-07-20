# Proxim (Remote Desktop Control)

Proxim is a privacy-focused remote desktop control system that allows a mobile device to securely connect to and control a desktop computer. The system is built with local-first principles, where all connections are peer-to-peer and initiated by the user via QR codes.

---

## Key Features

* Remote desktop control via mobile
* Secure peer-to-peer connection using WebRTC
* Session-based access with manual approval
* No account or login system required
* Cross-platform core:

  * Desktop (Rust)
  * Signaling Server (Go)
  * Mobile App (Kotlin + XML)

---

## Architecture Overview

```
 Mobile App         Signaling Server (Go)         Desktop App
     |                       |                          |
     |--- Pairing Request -->|                          |
     |                       |--> Notify Rust App --->  |
     |                       |<-- Accept/Reject ------- |
     |<------ WebRTC Negotiation via Signaling ------>  |
     |<-------------------- Direct P2P ---------------->|
```

* The signaling server is embedded as an executable and launched by the Rust desktop app.
* The mobile app communicates with the signaling server over the local network or internet (WebRTC + STUN).

---

## Happy Path Summary

1. User launches the desktop app.
2. User clicks "Generate QR Code".
3. User selects session type (One-time, Persistent, Temporary).
4. A QR code is generated containing UUID + session ID.
5. User scans QR code using the mobile app.
6. Mobile sends pairing request to the signaling server.
7. Desktop receives pairing request and displays device info.
8. User manually approves the request within 40 seconds.
9. Upon approval, full desktop control is granted to the mobile device.

> See `happy_path.md` for full internal flow.

---

## Session Types

* **One-time**: Expires after the first disconnect.
* **Persistent**: Remains until manually deleted.
* **Temporary**: Auto-expires after user-defined duration.

Each session is identified by a `sessionID` and bound to a `UUID` representing the desktop.

---

## Security Model

* No credentials or passwords are stored or exchanged.
* QR code acts as a one-time authentication token.
* Manual user approval is required before each connection.
* QR code and session auto-expire after 15 minutes or if ignored.

---

## Technologies Used

| Component          | Language / Stack             |
| ------------------ | ---------------------------- |
| Desktop App        | Rust                         |
| Signaling Server   | Go (`.exe` subprocess)       |
| Mobile App         | Kotlin + XML                 |
| Signaling Protocol | WebRTC + JSON over WebSocket |

---

## Development Status

Current phase: **Core signaling and session architecture**

Next steps:

* Finalize signaling message structure
* Implement in-memory session tracking
* Integrate WebRTC data channel establishment

---

## Internal Documentation

* [`ideal_flow.md`](./.docs/ideal_flow.md)
* [`system_architecture.md`](./.docs/system_architecture.md)
* [`signaling_protocol.md`](./.docs/signaling_protocol.md)
* [`session_management.md`](./.docs/session_management.md)
* [`webrtc_setup.md`](./.docs/webrtc_setup.md)
* [`tasks.md`](tasks.md)
