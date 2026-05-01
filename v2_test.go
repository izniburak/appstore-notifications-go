package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestServerNotificationV2(t *testing.T) {
	appStoreServerRequest := os.Getenv("APPLE_NOTIFICATION_REQUEST")
	if appStoreServerRequest == "" {
		t.Skip("APPLE_NOTIFICATION_REQUEST not set")
	}
	var request AppStoreServerRequest
	if err := json.Unmarshal([]byte(appStoreServerRequest), &request); err != nil {
		t.Fatal(err)
	}

	rootCert := os.Getenv("APPLE_CERT")
	if rootCert == "" {
		t.Fatal("APPLE_CERT not set")
	}

	asn, err := New(request.SignedPayload, rootCert)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if !asn.IsValid {
		t.Error("Payload is not valid")
	}

	fmt.Printf(
		"NotificationType: %s\nEnvironment: %s\nIsTest: %t\n",
		asn.Payload.NotificationType,
		asn.Payload.Data.Environment,
		asn.IsTest,
	)
}

func TestMalformedPayloads(t *testing.T) {
	cases := []struct {
		name    string
		payload string
	}{
		{"empty string", ""},
		{"not a jwt", "garbage"},
		{"two segments", "a.b"},
		{"three segments no x5c", "eyJhbGciOiJFUzI1NiJ9.e30.sig"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("New panicked: %v", r)
				}
			}()
			asn, err := New(tc.payload, "dummy-cert")
			if err == nil {
				t.Error("expected error, got nil")
			}
			if asn != nil && asn.IsValid {
				t.Error("IsValid should be false on error")
			}
		})
	}
}
