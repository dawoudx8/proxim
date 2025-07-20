# User Flow – Proxim (Internal Planning)

> This user flow describes a typical successful session initiation in Proxim from a desktop user’s perspective. It is intended for internal planning and prototyping purposes only.

---

## 1. Launching the Desktop App

- The user opens the Proxim desktop application.
- Regardless of what is displayed on the screen, the user clicks the **"Generate QR Code"** button.

---

## 2. Choosing a Session Type

- A prompt appears asking the user to select the type of session:
  - **One-time**: Expires immediately after the first disconnection.
  - **Persistent**: Stays active until manually deleted.
  - **Temporary**: Expires after a custom duration.

- If **Temporary** is selected:
  - The user is asked to specify how long the session should remain valid.
  - After this duration, the session will be forcefully invalidated.

---

## 3. QR Code Generation

- Once the session type is selected and confirmed, a QR code is generated and shown on screen.
- This QR code contains all necessary connection metadata.
- The code is valid for 15 minutes or until manually closed.

---

## 4. Scanning the QR Code (Mobile App)

- The user opens the Proxim mobile app and scans the displayed QR code.
- The app immediately sends a pairing request to the desktop application.
- The mobile device does not connect directly, but instead sends the request via the signaling server.

---

## 5. Pairing Request Prompt on Desktop

- The desktop app displays a pairing prompt showing:
  - Device name
  - Device model or identifier

- The prompt gives the user two options:
  - **Approve**
  - **Reject**

- If the user takes no action within **40 seconds**, the prompt disappears and the request is considered rejected.

---

## 6. Final Confirmation and Warning

- If the user approves the connection:
  - A secondary message appears warning that accepting this connection will grant full control of the system to the requesting mobile device.
  - Upon confirming, the session becomes active.

---

## 7. Active Session

- Once paired and approved, the mobile device gains full control over the desktop.
- The mobile user can:
  - Control mouse and keyboard
  - Interact with the operating system
  - Issue system-level commands (e.g., shutdown, restart, sleep)
- The session behavior depends on the originally selected session type.

---

## Notes

- This flow assumes everything works correctly and there are no network or technical errors.
- No account login or authentication is required during this process.
- This document is **not for publishing**, and is meant to capture the ideal user experience path during early design.

