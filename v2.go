package v2

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func New(payload string, appleRootCert string) (*AppStoreServerNotification, error) {
	asn := &AppStoreServerNotification{appleRootCert: appleRootCert}
	if err := asn.parseJwtSignedPayload(payload); err != nil {
		return nil, err
	}
	return asn, nil
}

func (asn *AppStoreServerNotification) extractHeaderByIndex(payload string, index int) ([]byte, error) {
	payloadArr := strings.Split(payload, ".")
	if len(payloadArr) < 3 {
		return nil, errors.New("payload must be a valid JWS token with 3 segments")
	}

	headerByte, err := base64.RawStdEncoding.DecodeString(payloadArr[0])
	if err != nil {
		return nil, err
	}

	var header NotificationHeader
	if err = json.Unmarshal(headerByte, &header); err != nil {
		return nil, err
	}

	if len(header.X5c) <= index {
		return nil, fmt.Errorf("x5c header has %d entries, need at least %d", len(header.X5c), index+1)
	}

	certByte, err := base64.StdEncoding.DecodeString(header.X5c[index])
	if err != nil {
		return nil, err
	}

	return certByte, nil
}

func (asn *AppStoreServerNotification) verifyCertificate(certByte []byte, intermediateCert []byte) error {
	roots := x509.NewCertPool()

	ok := roots.AppendCertsFromPEM([]byte(asn.appleRootCert))
	if !ok {
		return errors.New("root certificate couldn't be parsed")
	}

	interCert, err := x509.ParseCertificate(intermediateCert)
	if err != nil {
		return fmt.Errorf("intermediate certificate: %w", err)
	}
	intermediate := x509.NewCertPool()
	intermediate.AddCert(interCert)

	cert, err := x509.ParseCertificate(certByte)
	if err != nil {
		return err
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediate,
	}
	if _, err := cert.Verify(opts); err != nil {
		return err
	}

	return nil
}

func (asn *AppStoreServerNotification) extractPublicKeyFromPayload(payload string) (*ecdsa.PublicKey, error) {
	certStr, err := asn.extractHeaderByIndex(payload, 0)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certStr)
	if err != nil {
		return nil, err
	}

	switch pk := cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("appstore public key must be of type ecdsa.PublicKey")
	}
}

func (asn *AppStoreServerNotification) parseJwtSignedPayload(payload string) error {
	rootCertStr, err := asn.extractHeaderByIndex(payload, 2)
	if err != nil {
		return err
	}

	intermediateCertStr, err := asn.extractHeaderByIndex(payload, 1)
	if err != nil {
		return err
	}

	if err = asn.verifyCertificate(rootCertStr, intermediateCertStr); err != nil {
		return err
	}

	notificationPayload := &NotificationPayload{}
	_, err = jwt.ParseWithClaims(payload, notificationPayload, func(token *jwt.Token) (interface{}, error) {
		return asn.extractPublicKeyFromPayload(payload)
	})
	if err != nil {
		return err
	}
	asn.Payload = notificationPayload
	asn.IsTest = asn.Payload.NotificationType == "TEST"

	if sti := asn.Payload.Data.SignedTransactionInfo; sti != "" {
		transactionInfo := &TransactionInfo{}
		_, err = jwt.ParseWithClaims(sti, transactionInfo, func(token *jwt.Token) (interface{}, error) {
			return asn.extractPublicKeyFromPayload(sti)
		})
		if err != nil {
			return fmt.Errorf("parse signedTransactionInfo: %w", err)
		}
		asn.TransactionInfo = transactionInfo
	}

	if sri := asn.Payload.Data.SignedRenewalInfo; sri != "" {
		renewalInfo := &RenewalInfo{}
		_, err = jwt.ParseWithClaims(sri, renewalInfo, func(token *jwt.Token) (interface{}, error) {
			return asn.extractPublicKeyFromPayload(sri)
		})
		if err != nil {
			return fmt.Errorf("parse signedRenewalInfo: %w", err)
		}
		asn.RenewalInfo = renewalInfo
	}

	asn.IsValid = true
	return nil
}
