package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
)

func New(payload string, appleRootCert string) *AppStoreServerNotification {
	asn := &AppStoreServerNotification{}
	asn.isValid = false
	asn.appleRootCert = appleRootCert
	asn.parseJwtSignedPayload(payload)
	return asn
}

func (asn *AppStoreServerNotification) extractHeaderByIndex(payload string, index int) ([]byte, error) {
	// get header from token
	payloadArr := strings.Split(payload, ".")

	//convert header to byte
	headerByte, err := base64.RawStdEncoding.DecodeString(payloadArr[0])
	if err != nil {
		return nil, err
	}

	// bind byte to header structure
	var header NotificationHeader
	err = json.Unmarshal(headerByte, &header)
	if err != nil {
		return nil, err
	}

	// decode x.509 certificate headers to byte
	certByte, err := base64.StdEncoding.DecodeString(header.X5c[index])
	if err != nil {
		return nil, err
	}

	return certByte, nil
}

func (asn *AppStoreServerNotification) verifyCertificate(certByte []byte, intermediateCert []byte) error {
	// new empty set of certificate pool
	roots := x509.NewCertPool()

	// parse and append app store certificate to certPool
	ok := roots.AppendCertsFromPEM([]byte(asn.appleRootCert))
	if !ok {
		return errors.New("failed to parse root certificate")
	}

	// parse and append intermediate X5c certificate
	interCert, err := x509.ParseCertificate(intermediateCert)
	if err != nil {
		return errors.New("failed to parse intermediate certificate")
	}
	intermediate := x509.NewCertPool()
	intermediate.AddCert(interCert)

	// parse X5c certificate
	cert, err := x509.ParseCertificate(certByte)
	if err != nil {
		return err
	}

	// append certificate pool to verify options of x509
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediate,
	}

	if _, err := cert.Verify(opts); err != nil { // verify X5c certificate using app store certificate resides in opts
		return err
	}

	return nil
}

func (asn *AppStoreServerNotification) extractPublicKeyFromPayload(payload string) (*ecdsa.PublicKey, error) {
	// get certificate from X5c[0] header
	certStr, err := asn.extractHeaderByIndex(payload, 0)
	if err != nil {
		return nil, err
	}

	// parse certificate
	cert, err := x509.ParseCertificate(certStr)
	if err != nil {
		return nil, err
	}

	// get public key
	switch pk := cert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("appstore public key must be of type ecdsa.PublicKey")
	}
}

func (asn *AppStoreServerNotification) parseJwtSignedPayload(payload string) {
	// get root cert from X5c header
	rootCertStr, err := asn.extractHeaderByIndex(payload, 2)
	if err != nil {
		panic(err)
	}

	// get intermediate cert from X5c header
	intermediateCertStr, err := asn.extractHeaderByIndex(payload, 1)
	if err != nil {
		panic(err)
	}

	// verify certs. if not verified, throw panic
	if err = asn.verifyCertificate(rootCertStr, intermediateCertStr); err != nil {
		panic(err)
	}

	// payload data
	notificationPayload := &NotificationPayload{}
	_, err = jwt.ParseWithClaims(payload, notificationPayload, func(token *jwt.Token) (interface{}, error) {
		return asn.extractPublicKeyFromPayload(payload)
	})
	if err != nil {
		panic(err)
	}
	asn.Payload = notificationPayload

	// transaction info
	transactionInfo := &TransactionInfo{}
	payload = asn.Payload.Data.SignedTransactionInfo
	_, err = jwt.ParseWithClaims(payload, transactionInfo, func(token *jwt.Token) (interface{}, error) {
		return asn.extractPublicKeyFromPayload(payload)
	})
	if err != nil {
		panic(err)
	}
	asn.TransactionInfo = transactionInfo

	// renewal info
	renewalInfo := &RenewalInfo{}
	payload = asn.Payload.Data.SignedRenewalInfo
	_, err = jwt.ParseWithClaims(payload, renewalInfo, func(token *jwt.Token) (interface{}, error) {
		return asn.extractPublicKeyFromPayload(payload)
	})
	if err != nil {
		panic(err)
	}
	asn.RenewalInfo = renewalInfo

	// valid request
	asn.isValid = true
}
