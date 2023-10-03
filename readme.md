# App Store Server Notifications in Golang [![](https://github.com/izniburak/appstore-notifications-go/actions/workflows/go.yml/badge.svg)](https://github.com/izniburak/appstore-notifications-go/actions) [![PkgGoDev](https://pkg.go.dev/badge/github.com/izniburak/appstore-notifications-go)](https://pkg.go.dev/github.com/izniburak/appstore-notifications-go)

***appstore-notifications-go*** package helps you to handle and process the Apple's [App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications) in Golang.

## Install

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
	fmt.Printf("App Store Server Notification is valid?: %t\n", appStoreServerNotification.isValid)
}
```

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

