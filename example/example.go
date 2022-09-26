package example

import (
	"context"
	"errors"
	reverseip "github.com/whois-api-llc/reverse-ip-go"
	"log"
	"net"
	"time"
)

func GetData(apikey string) {
	client := reverseip.NewBasicClient(apikey)

	// Get parsed Reverse IP/DNS API response by IP address as a model instance.
	reverseIPResp, resp, err := client.Get(context.Background(),
		net.IP{8, 8, 8, 8},
		// this option is ignored, as the inner parser works with JSON only.
		reverseip.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *reverseip.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print some values from each returned record.
	for _, obj := range reverseIPResp.Result {
		log.Printf("Name: %s, First: %s, Last: %s\n",
			obj.Name,
			time.Unix(obj.FirstSeen, 0).Format(time.RFC3339),
			time.Unix(obj.LastVisit, 0).Format(time.RFC3339),
		)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetAllData(apikey string) {
	// Each response is limited to 300 records.
	const limit = 300

	from := "1"
	var results []reverseip.Result
	client := reverseip.NewBasicClient(apikey)

	for {
		reverseIPResp, _, err := client.Get(context.Background(), []byte{8, 8, 8, 8},
			//  This option results in the next page is retrieved.
			reverseip.OptionFrom(from))
		if err != nil {
			log.Println(err)
			return
		}

		// Store all returned records in the single slice.
		results = append(results, reverseIPResp.Result...)

		// Break the loop when the last page is reached.
		if reverseIPResp.Size < limit {
			break
		}
		from = reverseIPResp.Result[reverseIPResp.Size-1].Name
	}

	// Then print the count and some values from each record.
	log.Println(len(results))
	for _, obj := range results {
		log.Printf("Name: %s, First: %s, Last: %s\n",
			obj.Name,
			time.Unix(obj.FirstSeen, 0).Format(time.RFC3339),
			time.Unix(obj.LastVisit, 0).Format(time.RFC3339),
		)
	}
}

func GetRawData(apikey string) {
	client := reverseip.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.GetRaw(context.Background(),
		net.IP{1, 1, 1, 1},
		reverseip.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
