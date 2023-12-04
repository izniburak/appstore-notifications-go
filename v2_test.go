package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestServerNotificationV2(t *testing.T) {
	// {"signedPayload":"..."}
	appStoreServerRequest := os.Getenv("APPLE_NOTIFICATIONNOTIFICATION_REQUEST")
	if appStoreServerRequest == "" {
		t.Skip("No valid AppStoreServerRequest")
	}
	var request AppStoreServerRequest
	err := json.Unmarshal([]byte(appStoreServerRequest), &request) // bind byte to header structure
	if err != nil {
		panic(err)
	}

	// -----BEGIN CERTIFICATE----- ......
	rootCert := os.Getenv("APPLE_CERT")
	if rootCert == "" {
		t.Skip("Apple Root Cert not available")
	}

	appStoreServerNotification := New(request.SignedPayload, rootCert)

	if !appStoreServerNotification.IsValid {
		t.Error("Payload is not valid")
	}

	if appStoreServerNotification.Payload.Data.Environment != "Sandbox" {
		t.Errorf("got %s, want Sandbox", appStoreServerNotification.Payload.Data.Environment)
	}

	println(appStoreServerNotification.Payload.Data.BundleId)
	println(appStoreServerNotification.TransactionInfo.ProductId)
	fmt.Printf("Product Id: %s", appStoreServerNotification.RenewalInfo.ProductId)
}

func TestServerTestNotificationV2(t *testing.T) {
	// App store test notification issued after requesting a test notification
	// {"signedPayload":"..."}
	appStoreServerRequest := os.Getenv("APPLE_NOTIFICATION_TEST_REQUEST")
	if appStoreServerRequest == "" {
		t.Skip("No valid AppStoreServerTestRequest")
	}
	var request AppStoreServerRequest
	err := json.Unmarshal([]byte(appStoreServerRequest), &request)
	if err != nil {
		panic(err)
	}

	// -----BEGIN CERTIFICATE----- ......
	rootCert := os.Getenv("APPLE_CERT")
	if rootCert == "" {
		t.Skip("Apple Root Cert not available")
	}

	appStoreServerNotification := New(request.SignedPayload, rootCert)

	if !appStoreServerNotification.IsValid {
		t.Error("Payload is not valid")
	}

	if appStoreServerNotification.Payload.Data.Environment != "Sandbox" {
		t.Errorf("got %s, want Sandbox", appStoreServerNotification.Payload.Data.Environment)
	}

	println(appStoreServerNotification.Payload.Data.BundleId)
	fmt.Printf("Notification Type: %s\n", appStoreServerNotification.Payload.NotificationType)
	fmt.Printf("Notification Sub-Type: %s\n", appStoreServerNotification.Payload.Subtype)
}
