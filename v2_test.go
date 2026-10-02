package v2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestNewParsesNestedSignedPayloads(t *testing.T) {
	fixture := newTestSigningFixture(t)
	transactionInfo := &TransactionInfo{
		ProductId:            "product",
		AppTransactionId:     "app-transaction",
		OfferPeriod:          "P1M",
		RevocationType:       "REFUND_PRORATED",
		RevocationPercentage: 25000,
		AdvancedCommerceInfo: &AdvancedCommerceTransactionInfo{
			Descriptors: &AdvancedCommerceDescriptors{
				Description: "description",
				DisplayName: "display name",
			},
			EstimatedTax: 25,
			Items: []AdvancedCommerceTransactionItem{{
				SKU:   "sku",
				Price: 1000,
				Offer: &AdvancedCommerceOffer{
					Period:      "P1M",
					PeriodCount: 2,
					Price:       500,
					Reason:      "PROMOTIONAL",
				},
				Refunds: []AdvancedCommerceRefund{{
					RefundAmount: 100,
					RefundDate:   1700000000000,
					RefundReason: "CUSTOMER_REQUEST",
					RefundType:   "FULL",
				}},
			}},
			Period:            "P1M",
			TaxCode:           "tax-code",
			TaxExclusivePrice: 975,
		},
		BillingPlanType: "MONTHLY",
		CommitmentInfo: &TransactionCommitmentInfo{
			BillingPeriodNumber:   2,
			CommitmentExpiresDate: 1700000000000,
			CommitmentPrice:       1200,
			TotalBillingPeriods:   12,
		},
	}
	renewalInfo := &RenewalInfo{
		EligibleWinBackOfferIds: []string{"offer-one", "offer-two"},
		AppAccountToken:         "account-token",
		AppTransactionId:        "app-transaction",
		OfferPeriod:             "P1M",
		AdvancedCommerceInfo: &AdvancedCommerceRenewalInfo{
			ConsistencyToken: "consistency-token",
			Items: []AdvancedCommerceRenewalItem{{
				SKU:   "sku",
				Price: 1100,
				PriceIncreaseInfo: &AdvancedCommercePriceIncreaseInfo{
					DependentSKUs: []string{"dependent-sku"},
					Price:         1100,
					Status:        "SCHEDULED",
				},
			}},
			Period: "P1M",
		},
		CommitmentInfo: &RenewalCommitmentInfo{
			CommitmentAutoRenewProductId:     "product",
			CommitmentAutoRenewStatus:        1,
			CommitmentRenewalBillingPlanType: "MONTHLY",
			CommitmentRenewalDate:            1700000000000,
			CommitmentRenewalPrice:           1200,
		},
		RenewalBillingPlanType: "MONTHLY",
	}

	signedTransaction := fixture.sign(t, transactionInfo, jwt.SigningMethodES256, fixture.leafKey)
	signedRenewal := fixture.sign(t, renewalInfo, jwt.SigningMethodES256, fixture.leafKey)
	payload := &NotificationPayload{
		NotificationType: "DID_RENEW",
		Data: NotificationData{
			AppAppleId:            123456789,
			SignedTransactionInfo: signedTransaction,
			SignedRenewalInfo:     signedRenewal,
		},
	}
	signedNotification := fixture.sign(t, payload, jwt.SigningMethodES256, fixture.leafKey)

	asn, err := New(signedNotification, fixture.rootPEM)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if !asn.IsValid || asn.IsTest {
		t.Fatalf("unexpected notification state: IsValid=%t IsTest=%t", asn.IsValid, asn.IsTest)
	}
	if asn.TransactionInfo == nil || asn.TransactionInfo.ProductId != "product" ||
		asn.TransactionInfo.AppTransactionId != "app-transaction" ||
		asn.TransactionInfo.RevocationPercentage != 25000 ||
		asn.TransactionInfo.AdvancedCommerceInfo.Items[0].Refunds[0].RefundAmount != 100 ||
		asn.TransactionInfo.CommitmentInfo.TotalBillingPeriods != 12 {
		t.Fatalf("transaction fields were not parsed: %#v", asn.TransactionInfo)
	}
	if asn.RenewalInfo == nil || len(asn.RenewalInfo.EligibleWinBackOfferIds) != 2 ||
		asn.RenewalInfo.EligibleWinBackOfferIds[1] != "offer-two" ||
		asn.RenewalInfo.AppAccountToken != "account-token" ||
		asn.RenewalInfo.AdvancedCommerceInfo.Items[0].PriceIncreaseInfo.Status != "SCHEDULED" ||
		asn.RenewalInfo.CommitmentInfo.CommitmentRenewalPrice != 1200 {
		t.Fatalf("renewal fields were not parsed: %#v", asn.RenewalInfo)
	}
}

