package main

type Command struct {
	body []byte
}

//we need to make this flow for serialization
// client receives the wire-protocol message,
//  parses it,
//  that the client sends to the server.
