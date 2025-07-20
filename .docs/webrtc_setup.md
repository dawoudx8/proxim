# WebRTC & STUN Setup – Proxim

This document explains how to configure and use WebRTC with STUN for establishing a peer-to-peer connection between the Proxim mobile and desktop clients.

---

## Overview

WebRTC is used in Proxim to create a direct data channel between the mobile device and the desktop app. Because most devices are behind NAT (home routers or firewalls), we use **STUN** servers to help both peers discover their public-facing IP and port.

* **STUN** = Session Traversal Utilities for NAT
* **TURN** is intentionally not used due to budget and privacy reasons
* Proxim relies on WebRTC + STUN + hole punching

---

## Basic WebRTC Flow (Simplified)

```
Mobile App              Signaling Server             Desktop App
    |                         |                           |
    |--- pair_request ------> |                           |
    |                         |--- pair_prompt ---------> |
    |<---- pair_result ------ |<-- pair_response -------- |
    |                         |                           |
    |-- SDP Offer (signal) -> |-> SDP Offer ------------> |
    |<- SDP Answer (signal) --|<- SDP Answer ------------ |
    |-- ICE Candidates -----> |-> ICE Candidates -------->|
    |<-- ICE Candidates ------|<- ICE Candidates -------- |
    |                         |                           |
    |<-------- WebRTC Peer-to-Peer Connection --------->  |
```

---

## STUN Configuration

STUN servers must be provided to the WebRTC engine on both sides (Rust and Android). A public free STUN server (such as Google's) is sufficient for most cases.

### Recommended STUN server:

```text
stun:stun.l.google.com:19302
```

You can optionally add more for redundancy:

```text
stun:stun1.l.google.com:19302
stun:stun2.l.google.com:19302
```

---

## Desktop (Rust) – WebRTC Setup (Conceptual)

You will likely use a crate like [`webrtc`](https://github.com/webrtc-rs/webrtc) for the Rust side.

### Sample ICE config:

```rust
let config = RTCConfiguration {
    ice_servers: vec![
        RTCIceServer {
            urls: vec!["stun:stun.l.google.com:19302".to_string()],
            ..Default::default()
        }
    ],
    ..Default::default()
};
```

* Pass this `config` when creating the peer connection.
* Use async handlers to process incoming SDP and ICE messages.

---

## Mobile (Android Kotlin) – WebRTC Setup

WebRTC for Android is provided via the `org.webrtc` library (used in Google Meet, etc).

### Initialize PeerConnectionFactory:

```kotlin
val options = PeerConnectionFactory.InitializationOptions.builder(context)
    .createInitializationOptions()
PeerConnectionFactory.initialize(options)
```

### Create STUN configuration:

```kotlin
val iceServers = listOf(
    PeerConnection.IceServer.builder("stun:stun.l.google.com:19302").createIceServer()
)

val rtcConfig = PeerConnection.RTCConfiguration(iceServers)
```

* Use `rtcConfig` when creating the peer connection.
* Handle callbacks for `onIceCandidate`, `onAddStream`, etc.

---

## Summary

* STUN allows peer-to-peer communication across NAT/firewalls.
* WebRTC needs signaling first (handled by your Go server).
* After signaling, peers exchange SDP and ICE to connect directly.
* No TURN server is used. If NAT traversal fails, connection will fail.

This setup should work well for most home/office networks.
For rare edge cases (symmetric NAT), connection may fail, but this is acceptable in MVP stage.