func TestNewParsesNotificationVariants(t *testing.T) {
	fixture := newTestSigningFixture(t)
	appTransaction := &AppTransactionInfo{
		AppAppleId:                 123456789,
		BundleId:                   "com.example.app",
		ApplicationVersion:         "1.0",
		OriginalApplicationVersion: "1.0",
		AppTransactionId:           "app-transaction",
	}
	signedAppTransaction := fixture.sign(t, appTransaction, jwt.SigningMethodES256, fixture.leafKey)

	tests := []struct {
		name    string
		payload *NotificationPayload
		check   func(*testing.T, *AppStoreServerNotification)
	}{
		{
			name: "test",
			payload: &NotificationPayload{
				NotificationType: "TEST",
			},
			check: func(t *testing.T, asn *AppStoreServerNotification) {
				if !asn.IsTest || asn.TransactionInfo != nil || asn.RenewalInfo != nil {
					t.Fatalf("unexpected test notification state: %#v", asn)
				}
			},
		},
		{
			name: "external purchase token",
			payload: &NotificationPayload{
				NotificationType: "EXTERNAL_PURCHASE_TOKEN",
				Subtype:          "CREATED",
				ExternalPurchaseToken: ExternalPurchaseToken{
					ExternalPurchaseId:  "external-id",
					TokenCreationDate:   1700000000000,
					AppAppleId:          123456789,
					BundleId:            "com.example.app",
					TokenType:           "ACQUISITION",
					TokenExpirationDate: 1700003600000,
				},
			},
			check: func(t *testing.T, asn *AppStoreServerNotification) {
				token := asn.Payload.ExternalPurchaseToken
				if token.TokenType != "ACQUISITION" || token.TokenExpirationDate != 1700003600000 {
					t.Fatalf("external purchase token fields were not parsed: %#v", token)
				}
			},
		},
		{
			name: "rescind consent",
			payload: &NotificationPayload{
				NotificationType: "RESCIND_CONSENT",
				AppData: &AppData{
					AppAppleId:               123456789,
					BundleId:                 "com.example.app",
					Environment:              "Production",
					SignedAppTransactionInfo: signedAppTransaction,
				},
			},
			check: func(t *testing.T, asn *AppStoreServerNotification) {
				if asn.Payload.AppData == nil || asn.AppTransactionInfo == nil ||
					asn.AppTransactionInfo.AppTransactionId != "app-transaction" {
					t.Fatalf("app data was not parsed: %#v", asn)
				}
			},
		},
		{
			name: "summary",
			payload: &NotificationPayload{
				NotificationType: "RENEWAL_EXTENSION",
				Summary: NotificationSummary{
					AppAppleId:             123456789,
					StoreFrontCountryCodes: []string{"USA", "CAN"},
					FailedCount:            2,
					SucceededCount:         10,
				},
			},
			check: func(t *testing.T, asn *AppStoreServerNotification) {
				if asn.Payload.Summary.AppAppleId != 123456789 ||
					asn.Payload.Summary.SucceededCount != 10 ||
					len(asn.Payload.Summary.StoreFrontCountryCodes) != 2 {
					t.Fatalf("summary fields were not parsed: %#v", asn.Payload.Summary)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			signedNotification := fixture.sign(t, tc.payload, jwt.SigningMethodES256, fixture.leafKey)
			asn, err := New(signedNotification, fixture.rootPEM)
			if err != nil {
				t.Fatalf("New returned error: %v", err)
			}
			if !asn.IsValid {
				t.Fatal("notification was not marked valid")
			}
			tc.check(t, asn)
		})
	}
}

func TestNewRejectsInvalidSignedPayloads(t *testing.T) {
	fixture := newTestSigningFixture(t)
	transaction := &TransactionInfo{ProductId: "product"}
	validTransaction := fixture.sign(t, transaction, jwt.SigningMethodES256, fixture.leafKey)
	validPayload := &NotificationPayload{
		NotificationType: "DID_RENEW",
		Data: NotificationData{
			SignedTransactionInfo: validTransaction,
		},
	}

	otherFixture := newTestSigningFixture(t)
	untrustedTransaction := otherFixture.sign(t, transaction, jwt.SigningMethodES256, otherFixture.leafKey)
	untrustedPayload := *validPayload
	untrustedPayload.Data.SignedTransactionInfo = untrustedTransaction

	wrongAlgorithm := fixture.sign(t, &NotificationPayload{NotificationType: "TEST"}, jwt.SigningMethodHS256, []byte("secret"))
	badSignature := fixture.sign(t, &NotificationPayload{NotificationType: "TEST"}, jwt.SigningMethodES256, newTestKey(t))

	tests := []struct {
		name    string
		payload string
		root    string
	}{
		{
			name:    "untrusted root",
			payload: fixture.sign(t, &NotificationPayload{NotificationType: "TEST"}, jwt.SigningMethodES256, fixture.leafKey),
			root:    otherFixture.rootPEM,
		},
		{
			name:    "untrusted nested certificate",
			payload: fixture.sign(t, &untrustedPayload, jwt.SigningMethodES256, fixture.leafKey),
			root:    fixture.rootPEM,
		},
		{
			name:    "unsupported algorithm",
			payload: wrongAlgorithm,
			root:    fixture.rootPEM,
		},
		{
			name:    "invalid signature",
			payload: badSignature,
			root:    fixture.rootPEM,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			asn, err := New(tc.payload, tc.root)
			if err == nil {
				t.Fatal("expected New to reject the payload")
			}
			if asn != nil {
				t.Fatalf("expected nil notification on error, got %#v", asn)
			}
		})
	}
}

