//go:build tools
// +build tools

package tools

import (
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/ofabry/go-callvis"
	_ "github.com/yoheimuta/protolint/cmd/protolint"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
