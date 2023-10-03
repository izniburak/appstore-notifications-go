package v2

import "github.com/golang-jwt/jwt"

type AppStoreServerNotification struct {
	appleRootCert   string
	Payload         *NotificationPayload
	TransactionInfo *TransactionInfo
	RenewalInfo     *RenewalInfo
	isValid         bool
}

type AppStoreServerRequest struct {
	SignedPayload string `json:"signedPayload"`
}

type NotificationHeader struct {
	Alg string   `json:"alg"`
	X5c []string `json:"x5c"`
}

type NotificationPayload struct {
	jwt.StandardClaims
	NotificationType    string           `json:"notificationType"`
	Subtype             string           `json:"subtype"`
	NotificationUUID    string           `json:"notificationUUID"`
	NotificationVersion string           `json:"notificationVersion"`
	Data                NotificationData `json:"data"`
}

type NotificationData struct {
	AppAppleID            int    `json:"appAppleId"`
	BundleID              string `json:"bundleId"`
	BundleVersion         string `json:"bundleVersion"`
	Environment           string `json:"environment"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

type TransactionInfo struct {
	jwt.StandardClaims
	TransactionId               string `json:"transactionId"`
	OriginalTransactionID       string `json:"originalTransactionId"`
	WebOrderLineItemID          string `json:"webOrderLineItemId"`
	BundleID                    string `json:"bundleId"`
	ProductID                   string `json:"productId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	PurchaseDate                int    `json:"purchaseDate"`
	OriginalPurchaseDate        int    `json:"originalPurchaseDate"`
	ExpiresDate                 int    `json:"expiresDate"`
	Type                        string `json:"type"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	SignedDate                  int    `json:"signedDate"`
	Environment                 string `json:"environment"`
}

type RenewalInfo struct {
	jwt.StandardClaims
	OriginalTransactionID  string `json:"originalTransactionId"`
	ExpirationIntent       int    `json:"expirationIntent"`
	AutoRenewProductId     string `json:"autoRenewProductId"`
	ProductID              string `json:"productId"`
	AutoRenewStatus        int    `json:"autoRenewStatus"`
	IsInBillingRetryPeriod bool   `json:"isInBillingRetryPeriod"`
	SignedDate             int    `json:"signedDate"`
	Environment            string `json:"environment"`
}
