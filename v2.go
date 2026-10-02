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

func (asn *AppStoreServerNotification) extractCertificateChain(payload string) ([]*x509.Certificate, error) {
	segments := strings.Split(payload, ".")
	if len(segments) != 3 {
		return nil, errors.New("payload must be a valid JWS token with 3 segments")
	}

	headerByte, err := base64.RawURLEncoding.DecodeString(segments[0])
	if err != nil {
		return nil, err
	}

	var header NotificationHeader
	if err = json.Unmarshal(headerByte, &header); err != nil {
		return nil, err
	}

	if len(header.X5c) < 3 {
		return nil, fmt.Errorf("x5c header has %d entries, need at least 3", len(header.X5c))
	}

	certificates := make([]*x509.Certificate, 0, len(header.X5c))
	for i, encodedCert := range header.X5c {
		certByte, err := base64.StdEncoding.DecodeString(encodedCert)
		if err != nil {
			return nil, fmt.Errorf("decode x5c certificate %d: %w", i, err)
		}
		cert, err := x509.ParseCertificate(certByte)
		if err != nil {
			return nil, fmt.Errorf("parse x5c certificate %d: %w", i, err)
		}
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

func (asn *AppStoreServerNotification) verifyCertificate(cert *x509.Certificate, chain []*x509.Certificate) error {
	roots := x509.NewCertPool()

	ok := roots.AppendCertsFromPEM([]byte(asn.appleRootCert))
	if !ok {
		return errors.New("root certificate couldn't be parsed")
	}

	intermediate := x509.NewCertPool()
	for _, cert := range chain {
		intermediate.AddCert(cert)
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediate,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("verify certificate chain: %w", err)
	}

	return nil
}

func (asn *AppStoreServerNotification) parseSignedJWT(payload string, claims jwt.Claims) error {
	certificates, err := asn.extractCertificateChain(payload)
	if err != nil {
		return err
	}

	if err = asn.verifyCertificate(certificates[0], certificates[1:]); err != nil {
		return err
	}

	publicKey, ok := certificates[0].PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("appstore public key must be of type ecdsa.PublicKey")
	}

	token, err := jwt.ParseWithClaims(payload, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}))
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid signed payload")
	}
	return nil
}

func (asn *AppStoreServerNotification) parseJwtSignedPayload(payload string) error {
	notificationPayload := &NotificationPayload{}
	if err := asn.parseSignedJWT(payload, notificationPayload); err != nil {
		return err
	}
	asn.Payload = notificationPayload
	asn.IsTest = asn.Payload.NotificationType == "TEST"

	if sti := asn.Payload.Data.SignedTransactionInfo; sti != "" {
		transactionInfo := &TransactionInfo{}
		if err := asn.parseSignedJWT(sti, transactionInfo); err != nil {
			return fmt.Errorf("parse signedTransactionInfo: %w", err)
		}
		asn.TransactionInfo = transactionInfo
	}

	if sri := asn.Payload.Data.SignedRenewalInfo; sri != "" {
		renewalInfo := &RenewalInfo{}
		if err := asn.parseSignedJWT(sri, renewalInfo); err != nil {
			return fmt.Errorf("parse signedRenewalInfo: %w", err)
		}
		asn.RenewalInfo = renewalInfo
	}

	if appTransaction := asn.Payload.AppData; appTransaction != nil && appTransaction.SignedAppTransactionInfo != "" {
		appTransactionInfo := &AppTransactionInfo{}
		if err := asn.parseSignedJWT(appTransaction.SignedAppTransactionInfo, appTransactionInfo); err != nil {
			return fmt.Errorf("parse signedAppTransactionInfo: %w", err)
		}
		asn.AppTransactionInfo = appTransactionInfo
	}

	asn.IsValid = true
	return nil
}
