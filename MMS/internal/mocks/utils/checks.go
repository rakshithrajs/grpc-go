package mocks

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

// CheckData asserts that actual matches expected after a JSON round-trip.
func CheckData(t *testing.T, actual any, expected any) {
	t.Helper()

	var got any
	switch v := actual.(type) {
	case *httptest.ResponseRecorder:
		body := v.Body.Bytes()
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
	default:
		actualJSON, err := json.Marshal(actual)
		if err != nil {
			t.Fatalf("failed to marshal actual data: %v", err)
		}
		if err := json.Unmarshal(actualJSON, &got); err != nil {
			t.Fatalf("failed to unmarshal actual data: %v", err)
		}
	}

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected data: %v", err)
	}

	var want any
	if err := json.Unmarshal(expectedJSON, &want); err != nil {
		t.Fatalf("failed to unmarshal expected data: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
