[![reverse-ip-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![reverse-ip-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/reverse-ip-go)
[![reverse-ip-go test](https://github.com/whois-api-llc/reverse-ip-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/reverse-ip-go/actions/)

# Overview

The client library for
[Reverse IP/DNS API](https://dns-history.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/reverse-ip-go
```

# Examples

Full API documentation available [here](https://dns-history.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := reverseip.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := reverseip.NewClient(apiKey, reverseip.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Reverse IP/DNS API lets you get a list of all domains associated with an IP address.

```go

// Make request to get a list of all domains by IP address as a model instance.
reverseIPResp, _, err := client.Get(ctx, []byte{8,8,8,8})
if err != nil {
    log.Fatal(err)
}

for _, obj := range reverseIPResp.Result {
    log.Println(obj.Name, obj.FirstSeen, obj.LastVisit)
}

// Make request to get raw data in XML.
resp, err := client.GetRaw(context.Background(), net.IP{1, 1, 1, 1},
    reverseip.OptionOutputFormat("XML"))
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))

```
## Advanced usage

Pagination

```go

// Each response is limited to 300 records.
const limit = 300
from := "1"

for {
    reverseIPResp, _, err := client.Get(context.Background(), []byte{8, 8, 8, 8},
    //  This option results in the next page is retrieved.
    reverseip.OptionFrom(from))
    if err != nil {
        log.Fatal(err)
    }

    for _, obj := range reverseIPResp.Result {
        log.Println(obj.Name)
    }

    // Break the loop when the last page is reached.
    if reverseIPResp.Size < limit {
        break
    }
    from = reverseIPResp.Result[reverseIPResp.Size-1].Name
}

```