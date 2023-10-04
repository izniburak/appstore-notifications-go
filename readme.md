# App Store Server Notifications in Golang [![](https://github.com/izniburak/appstore-notifications-go/actions/workflows/go.yml/badge.svg)](https://github.com/izniburak/appstore-notifications-go/actions) [![PkgGoDev](https://pkg.go.dev/badge/github.com/izniburak/appstore-notifications-go)](https://pkg.go.dev/github.com/izniburak/appstore-notifications-go)

***appstore-notifications-go*** is a Golang package, which helps you to handle, parse and verify the Apple's [App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications).

> App Store Server Notifications is a service provided by Apple for its App Store. It's designed to notify developers about key events and changes related to their app's in-app purchases and subscriptions. By integrating this service into their server-side logic, developers can receive real-time (I think, almost real-time) updates about various events without having to repeatedly poll the Apple servers.

## Install
To install the package, you can use following command on your terminal in your project directory:

```bash
go get github.com/izniburak/appstore-notifications-go
```

## Examples
```go
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	appstore "github.com/izniburak/appstore-notifications-go"
)

func main() {
	// App Store Server Notification Request JSON String
	appStoreServerRequest := "..." // {"signedPayload":"..."}
	var request appstore.AppStoreServerRequest
	err := json.Unmarshal([]byte(appStoreServerRequest), &request) // bind byte to header structure
	if err != nil {
		panic(err)
	}

	// Apple Root CA - G3 Root certificate
	// for details: https://www.apple.com/certificateauthority/
	// you need download it and covert it to a valid pem file in order to verify X5c certificates
	// `openssl x509 -in AppleRootCA-G3.cer -out cert.pem`
	rootCert := "-----BEGIN CERTIFICATE----- ......"
	if rootCert == "" {
		panic("Apple Root Cert not valid")
	}

	appStoreServerNotification := appstore.New(request.SignedPayload, rootCert)
	fmt.Printf("App Store Server Notification is valid?: %t\n", appStoreServerNotification.IsValid)
	fmt.Printf("Product Id: %s\n", appStoreServerNotification.TransactionInfo.ProductId)
}
```
You can access the all data in the payload by using one of the 4 params in instance of the `AppStoreServerNotification`:

- _instance_***.Payload***: Access the [Payload](https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload).
- _instance_***.TransactionInfo***: Access the [Transaction Info](https://developer.apple.com/documentation/appstoreservernotifications/jwstransactiondecodedpayload).
- _instance_***.RenewalInfo***: Access the [Renewal Info](https://developer.apple.com/documentation/appstoreservernotifications/jwsrenewalinfodecodedpayload).
- _instance_***.IsValid***: Check the payload parsed and verified successfully.

## Contributing

1. Fork [this repo](https://github.com/izniburak/appstore-notifications-go/fork)
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create a new Pull Request

## Contributors

- [izniburak](https://github.com/izniburak) İzni Burak Demirtaş - creator, maintainer

## License
The MIT License (MIT) - see [`license.md`](https://github.com//izniburak/appstore-notifications-go/blob/main/license.md) for more details

