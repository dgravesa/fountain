package gcp

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/dgravesa/fountain/pkg/fountain"
	"google.golang.org/api/iterator"
)

// DatastoreClient implements fountain data interfaces for GCP datastore
type DatastoreClient struct {
	projectName string
}

// NewDatastoreClient creates a new GCP datastore client instance
func NewDatastoreClient(projectName string) *DatastoreClient {
	client := new(DatastoreClient)
	client.projectName = projectName
	return client
}

func (c DatastoreClient) get(ctx context.Context) (*datastore.Client, error) {
	return datastore.NewClient(ctx, c.projectName)
}

func userKey(userID string) *datastore.Key {
	return datastore.NameKey("Users", userID, nil)
}

// User retrieves a user from datastore
func (c DatastoreClient) User(userID string) (*fountain.User, error) {
	ctx := context.Background()

	cl, err := c.get(ctx)
	if err != nil {
		return nil, err
	}

	// retrieve user info from datastore
	user := new(fountain.User)
	k := userKey(userID)
	if err = cl.Get(ctx, k, user); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// PutUser creates a new user or updates an existing user with the ID
func (c DatastoreClient) PutUser(user *fountain.User) error {
	ctx := context.Background()

	cl, err := c.get(ctx)
	if err != nil {
		return err
	}

	// write user info
	k := userKey(user.ID)
	_, err = cl.Put(ctx, k, user)
	return err
}

// WriteWl writes a new user waterlog to datastore
func (c DatastoreClient) WriteWl(userID string, wl *fountain.WaterLog) error {
	ctx := context.Background()

	cl, err := c.get(ctx)
	if err != nil {
		return err
	}

	// insert new item
	uKey := userKey(userID)
	wlKey := datastore.IDKey("WaterLogs", wl.Unix(), uKey)
	if _, err := cl.Put(ctx, wlKey, wl); err != nil {
		return err
	}

	return nil
}

// UserWls retrieves waterlogs for a user from GCP datastore
func (c DatastoreClient) UserWls(userID string) ([]*fountain.WaterLog, error) {
	var wls []*fountain.WaterLog
	ctx := context.Background()

	client, err := c.get(ctx)
	if err != nil {
		return wls, err
	}

	uKey := userKey(userID)
	q := datastore.NewQuery("WaterLogs").Ancestor(uKey)
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
