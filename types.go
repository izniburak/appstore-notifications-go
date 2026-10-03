package v2

import "github.com/golang-jwt/jwt/v5"

type AppStoreServerNotification struct {
	appleRootCert      string
	Payload            *NotificationPayload
	TransactionInfo    *TransactionInfo
	RenewalInfo        *RenewalInfo
	AppTransactionInfo *AppTransactionInfo
	IsValid            bool
	IsTest             bool
}

type AppStoreServerRequest struct {
	SignedPayload string `json:"signedPayload"`
}

type NotificationHeader struct {
	Alg string   `json:"alg"`
	X5c []string `json:"x5c"`
}

type NotificationPayload struct {
	jwt.RegisteredClaims
	NotificationType      string                `json:"notificationType"`
	Subtype               string                `json:"subtype"`
	NotificationUUID      string                `json:"notificationUUID"`
	Version               string                `json:"version"`
	SignedDate            int64                 `json:"signedDate"`
	Summary               NotificationSummary   `json:"summary,omitempty"`
	Data                  NotificationData      `json:"data,omitempty"`
	ExternalPurchaseToken ExternalPurchaseToken `json:"externalPurchaseToken,omitempty"`
	AppData               *AppData              `json:"appData,omitempty"`
}

type ExternalPurchaseToken struct {
	ExternalPurchaseId  string `json:"externalPurchaseId"`
	TokenCreationDate   int64  `json:"tokenCreationDate"`
	AppAppleId          int64  `json:"appAppleId"`
	BundleId            string `json:"bundleId"`
	TokenType           string `json:"tokenType"`
	TokenExpirationDate int64  `json:"tokenExpirationDate"`
}

type AppData struct {
	AppAppleId               int64  `json:"appAppleId"`
	BundleId                 string `json:"bundleId"`
	Environment              string `json:"environment"`
	SignedAppTransactionInfo string `json:"signedAppTransactionInfo"`
}

type NotificationSummary struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	AppAppleId             int64    `json:"appAppleId"`
	BundleId               string   `json:"bundleId"`
	ProductId              string   `json:"productId"`
	Environment            string   `json:"environment"`
	StoreFrontCountryCodes []string `json:"storefrontCountryCodes"`
	FailedCount            int64    `json:"failedCount"`
	SucceededCount         int64    `json:"succeededCount"`
}

type NotificationData struct {
	AppAppleId               int64  `json:"appAppleId"`
	BundleId                 string `json:"bundleId"`
	BundleVersion            string `json:"bundleVersion"`
	Environment              string `json:"environment"`
	SignedRenewalInfo        string `json:"signedRenewalInfo"`
	SignedTransactionInfo    string `json:"signedTransactionInfo"`
	Status                   int32  `json:"status"`
	ConsumptionRequestReason string `json:"consumptionRequestReason,omitempty"`
}

type TransactionInfo struct {
	jwt.RegisteredClaims
	AppAccountToken             string                           `json:"appAccountToken"`
	BundleId                    string                           `json:"bundleId"`
	Currency                    string                           `json:"currency,omitempty"`
	Environment                 string                           `json:"environment"`
	ExpiresDate                 int64                            `json:"expiresDate"`
	InAppOwnershipType          string                           `json:"inAppOwnershipType"`
	IsUpgraded                  bool                             `json:"isUpgraded"`
	OfferDiscountType           string                           `json:"offerDiscountType"`
	OfferIdentifier             string                           `json:"offerIdentifier"`
	OfferType                   int32                            `json:"offerType"`
	OriginalPurchaseDate        int64                            `json:"originalPurchaseDate"`
	OriginalTransactionId       string                           `json:"originalTransactionId"`
	Price                       int64                            `json:"price,omitempty"`
	ProductId                   string                           `json:"productId"`
	PurchaseDate                int64                            `json:"purchaseDate"`
	Quantity                    int32                            `json:"quantity"`
	RevocationDate              int64                            `json:"revocationDate"`
	RevocationReason            int32                            `json:"revocationReason"`
	SignedDate                  int64                            `json:"signedDate"`
	StoreFront                  string                           `json:"storefront"`
	StoreFrontId                string                           `json:"storefrontId"`
	SubscriptionGroupIdentifier string                           `json:"subscriptionGroupIdentifier"`
	TransactionId               string                           `json:"transactionId"`
	TransactionReason           string                           `json:"transactionReason"`
	Type                        string                           `json:"type"`
	WebOrderLineItemId          string                           `json:"webOrderLineItemId"`
	AppTransactionId            string                           `json:"appTransactionId"`
	OfferPeriod                 string                           `json:"offerPeriod"`
	RevocationType              string                           `json:"revocationType"`
	RevocationPercentage        int64                            `json:"revocationPercentage"`
	AdvancedCommerceInfo        *AdvancedCommerceTransactionInfo `json:"advancedCommerceInfo"`
	BillingPlanType             string                           `json:"billingPlanType"`
	CommitmentInfo              *TransactionCommitmentInfo       `json:"commitmentInfo"`
}

