package hafenhause

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

const (
	// cody      string = "Cody"
	// julia     string = "Julia"
	brannigan string = "Brannigan"
	malcolm   string = "Malcolm"
)

type hafenhausedb struct {
	db         *firestore.CollectionRef
	collection string

	projectID string
	client    *firestore.Client
}

func newHafenhausedb() hafenhausedb {
	ctx := context.Background()

	projectID := os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("Set Firebase project ID via GCLOUD_PROJECT env variable")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	collection := "hafenhause"

	collectionRef := client.Collection(collection)

	return hafenhausedb{
		db:         collectionRef,
		collection: collection,
		projectID:  projectID,
		client:     client,
	}
}

type defaults struct {
	Hour   int `firestore:"hour"`
	Minute int `firestore:"minute"`
}

func (h hafenhausedb) createBedtimes(names []string) error {
	ctx := context.Background()

	var defaults defaults
	defaultsSnap, err := h.client.Doc(h.collection + "/bedtimes/static/defaults").Get(ctx)
	if err != nil {
		return err
	}
	if err := defaultsSnap.DataTo(&defaults); err != nil {
		return err
	}

	bedtimes := []bedtime{}
	for _, name := range names {
		bedtimes = append(bedtimes, bedtime{
			Name: name, Hour: defaults.Hour, Minute: defaults.Minute,
		})
	}
	return h.updateBedtimes(bedtimes)
}

func (h hafenhausedb) readBedtimes(names []string) ([]bedtime, error) {
	ctx := context.Background()

	var docs []*firestore.DocumentRef
	for _, name := range names {
		docs = append(docs, h.client.Doc(h.collection+"/bedtimes/people/"+name))
	}

	docsnaps, err := h.client.GetAll(ctx, docs)
	if err != nil {
		return nil, err
	}

	var bedtimes []bedtime

	for _, ds := range docsnaps {
		var bedtime bedtime

		if err := ds.DataTo(&bedtime); err != nil {
			return nil, err
		}

		bedtime.Name = ds.Ref.ID

		bedtimes = append(bedtimes, bedtime)
	}

	return bedtimes, nil
}

func (h hafenhausedb) updateBedtimes(bedtimes []bedtime) error {
	ctx := context.Background()

	now := time.Now()

	batch := h.client.Batch()

	for _, b := range bedtimes {
		personRef := h.client.Doc(h.collection + "/bedtimes/people/" + b.Name)
		batch.Set(personRef, map[string]interface{}{
			"hour":    b.Hour,
			"minute":  b.Minute,
			"updated": now,
		})
	}

	_, err := batch.Commit(ctx)

	return err
}

func (h hafenhausedb) deleteBedtimes(recordNames []string) error {
	return fmt.Errorf("Not yet implemented")
}
