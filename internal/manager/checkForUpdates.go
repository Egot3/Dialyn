package manager

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	banye "github.com/Egot3/Banye"
	"github.com/Egot3/Zhao/queues"
)

func CheckForUpdates(ctx context.Context, cl *banye.Client, ch chan<- []*queues.QueueStruct) {
	tParsed, _ := time.ParseDuration(os.Getenv("UPDATEINTERVAL"))
	interval := tParsed * time.Second
	var upd [32]uint8
	var runQ []*queues.QueueStruct

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Guess the updateChan is closed")
			return
		case <-ticker.C:
			req, _ := http.NewRequest("GET", "http://localhost:15672/api/queues", nil)
			req.SetBasicAuth(os.Getenv("RABBIT_LOGIN"), os.Getenv("RABBIT_PASSWORD"))

			resp, err := cl.Do(req)
			if err != nil {
				log.Printf("failed to send a request to queues ep: %v", err)
				continue
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("failed to read body: %v", err)
				continue
			}

			updN := sha256.Sum256(bodyBytes)

			if updN != upd {
				upd = updN
				json.Unmarshal(bodyBytes, &runQ)
				//queueNames
				ch <- runQ
			}

			resp.Body.Close()
		}
	}
}
