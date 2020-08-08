package gcp

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/fountain/pkg/fountain"
	"google.golang.org/api/iterator"
)

// DatastoreClient implements Fountain data interfaces for GCP datastore
type DatastoreClient struct{}

// WriteWl writes a new user waterlog to datastore
func (DatastoreClient) WriteWl(userID string, wl *fountain.WaterLog) error {
	ctx := context.Background()

	cl, err := datastore.NewClient(ctx, "water-you-logging-for")
	if err != nil {
		return err
	}

	// insert new item
	userKey := datastore.NameKey("Users", userID, nil)
	wlKey := datastore.IDKey("WaterLogs", wl.Unix(), userKey)
	if _, err := cl.Put(ctx, wlKey, wl); err != nil {
		return err
	}

	return nil
}

// UserWls retrieves waterlogs for a user from GCP datastore
func (DatastoreClient) UserWls(userID string) ([]*fountain.WaterLog, error) {
	var wls []*fountain.WaterLog
	ctx := context.Background()

	client, err := datastore.NewClient(ctx, "water-you-logging-for")
	if err != nil {
		return wls, err
	}

	userKey := datastore.NameKey("Users", userID, nil)
	q := datastore.NewQuery("WaterLogs").Ancestor(userKey)
	qResult := client.Run(ctx, q)

	for {
		wl := new(fountain.WaterLog)

		// retrieve next log
		if _, err = qResult.Next(wl); err != nil {
			if err == iterator.Done {
				break
			} else {
				return wls, err
			}
		}

		wls = append(wls, wl)
	}

	return wls, nil
}
