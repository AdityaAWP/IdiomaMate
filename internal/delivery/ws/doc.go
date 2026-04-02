// Package ws provides the real-time WebSocket communication layer for IdiomaMate.
//
// This package is responsible for:
// 1. Handling the WebSocket upgrade handshake from HTTP.
// 2. Managing active client connections (Clients) in a thread-safe Hub.
// 3. Routing incoming JSON messages (Matchmaking, Chat) to the corresponding Services.
// 4. Broadcasting real-time notifications (Join Requests, Approvals, Kicks) back to users.
//
// By isolating this in the delivery layer, the core domain and services remain unaware
// of WebSocket specifics (Cleaner Architecture). The Hub implements the NotificationService
// interface to allow the application to push real-time events without circular dependencies.
package ws
