package hafenhause

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func (h hafenhausedb) create(newEntries map[docpath]interface{}) error {
	ctx := context.Background()

	batch := h.client.Batch()

	for path, fields := range newEntries {
		docRef := h.client.Doc(string(path))
		batch.Create(docRef, fields)
	}

	_, err := batch.Commit(ctx)

	if grpc.Code(err) == codes.AlreadyExists {
		return errAlreadyExists
	}

	return err
}

func (h hafenhausedb) read(paths ...docpath) ([]*firestore.DocumentSnapshot, error) {
	ctx := context.Background()

	var docs []*firestore.DocumentRef

	for _, path := range paths {
		docs = append(docs, h.client.Doc(string(path)))
	}

	return h.client.GetAll(ctx, docs)
}

func (h hafenhausedb) update(updateEntries map[docpath]interface{}) error {
	ctx := context.Background()

	batch := h.client.Batch()
	for path, fields := range updateEntries {
		// TODO: this creates
		batch.Set(h.client.Doc(string(path)), fields)
	}

	_, err := batch.Commit(ctx)

	if grpc.Code(err) == codes.NotFound {
		return errNotFound
	}

	return err
}

func (h hafenhausedb) delete(paths ...docpath) error {
	ctx := context.Background()

	batch := h.client.Batch()

	for _, path := range paths {
		batch.Delete(h.client.Doc(string(path)))
	}

	_, err := batch.Commit(ctx)

	if grpc.Code(err) == codes.NotFound {
		return errNotFound
	}

	return err
}
