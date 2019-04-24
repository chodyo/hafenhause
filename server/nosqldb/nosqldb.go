package nosqldb

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type nosqldb struct {
	projectID string
	client    *firestore.Client
}

// NosqldbContract is the public interface for the database functions
type NosqldbContract interface {
	Create(docName string, docContents interface{}) (err error)
	Read(docName string) (docContents map[string]interface{}, err error)
	Query(collectionName string, filters ...Filter) (docContents map[string]interface{}, err error)
	Update(docName string, docContents interface{}) (err error)
	Delete(docName string, fields ...string) (err error)
}

var (
	ErrAlreadyExists = errors.New("Resource already exists")
	ErrNotFound      = errors.New("Resource not found")
)

func NewNosqldb() nosqldb {
	ctx := context.Background()

	projectID := os.Getenv("GCLOUD_PROJECT")
	if projectID == "" {
		log.Fatalf("Set Firebase project ID via GCLOUD_PROJECT env variable")
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create client: %v", err)
	}

	return nosqldb{
		projectID: projectID,
		client:    client,
	}
}

func (n nosqldb) Create(docName string, docContents interface{}) (err error) {
	ctx := context.Background()

	doc := n.client.Doc(docName)

	// TODO: check if field already exists
	// Overwrite only the fields in the map; preserve all others
	_, err = doc.Set(ctx, docContents, firestore.MergeAll)

	return
}

func (n nosqldb) Read(docName string) (docContents map[string]interface{}, err error) {
	ctx := context.Background()

	var docSnapshot *firestore.DocumentSnapshot
	if docSnapshot, err = n.client.Doc(docName).Get(ctx); grpc.Code(err) == codes.NotFound {
		err = ErrNotFound
		return
	} else if err != nil {
		return
	}

	docContents = docSnapshot.Data()

	return
}

func (n nosqldb) Query(collectionName string, filters ...Filter) (docContents map[string]interface{}, err error) {
	ctx := context.Background()

	docContents = map[string]interface{}{}

	col := n.client.Collection(collectionName)

	var query firestore.Query
	for i, filter := range filters {
		key, op, value := filter.Tuple()
		if i == 0 {
			query = col.Where(key, op, value)
		} else {
			query = query.Where(key, op, value)
		}
	}

	iter := query.Documents(ctx)

	for {
		var doc *firestore.DocumentSnapshot
		doc, err = iter.Next()
		if err == iterator.Done {
			err = nil
			break
		}
		if err != nil {
			return
		}
		docContents[doc.Ref.ID] = doc.Data()
	}

	return
}

func (n nosqldb) Update(docName string, docContents interface{}) (err error) {
	ctx := context.Background()

	doc := n.client.Doc(docName)

	// TODO: check if field is missing
	// Overwrite only the fields in the map; preserve all others
	_, err = doc.Set(ctx, docContents, firestore.MergeAll)

	return
}

// Delete will delete a document if no fields are given, or it will delete all
// given fields in a document.
func (n nosqldb) Delete(docName string, fields ...string) (err error) {
	ctx := context.Background()

	doc := n.client.Doc(docName)

	if fields != nil {
		var updates []firestore.Update

		for _, field := range fields {
			update := firestore.Update{Path: field, Value: firestore.Delete}

			updates = append(updates, update)
		}

		_, err = doc.Update(ctx, updates)

		return
	}

	_, err = doc.Delete(ctx)

	return
}
