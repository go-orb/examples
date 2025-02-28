// Package proto ...
package proto

// go install github.com/go-orb/plugins/server/cmd/protoc-gen-go-orb@latest for --go-orb_out
//go:generate protoc -I . --go-orb_out=paths=source_relative:. --go-orb_opt=supported_servers=drpc ./auth_v1/auth.proto
