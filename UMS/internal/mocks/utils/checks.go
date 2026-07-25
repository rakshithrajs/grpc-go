package mocks

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ErrorResponse struct {
	Error any `json:"error"`
}

func CheckError(t *testing.T, w *httptest.ResponseRecorder, expected any) {
	t.Helper()

	body := w.Body.Bytes()

	var resp ErrorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	switch want := expected.(type) {
	case string:
		got, ok := resp.Error.(string)
		if !ok {
			t.Fatalf("expected single error %q, got keyed error %v", want, resp.Error)
		}
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	case map[string]string:
		got, ok := resp.Error.(map[string]any)
		if !ok {
			t.Fatalf("expected keyed validation errors, got %v", resp.Error)
		}
		if len(got) != len(want) {
			t.Errorf("expected %d errors, got %d (%v)", len(want), len(got), got)
		}
		for key, wantErr := range want {
			gotVal, exists := got[key]
			if !exists {
				t.Errorf("expected error for field %q", key)
				continue
			}
			gotStr, ok := gotVal.(string)
			if !ok {
				t.Errorf("expected string error for field %q, got %v", key, gotVal)
				continue
			}
			if gotStr != wantErr {
				t.Errorf("expected field %q to have error %q, got %q", key, wantErr, gotStr)
			}
		}
	default:
		t.Fatalf("unsupported expected error type %T", expected)
	}
}

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
