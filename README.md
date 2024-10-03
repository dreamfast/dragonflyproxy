# Go HTTP Proxy

This is a simple HTTP proxy written in Go. It was originally used to develop the DSynth Dashboard locally, but it can be used for various local development scenarios.

It listens on a specified port (default: 8899) and forwards all requests to a specified target URL. The proxy also handles CORS (Cross-Origin Resource Sharing) requests, making it suitable for local development.

## Features

- Forwards HTTP requests to a specified target URL
- Handles CORS requests
- Configurable target URL and port via command-line arguments
- Sets X-Forwarded-For header for proper client IP forwarding

## Requirements

- Go 1.16 or later

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/dreamfast/dragonflyproxy.git
   cd dragonflyproxy
   ```

2. Build the binary:
   ```
   go build
   ```

## Usage

Run the proxy with default settings:

```
./dragonflyproxy
```

Specify a custom target URL and port:

```
./dragonflyproxy -target https://example.com -port 8080
```

### Command-line Options

- `-target`: The target URL to forward requests to (default: https://ironman.dragonflybsd.org)
- `-port`: The port to run the proxy server on (default: 8899)

## Example

```
./dragonflyproxy -target https://api.example.com -port 3000
```

This will start the proxy server on `http://localhost:3000`, forwarding all requests to `https://api.example.com`.

## License

This project is open source and available under the [MIT License](LICENSE).