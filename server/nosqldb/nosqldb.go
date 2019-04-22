package nosqldb

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

type nosqldb struct {
	collection string
	projectID  string
	client     *firestore.Client
}

// NosqldbContract is the public interface for the database functions
type NosqldbContract interface {
	Create(docName, key string, value interface{}) error
	Read(docName string, keys []string) ([]interface{}, error)
	Update(docName, key string, value interface{}) error
	Delete(docName, key string) error
}

var (
	ErrAlreadyExists = errors.New("Resource already exists")
	ErrNotFound      = errors.New("Resource not found")
)

func NewNosqldb(collectionName string) nosqldb {
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
		collection: collectionName,
		projectID:  projectID,
		client:     client,
	}
}

func (n nosqldb) Create(docName, key string, value interface{}) (err error) {
	ctx := context.Background()

	doc := n.client.Collection(n.collection).Doc(docName)

	var docSnapshot *firestore.DocumentSnapshot
	if docSnapshot, err = doc.Get(ctx); err != nil {
		return err
	}

	if _, err = docSnapshot.DataAt(key); err == nil {
		return ErrAlreadyExists
	}

	_, err = doc.Update(ctx, []firestore.Update{{
		Path:  key,
		Value: value,
	}})

	return
}

func (n nosqldb) Read(docName string, keys []string) (objects []interface{}, err error) {
	ctx := context.Background()

	var statesSnapshot *firestore.DocumentSnapshot
	if statesSnapshot, err = n.client.Collection(n.collection).Doc(docName).Get(ctx); err != nil {
		return
	}

	for _, key := range keys {
		var object interface{}
		if object, err = statesSnapshot.DataAt(key); err != nil {
			return nil, ErrNotFound
		}

		objects = append(objects, object)
	}

	return
}

func (n nosqldb) Update(docName, key string, value interface{}) (err error) {
	ctx := context.Background()

	doc := n.client.Collection(n.collection).Doc(docName)

	var docSnapshot *firestore.DocumentSnapshot
	if docSnapshot, err = doc.Get(ctx); err != nil {
		return err
	}

	if _, err = docSnapshot.DataAt(key); err != nil {
		return ErrNotFound
	}

	_, err = doc.Update(ctx, []firestore.Update{{
		Path:  key,
		Value: value,
	}})

	return err
}

func (n nosqldb) Delete(docName, key string) error {
	var empty struct{}
	return n.Update(docName, key, empty)
}
