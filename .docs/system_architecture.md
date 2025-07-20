# System Architecture – Proxim

This document describes the high-level system architecture of the Proxim remote desktop control system. It outlines how the desktop application, mobile application, and embedded signaling server interact to establish a secure peer-to-peer control session.

---

## Components Overview

### 1. Desktop Application (Rust)
- Provides GUI for user interaction
- Generates and displays QR codes
- Launches the embedded signaling server (`proxim-server.exe`)
- Connects to the signaling server as a desktop client
- Receives pairing requests and manages session approval
- Executes remote control commands (mouse, keyboard, system operations)

### 2. Signaling Server (Go)
- Runs as a subprocess launched by the desktop app
- Handles session registration, pairing, and expiration
- Forwards signaling messages between desktop and mobile
- Serves as a relay for WebRTC negotiation only (not data transfer)

### 3. Mobile Application (Kotlin + XML)
- Scans the desktop-generated QR code
- Connects to the signaling server using IP/Port from QR
- Sends pairing requests to the signaling server
- Waits for approval from the desktop app
- Upon approval, establishes WebRTC connection for control

---

## Communication Flow

```
  Mobile App         Signaling Server           Desktop App
      |                    |                        |
      |--- scan QR ------>|                        |
      |--- connect ------>|                        |
      |--- pair request -->                        |
      |                    |--- notify ----------->|
      |                    |<-- approval/reject ---|
      |<--- signaling exchange (SDP, ICE) via ---->|
      |<---------- WebRTC peer-to-peer connection --------->|
```

- The signaling server is used only during connection setup.
- After signaling, the desktop and mobile communicate directly over WebRTC (peer-to-peer).

---

## Session Initiation

1. Rust app starts and launches the signaling server
2. User clicks "Generate QR Code"
3. Rust generates UUID + sessionID and opens the signaling port
4. QR code includes: `sessionID`, `UUID`, optionally `IP:port`
5. Mobile app scans the QR and connects to signaling server
6. Session is created and pairing begins

---

## Session Approval and Connection

- The desktop app shows device info and waits 40 seconds for manual approval
- If approved, signaling messages are exchanged (SDP, ICE)
- After successful WebRTC setup, full control session begins
- If no approval is given within 40 seconds, session is deleted and QR invalidated

---

## Deployment Assumptions

- One session per desktop instance at a time
- Signaling server runs locally and is not publicly hosted
- All sessions and states are held in memory (no persistent storage)
- WebRTC uses STUN for NAT traversal; TURN is not used
- Port forwarding is not required; system relies on hole punching where possible

---

## Notes

- The signaling server is tightly coupled to the desktop app and not intended for reuse across machines
- The architecture prioritizes simplicity, privacy, and local-first control
- Future iterations may support external signaling servers, TURN fallback, or account systems

