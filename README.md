# GOAuth

An experimental, self-hostable, OAuth 2 server.

## ⚡️ Requirements

```sh
brew install just
brew install encoredev/tap/encore

go get gotest.tools/gotestsum@latest
go install gotest.tools/gotestsum@latest
go install github.com/go-delve/delve/cmd/dlv@latest

encore auth login

just run
just dev
just test
just test-watch
```

## Debugging

https://encore.dev/docs/go/how-to/debug
