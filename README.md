# Go HTTP Proxy

This is a simple HTTP proxy written in Go. It was used to develop the DSynth Dashboard locally.

It listens on `localhost` and forwards all requests to a specified target URL. The proxy also handles CORS (Cross-Origin Resource Sharing) requests, making it suitable for local development.

## Features

- Forwards HTTP requests to a specified target URL
- Handles CORS requests
- Configurable target URL via command-line argument

## Requirements

- Go 1.16 or later

## Example

`./dragonflyproxy -target https://ironman.dragonflybsd.org`