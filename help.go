package main

const Help = `Ping-Pong (v%v) | tpkn.me

Running mock server on user defined port. Returns back all requests with timestamp and payload.

Usage:
  ping_pong -p <port>

Options:
  -p, --port      Server port (default: 8181")
  -j, --json      Pong in JSON format (default: tsv)
  -s, --silent    Don't print requests to STDOUT
  --cpu           Maximum number of CPU cores used by server
  -h, --help      Help
  -v, --version   Version
`