type RenewalInfo struct {
	jwt.RegisteredClaims
	AutoRenewProductId          string                       `json:"autoRenewProductId"`
	AutoRenewStatus             int32                        `json:"autoRenewStatus"`
	Environment                 string                       `json:"environment"`
	Currency                    string                       `json:"currency"`
	EligibleWinBackOfferIds     []string                     `json:"eligibleWinBackOfferIds"`
	ExpirationIntent            int32                        `json:"expirationIntent"`
	GracePeriodExpiresDate      int64                        `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      bool                         `json:"isInBillingRetryPeriod"`
	OfferDiscountType           string                       `json:"offerDiscountType"`
	OfferIdentifier             string                       `json:"offerIdentifier"`
	OfferType                   int32                        `json:"offerType"`
	OriginalTransactionId       string                       `json:"originalTransactionId"`
	PriceIncreaseStatus         int32                        `json:"priceIncreaseStatus"`
	ProductId                   string                       `json:"productId"`
	RecentSubscriptionStartDate int64                        `json:"recentSubscriptionStartDate"`
	RenewalDate                 int64                        `json:"renewalDate"`
	RenewalPrice                int64                        `json:"renewalPrice"`
	SignedDate                  int64                        `json:"signedDate"`
	AppAccountToken             string                       `json:"appAccountToken"`
	AppTransactionId            string                       `json:"appTransactionId"`
	OfferPeriod                 string                       `json:"offerPeriod"`
	AdvancedCommerceInfo        *AdvancedCommerceRenewalInfo `json:"advancedCommerceInfo"`
	CommitmentInfo              *RenewalCommitmentInfo       `json:"commitmentInfo"`
	RenewalBillingPlanType      string                       `json:"renewalBillingPlanType"`
}

type AppTransactionInfo struct {
	jwt.RegisteredClaims
	ReceiptType                string `json:"receiptType"`
	AppAppleId                 int64  `json:"appAppleId"`
	BundleId                   string `json:"bundleId"`
	ApplicationVersion         string `json:"applicationVersion"`
	VersionExternalIdentifier  int64  `json:"versionExternalIdentifier"`
	ReceiptCreationDate        int64  `json:"receiptCreationDate"`
	OriginalPurchaseDate       int64  `json:"originalPurchaseDate"`
	OriginalApplicationVersion string `json:"originalApplicationVersion"`
	DeviceVerification         string `json:"deviceVerification"`
	DeviceVerificationNonce    string `json:"deviceVerificationNonce"`
	PreorderDate               int64  `json:"preorderDate"`
	AppTransactionId           string `json:"appTransactionId"`
	OriginalPlatform           string `json:"originalPlatform"`
	StoreType                  string `json:"storeType"`
}

type TransactionCommitmentInfo struct {
	BillingPeriodNumber   int32 `json:"billingPeriodNumber"`
	CommitmentExpiresDate int64 `json:"commitmentExpiresDate"`
	CommitmentPrice       int64 `json:"commitmentPrice"`
	TotalBillingPeriods   int32 `json:"totalBillingPeriods"`
}

type RenewalCommitmentInfo struct {
	CommitmentAutoRenewProductId     string `json:"commitmentAutoRenewProductId"`
	CommitmentAutoRenewStatus        int32  `json:"commitmentAutoRenewStatus"`
	CommitmentRenewalBillingPlanType string `json:"commitmentRenewalBillingPlanType"`
	CommitmentRenewalDate            int64  `json:"commitmentRenewalDate"`
	CommitmentRenewalPrice           int64  `json:"commitmentRenewalPrice"`
}

type AdvancedCommerceDescriptors struct {
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
}

type AdvancedCommerceOffer struct {
	Period      string `json:"period"`
	PeriodCount int32  `json:"periodCount"`
	Price       int64  `json:"price"`
	Reason      string `json:"reason"`
}

type AdvancedCommerceRefund struct {
	RefundAmount int64  `json:"refundAmount"`
	RefundDate   int64  `json:"refundDate"`
	RefundReason string `json:"refundReason"`
	RefundType   string `json:"refundType"`
}

type AdvancedCommerceTransactionItem struct {
	SKU            string                   `json:"SKU"`
	Description    string                   `json:"description"`
	DisplayName    string                   `json:"displayName"`
	Offer          *AdvancedCommerceOffer   `json:"offer"`
	Price          int64                    `json:"price"`
	Refunds        []AdvancedCommerceRefund `json:"refunds"`
	RevocationDate int64                    `json:"revocationDate"`
}

type AdvancedCommerceTransactionInfo struct {
	Descriptors        *AdvancedCommerceDescriptors      `json:"descriptors"`
	EstimatedTax       int64                             `json:"estimatedTax"`
	Items              []AdvancedCommerceTransactionItem `json:"items"`
	Period             string                            `json:"period"`
	RequestReferenceId string                            `json:"requestReferenceId"`
	TaxCode            string                            `json:"taxCode"`
	TaxExclusivePrice  int64                             `json:"taxExclusivePrice"`
	TaxRate            string                            `json:"taxRate"`
}

type AdvancedCommercePriceIncreaseInfo struct {
	DependentSKUs []string `json:"dependentSKUs"`
	Price         int64    `json:"price"`
	Status        string   `json:"status"`
}

type AdvancedCommerceRenewalItem struct {
	SKU               string                             `json:"SKU"`
	Description       string                             `json:"description"`
	DisplayName       string                             `json:"displayName"`
	Offer             *AdvancedCommerceOffer             `json:"offer"`
	Price             int64                              `json:"price"`
	PriceIncreaseInfo *AdvancedCommercePriceIncreaseInfo `json:"priceIncreaseInfo"`
}

type AdvancedCommerceRenewalInfo struct {
	ConsistencyToken   string                        `json:"consistencyToken"`
	Descriptors        *AdvancedCommerceDescriptors  `json:"descriptors"`
	Items              []AdvancedCommerceRenewalItem `json:"items"`
	Period             string                        `json:"period"`
	RequestReferenceId string                        `json:"requestReferenceId"`
	TaxCode            string                        `json:"taxCode"`
}
