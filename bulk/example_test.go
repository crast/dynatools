package bulk_test

import (
	"fmt"
	"log"

	"github.com/crast/dynatools/bulk"
	"gopkg.in/underarmour/dynago.v1"
)

// replace with a real executor here
var executor = &dynago.MockExecutor{}

func Example() {
	client := dynago.NewClient(executor)
	writer := bulk.New(bulk.Config{
		Client: client,
		Table:  "people",
	})

	// Start a goroutine to read results
	go resultReader(writer.Results())

	// Start putting 1000 records of bulk data
	for i := 0; i < 1000; i++ {
		writer.Write(dynago.Document{
			"Id":   i,
			"Name": fmt.Sprintf("User%d", i),
		})
	}

	writer.CloseWait()
}

func resultReader(results <-chan bulk.Result) {
	for result := range results {
		if result.Error != nil {
			log.Printf("Got error writing documents: %v: %v", result.Documents, result.Error)
		} else {
			log.Printf("Successfully wrote %d documents", len(result.Documents))
		}
	}
}
