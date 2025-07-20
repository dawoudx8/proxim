# Development Tasks – Proxim

This document outlines the development phases and tasks required to implement the Proxim system in a structured and incremental way.

---

> **Note:** Go handles the signaling server, Rust powers the desktop client, and Kotlin is used for the Android app.

---

## Phase 1: Core Infrastructure (Signaling & QR)
> **Info:** Tasks 1–5 should be implemented using Go, and tasks 6-9 using Rust.

### Task 1: Create basic Go signaling server
- [x] Set up a basic HTTP server that upgrades to WebSocket 
- [x] Accept incoming WebSocket connections and assign client roles (mobile or desktop)
- [ ] Build a message routing system based on sessionID and role 
- [ ] Define basic JSON input/output handling with decoding and validation 

### Task 2: Define signaling message types and schema 
- [ ] Create Go structs for each message type: `pair_request`, `pair_prompt`, `pair_response`, `pair_result`, `signal`
- [ ] Implement JSON unmarshalling with validation for each struct 
- [ ] Add unit tests for message validation and structure enforcement 
- [ ] Implement error responses for unknown or invalid message types 

### Task 3: Implement in-memory session store 
- [ ] Define a `Session` struct: includes sessionID, UUID, createdAt, expiresAt, status, deviceInfo 
- [ ] Use a thread-safe `map[string]Session` with `sync.Mutex` or `sync.Map` 
- [ ] Implement `createSession(sessionID, UUID, expiresAt)` 
- [ ] Implement `getSession(sessionID)`, `updateSession(sessionID, newState)`, `deleteSession(sessionID)` 
- [ ] Build a background task to periodically clean expired sessions 

### Task 4: Handle pairing requests 
- [ ] On receiving `pair_request`, check if sessionID exists and is in `pending` state 
- [ ] Attach `deviceInfo` to the session 
- [ ] Forward a `pair_prompt` message to the appropriate desktop connection 
- [ ] Store timestamp of the pairing attempt 
- [ ] Enforce single active request per sessionID 

### Task 5: Auto-expire session and pairing 
- [ ] Implement a 40-second pairing request timeout timer per session
- [ ] Mark sessions as `rejected` if no `pair_response` is received 
- [ ] Mark sessions as `expired` after global `expiresAt` timestamp passes 
- [ ] Remove fully expired sessions from memory in cleanup loop 
- [ ] Log expiration events for debugging 

### Task 6: Create minimal Rust desktop app stub 
- [ ] Initialize new Rust project with `cargo new proxim-desktop` 
- [ ] Add basic CLI using `clap` or minimal window using `egui` or `fltk` 
- [ ] Print debug logs to console when app launches 
- [ ] Display placeholder text or simple UI to indicate running state 

### Task 7: Generate UUID and sessionID 
- [ ] On first launch, check for existing UUID file, otherwise generate new UUID using `uuid` crate 
- [ ] Save UUID to disk (e.g., `~/.proxim/uuid`) 
- [ ] Generate a random sessionID on each pairing using alphanumeric string 
- [ ] Associate sessionID and UUID to new session request 
- [ ] Prepare struct to send to signaling server 

### Task 8: Launch signaling server as subprocess 
- [ ] Include `proxim-server.exe` in project directory or embed into build artifacts 
- [ ] Use `std::process::Command` to start signaling server from Rust 
- [ ] Capture and log stdout/stderr from subprocess 
- [ ] Handle automatic restart if subprocess crashes unexpectedly 
- [ ] Kill process gracefully on app exit 

### Task 9: Display QR code 
- [ ] Create QR content object: includes UUID, sessionID, and optional IP:port 
- [ ] Serialize content as JSON string 
- [ ] Use `qrcode` crate to generate QR bitmap 
- [ ] Render in terminal (for CLI) or in image widget (for GUI) 
- [ ] Handle QR expiration or regeneration when session is removed 

---

## Phase 2: Mobile Pairing & Approval Flow
> **Info:** Tasks 1–4 should be implemented using Kotlin, tasks 5-6 using Rust, and task 7 using Go & Kotlin.

### Task 1: Build Kotlin mobile app skeleton 
- [ ] Create Android Studio project with minimum SDK 24+ 
- [ ] Set up single Activity and navigation controller (if needed) 
- [ ] Add required permissions for camera and internet 
- [ ] Test project builds and launches on emulator/device

### Task 2: Implement QR code scanner 
- [ ] Add camera preview using CameraX or ZXing 
- [ ] Scan and decode QR code into structured data 
- [ ] Parse sessionID, UUID, and optional IP:port from scanned result 
- [ ] Show confirmation screen after scanning

