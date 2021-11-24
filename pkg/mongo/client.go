package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var (
	clients          = make(map[string]*mongo.Client)
	clientsConnected = make(map[string]*mongo.Client)
	clientsMutex     = new(sync.Mutex)
)

func Init(conf map[string]*Config) error {
	Conf = conf
	for conn, conf := range Conf {
		c, err := mongo.NewClient(options.Client().ApplyURI(conf.URI))
		if err != nil {
			return err
		}
		clients[conn] = c
	}
	return nil
}

func Disconnect() error {
	disconnect := func(client *mongo.Client) error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			return err
		}
		return nil
	}

	if len(clientsConnected) > 0 {
		for _, c := range clientsConnected {
			err := disconnect(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetDefaultDB() (*mongo.Database, error) {
	return GetDB("default")
}

func GetClient(connName string) (*mongo.Client, error) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if client, ok := clientsConnected[connName]; ok {
		return client, nil
	} else if client, ok = clients[connName]; ok {
		err := (func(client *mongo.Client) error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := client.Connect(ctx)
			if err != nil {
				return err
			}
			return nil
		})(client)

		if err != nil {
			return nil, err
		}

		clientsConnected[connName] = client

		return clientsConnected[connName], nil
	}
	return nil, fmt.Errorf("conn named [%s] not configed", connName)
}

func GetDB(connName string) (*mongo.Database, error) {
	client, err := GetClient(connName)
	if err != nil {
		return nil, err
	}

	db := client.Database(Conf[connName].DBName)
	if db == nil {
		return nil, fmt.Errorf("database named [%s] not found", Conf[connName].DBName)
	}
	return db, nil
}

func TimeoutCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
