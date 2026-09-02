# Pinoy Henyo

## For testing purposes
## Requirements
#### Bun

To install bun, paste this command in a terminal
`powershell -c "irm bun.sh/install.ps1|iex"`

Then check to see if bun is installed
`bun --version`

#### Go
To install Go, go to their [website](https://go.dev)
Go to downloads page find your Operating System, download installer
PATH should also be updated already no manual updates needed. I think

Check to see if Go is installed
`go version`

## Setup
### Backend
```bash
cd backend
go mod tidy
go run .
# runs on :8080
```

### Frontend
```bash
cd frontend
bun install
bun run dev
# runs on localhost:5173, open in browser
```
