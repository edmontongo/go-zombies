// Command server launches an HTTP server for the zombie simulation API.
//
// Usage
//
// By default the server is launched on port 11235. This can be set with the
// -port flag. Currently only one gameroom is supported, the status of which
// can be viewed at http://localhost:11235/status.
//
// API
//
// The simple HTTP based API is consists of these urls:
//   /register     - Registers a new player, returns their unique id
//   /collision    - Accepts a unique ID and checks if there has been a collision with another player
//
package main
