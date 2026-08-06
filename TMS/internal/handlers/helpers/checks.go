package handlers

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckGRPCStatus asserts that err is a gRPC status error with the expected code and message.
func CheckGRPCStatus(t *testing.T, err error, wantCode codes.Code, wantMsg string) {
	t.Helper()

	st, _ := status.FromError(err)

	if st.Code() != wantCode {
		t.Errorf("expected code %v, got %v", wantCode, st.Code())
	}

	if st.Message() != wantMsg {
		t.Errorf("expected message %q, got %q", wantMsg, st.Message())
	}
}

// CheckTokenClaims verifies the expected static claims and that issuedAt is a recent timestamp.
func CheckTokenClaims(t *testing.T, got *TMSpb.TokenClaims, wantUserID, wantIssuer string) {
	t.Helper()

	if got.GetUserID() != wantUserID {
		t.Errorf("expected userID %q, got %q", wantUserID, got.GetUserID())
	}
	if got.GetIssuer() != wantIssuer {
		t.Errorf("expected issuer %q, got %q", wantIssuer, got.GetIssuer())
	}
	if got.GetIssuedAt() == 0 {
		t.Errorf("expected non-zero issuedAt")
	}
	issuedAt := time.Unix(got.GetIssuedAt(), 0)
	if time.Since(issuedAt) > time.Minute {
		t.Errorf("expected issuedAt within the last minute, got %v", issuedAt)
	}
}

// CheckData compares actual and expected values by round-tripping through JSON.
func CheckData(t *testing.T, actual any, expected any) {
	t.Helper()

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("failed to marshal actual data: %v", err)
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected data: %v", err)
	}

	var got, want any
	if err := json.Unmarshal(actualJSON, &got); err != nil {
		t.Fatalf("failed to unmarshal actual data: %v", err)
	}
	if err := json.Unmarshal(expectedJSON, &want); err != nil {
		t.Fatalf("failed to unmarshal expected data: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
