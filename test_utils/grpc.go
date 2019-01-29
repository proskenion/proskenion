package test_utils

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func RandomMockServerStream() grpc.ServerStream {
	return &MockServerStream{&MockStream{}}
}

func RandomMockClientStream() grpc.ClientStream {
	return &MockClientStream{&MockStream{}}
}

type MockServerStream struct {
	*MockStream
}

func (*MockServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (*MockServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (*MockServerStream) SetTrailer(metadata.MD) {}

type MockClientStream struct {
	*MockStream
}

func (*MockClientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (*MockClientStream) Trailer() metadata.MD {
	return nil
}

func (*MockClientStream) CloseSend() error {
	return io.EOF
}

type MockStream struct{}

func (s *MockStream) Context() context.Context {
	return context.TODO()
}

func (s *MockStream) SendMsg(m interface{}) error {
	return nil
}

func (s *MockStream) RecvMsg(m interface{}) error {
	return nil
}
