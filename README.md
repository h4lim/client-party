Client Party is simple go http client library

### Client Party support many content type like below

- application/json
- application/x-www-form-urlencoded
- form data
- and more

### Getting the library

With [Go module](https://github.com/golang/go/wiki/Modules) support, simply add the following import

```
import "github.com/h4lim/client-party"
```

to your code, and then `go [build|run|test]` will automatically fetch the necessary dependencies.

Otherwise, run the following Go command to install the `qr` package:

```sh
$ go get -u github.com/h4lim/client-party
```

### JSON TEST

First you need to import client-party , one simplest example likes the follow `example.go`:

```go
package main

import (
	"fmt"
	cp "github.com/h4lim/client-party"
)

func main() {
	method := cp.MethodGet
	url := "http://facebook.com"
	response, err := cp.NewClientParty(method, url).HitClient()
	if err != nil {
		fmt.Println(*err)
		return
	}

	fmt.Println(response)
}
```

And use the Go command to run the demo:

```
# run example.go
$ go run example.go
```

The output will be like below:

```
{"tag_00":{"version":"01","type":"11","tag_52":"5072","tag_53":"360","tag_58":"ID","tag_61":"40271","tag_62":"0703A01","amount":0,"merchant_owner":"PERKAKASKU","merchant_address":"BANDUNG","checksum":"4D4A"},"tag_26":{"qr_owner":"ID.CO.BCA.WWW","merchant_id":"936000140000940453","merchant_acquirer_id":"000885000940453","merchant_scale":"UKE"},"tag_51":{"qris_web":"ID.CO.QRIS.WWW","qris_id":"ID2020034073193","scale":"UKE"}}
```
