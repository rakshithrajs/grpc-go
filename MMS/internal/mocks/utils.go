package mocks

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CheckGRPCError(t *testing.T, err error, wantCode codes.Code, wantMsg string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected status error, got %v", err)
	}

	if st.Code() != wantCode {
		t.Errorf("expected code %v, got %v", wantCode, st.Code())
	}

	if st.Message() != wantMsg {
		t.Errorf("expected message %q, got %q", wantMsg, st.Message())
	}
}