func TestMalformedPayloads(t *testing.T) {
	cases := []struct {
		name    string
		payload string
	}{
		{"empty string", ""},
		{"not a jwt", "garbage"},
		{"two segments", "a.b"},
		{"invalid header", "******.e30.c2ln"},
		{"header without x5c", "e30.e30.c2ln"},
		{"too many segments", "e30.e30.c2ln.extra"},
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

type testSigningFixture struct {
	rootPEM   string
	leafKey   *ecdsa.PrivateKey
	certChain []string
}

func newTestSigningFixture(t *testing.T) *testSigningFixture {
	t.Helper()
	now := time.Now()

	rootKey := newTestKey(t)
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test Root"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("create test root certificate: %v", err)
	}

	intermediateKey := newTestKey(t)
	intermediateTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Test Intermediate"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	intermediateDER, err := x509.CreateCertificate(rand.Reader, intermediateTemplate, rootTemplate, &intermediateKey.PublicKey, rootKey)
	if err != nil {
		t.Fatalf("create test intermediate certificate: %v", err)
	}

	leafKey := newTestKey(t)
	leafTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(3),
		Subject:               pkix.Name{CommonName: "Test Leaf"},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour),
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, intermediateTemplate, &leafKey.PublicKey, intermediateKey)
	if err != nil {
		t.Fatalf("create test leaf certificate: %v", err)
	}

	rootCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootDER})
	return &testSigningFixture{
		rootPEM: string(rootCert),
		leafKey: leafKey,
		certChain: []string{
			base64.StdEncoding.EncodeToString(leafDER),
			base64.StdEncoding.EncodeToString(intermediateDER),
			base64.StdEncoding.EncodeToString(rootDER),
		},
	}
}

func (fixture *testSigningFixture) sign(t *testing.T, claims jwt.Claims, method jwt.SigningMethod, key interface{}) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	token.Header["x5c"] = fixture.certChain
	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign test JWT: %v", err)
	}
	return signed
}

func newTestKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate test key: %v", err)
	}
	return key
}
