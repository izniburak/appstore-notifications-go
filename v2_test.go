package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestServerNotificationV2(t *testing.T) {
	// {"signedPayload":"..."}
	appStoreServerRequest := os.Getenv("APPLE_NOTIFICATION_REQUEST")
	if appStoreServerRequest == "" {
		panic("No valid AppStoreServerRequest")
	}
	var request AppStoreServerRequest
	err := json.Unmarshal([]byte(appStoreServerRequest), &request) // bind byte to header structure
	if err != nil {
		panic(err)
	}

	// -----BEGIN CERTIFICATE----- ......
	rootCert := os.Getenv("APPLE_CERT")
	if rootCert == "" {
		panic("Apple Root Cert not available")
	}

	appStoreServerNotification := New(request.SignedPayload, rootCert)

	if !appStoreServerNotification.IsValid {
		t.Error("Payload is not valid")
	}

	fmt.Printf(
		"NotificationType: %s\nEnvironment: %s\nIsTest: %t\n",
		appStoreServerNotification.Payload.NotificationType,
		appStoreServerNotification.Payload.Data.Environment,
		appStoreServerNotification.IsTest,
	)
}
