package bedtimes

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

const collection string = "hafenhaus"
const doc string = "bedtimes"

var projectID string
var client *firestore.Client

func init() {
	ctx := context.Background()

	projectID = os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("Set Firebase project ID via GCLOUD_PROJECT env variable")
	}

	var err error
	client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}
}
