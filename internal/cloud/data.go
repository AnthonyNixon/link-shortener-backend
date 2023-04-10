package data

import (
	datastore "cloud.google.com/go/datastore"
	"context"
	"errors"
	"fmt"
	"github.com/anthonynixon/link-shortener-backend/internal/types"
	"log"
	"os"
	"strings"
)

var datastoreClient *datastore.Client
var ctx context.Context
var namespace string

var AlreadyExistsErr error

func Initialize() {
	projID := os.Getenv("DATASTORE_PROJECT_ID")
	if projID == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_PROJECT_ID"`)
	}

	namespace = os.Getenv("DATASTORE_NAMESPACE")
	if namespace == "" {
		log.Fatal(`You need to set the environment variable "DATASTORE_NAMESPACE"`)
	}

	ctx = context.Background()
	client, err := datastore.NewClient(ctx, projID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	datastoreClient = client

	AlreadyExistsErr = errors.New("entity already exists")
}

func GetLink(short string) (link types.Link, err error) {
	query := datastore.NewQuery("link").
		FilterField("Short", "=", strings.ToUpper(short)).
		Limit(1).
		Namespace(namespace)
	var links []types.Link
	_, err = datastoreClient.GetAll(ctx, query, &links)
	if err != nil {
		return link, err
	}

	fmt.Printf("%v\n", links)
	if len(links) > 0 {
		link = links[0]
	}

	return
}

func NewLink(newLink types.Link) (err error) {
	if !strings.HasPrefix(newLink.Long, "https://") && !strings.HasPrefix(newLink.Long, "http://") {
		newLink.Long = fmt.Sprintf("https://%s", newLink.Long)
	}
	linkKey := datastore.NameKey("link", strings.ToUpper(newLink.Short), nil)
	linkKey.Namespace = "links.ajn.me"
	_, err = datastoreClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		// We first check that there is no entity stored with the given key.
		var empty types.Link
		if err = tx.Get(linkKey, &empty); err != datastore.ErrNoSuchEntity {
			fmt.Printf("empty?: %v\n", empty)
			return AlreadyExistsErr
		}

		// If there was no matching entity, store it now.
		newLink.Short = strings.ToUpper(newLink.Short)
		_, err = tx.Put(linkKey, &newLink)
		return err
	})

	return
}

func IncrementCountInDatastore(link types.Link) {
	linkKey := datastore.NameKey("link", strings.ToUpper(link.Short), nil)
	linkKey.Namespace = "links.ajn.me"

	link.Clicks++

	_, err := datastoreClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		_, err := tx.Put(linkKey, &link)
		return err
	})

	if err != nil {
		fmt.Printf("Error updating click count: %s\n", err.Error())
	}
}
