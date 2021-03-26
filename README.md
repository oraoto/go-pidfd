# pidfd for go

[![Go Reference](https://pkg.go.dev/badge/github.com/oraoto/go-pidfd.svg)](https://pkg.go.dev/github.com/oraoto/go-pidfd)
[![Go Report Card](https://goreportcard.com/badge/github.com/oraoto/go-pidfd)](https://goreportcard.com/report/github.com/oraoto/go-pidfd)

Go bindings to `pidfd_open`, `pidfd_getfd`, `pidfd_send_signal`on Linux 5.6+.

## Example Usages

- [Share listening socket](./examples/share-listen-fd), without fork or sendmsg
- Process management