### Task 3: Connect to signaling server from mobile
- [ ] Establish WebSocket connection to scanned IP/port 
- [ ] Handle connection success/failure with UI feedback 
- [ ] Keep connection alive during pairing attempt 

### Task 4: Send `pair_request` message 
- [ ] Construct JSON message with type `pair_request` 
- [ ] Include sessionID and deviceInfo (model, name, ID) 
- [ ] Send over WebSocket and await response 

### Task 5: Handle pairing prompt in Rust app 
- [ ] Listen for `pair_prompt` message from signaling server
- [ ] Display prompt UI with device name + identifier 
- [ ] Start 40-second countdown timer 
- [ ] Show Approve / Reject buttons to user 

### Task 6: Auto-expire or manually approve 
- [ ] On Approve → send `pair_response` with `approved = true` 
- [ ] On Reject → send `pair_response` with `approved = false` 
- [ ] If 40 seconds pass → mark session as `expired`, no message sent 

### Task 7: Notify mobile app of result 
- [ ] Server forwards `pair_result` to mobile client
- [ ] If approved, mobile transitions to next state (WebRTC negotiation) 
- [ ] If rejected or expired, show error and option to rescan 

---

## Phase 3: WebRTC Channel Establishment
> **Note:** WebRTC and STUN setup is entirely handled by the desktop (Rust) and mobile (Android) apps. The signaling server (Go) only relays messages and does not interact with STUN or WebRTC directly.
> **Info:** Task 1 should be implemented using Rust, task 2 using Kotlin, task 3 using Go, and task 4 using Rust & Kotlin.

### Task 1: Integrate WebRTC in Rust 
- [ ] Choose and configure `webrtc-rs` or equivalent crate 
- [ ] Set up SDP and ICE handling
- [ ] Expose signaling handlers (send/receive over WebSocket) 

### Task 2: Integrate WebRTC in Android 
- [ ] Add WebRTC dependency via gradle 
- [ ] Create peer connection and data channel setup logic 
- [ ] Handle ICE candidates and connection state 

### Task 3: Implement signaling relay via server 
- [ ] Forward `signal` type messages between peers via signaling server 
- [ ] Include payload for SDP and ICE 
- [ ] Ensure message order and delivery correctness 

### Task 4: Establish working P2P connection 
- [ ] Confirm that peer connection forms over LAN 
- [ ] Test NAT traversal with public STUN servers 
- [ ] Log ICE connection status on both ends 

---

## Phase 4: Remote Control Functionality
> **Info:** Tasks 1-4 should be implemented using Kotlin & Rust, and task 5 using Rust.

### Task 1: Handle mouse input from mobile
- [ ] Capture touch/motion events on Android 
- [ ] Translate to relative or absolute cursor position 
- [ ] Send as WebRTC data message to desktop
- [ ] Move system cursor using Rust on Windows APIs 

### Task 2: Handle keyboard input from mobile 
- [ ] Capture soft keyboard input from Android 
- [ ] Package key codes and modifiers 
- [ ] Send as WebRTC message 
- [ ] Simulate key press in Rust using native libraries 

### Task 3: System-level commands 
- [ ] Add shutdown, restart, sleep buttons to mobile UI 
- [ ] Send command via WebRTC 
- [ ] Execute corresponding system command on desktop (with confirmation) 

### Task 4: Display connection status
- [ ] Show active/disconnected status in mobile and desktop 
- [ ] Handle reconnection or unexpected disconnects

### Task 5: Allow session termination 
- [ ] Add "End Session" button in desktop UI 
- [ ] On click: close WebRTC, mark session inactive, notify peer 

---

## Phase 5: Stability & Polishing
> **Info:** Tasks 1–2 should be implemented Go, Rust & Kotlin, task 3 using Rust, task 4 using Rust & build tools, and task 5 using Kotlin & Rust.

### Task 1: Improve error handling 
- [ ] Add retries and backoff for WebSocket failures 
- [ ] Validate message fields rigorously 

### Task 2: Add optional session logs 
- [ ] Track session creation, approval, expiration events 
- [ ] Store temporarily in memory or optional local file 

### Task 3: Optimize startup/shutdown
- [ ] Reduce launch time for server subprocess 
- [ ] Graceful shutdown for all components 

### Task 4: Package and distribute 
- [ ] Bundle desktop app and signaling server 
- [ ] Create installer or archive with versioning 

### Task 5: UI design polish
- [ ] Unify mobile and desktop UI style 
- [ ] Add loading indicators, confirmation dialogs, and error prompts 

---

## Notes

- Phases are ordered logically but may overlap
- Each major step should be followed by Git commit (with message like: `feat: Complete X`)
- Minimum viable product (MVP) target = Phase 1 + Phase 2 + basic WebRTC link

