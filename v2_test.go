package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestServerNotificationV2(t *testing.T) {
	// {"signedPayload":"..."}
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("Test error. Could not load .env file: %w", err))
	}
	appStoreServerRequest := os.Getenv("SUBSCRIBED_NOTIFICATION")
	if appStoreServerRequest == "" {
		t.Skip("No valid AppStoreServerRequest")
	}
	var request AppStoreServerRequest
	err = json.Unmarshal([]byte(appStoreServerRequest), &request) // bind byte to header structure
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
	fmt.Printf("Product Id: %s\n", appStoreServerNotification.RenewalInfo.ProductId)
	fmt.Printf("DateSigned: %d\n", appStoreServerNotification.Payload.SignedDate)
	fmt.Printf("IssuedAt: %d\n", appStoreServerNotification.Payload.IssuedAt)
}

func TestServerTestNotificationV2(t *testing.T) {
	// App store test notification issued after requesting a test notification
	// {"signedPayload":"..."}
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("Test error. Could not load .env file: %w", err))
	}
	appStoreServerRequest := os.Getenv("TEST_NOTIFICATION")
	if appStoreServerRequest == "" {
		t.Skip("No valid AppStoreServerTestRequest")
	}
	var request AppStoreServerRequest
	err = json.Unmarshal([]byte(appStoreServerRequest), &request)
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
	fmt.Printf("DateSigned: %d\n", appStoreServerNotification.Payload.SignedDate)
	fmt.Printf("IssuedAt: %d\n", appStoreServerNotification.Payload.IssuedAt)
}
