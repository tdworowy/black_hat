package main

import (
	"black_hat_go/databases/dbminer"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMiner struct {
	Host    string
	Client  *mongo.Client
	Context context.Context
}

func new(host string) (*MongoMiner, error) {
	m := MongoMiner{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MongoMiner) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Host))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	m.Context = ctx
	m.Client = client
	return nil

}

func (m *MongoMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema) // how to fix it ?

	dbnames, err := m.Client.ListDatabaseNames(m.Context, bson.D{})
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbnames {
		db := dbminer.Database{Name: dbname, Tables: []dbminer.Table{}}
		collections, err := m.Client.Database(dbname).ListCollectionNames(m.Context, bson.D{})
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			table := dbminer.Table{Name: collection, Columns: []string{}}

			var docRaw bson.Raw
			_, err := m.Client.Database(dbname).Collection(collection).Find(m.Context, bson.D{})
			if err != nil {
				return nil, err
			}

			var doc bson.Raw
			if err := docRaw.Unmarshal(&doc); err != nil { // how to fix it ?
				if err != nil {
					return nil, err
				}
			}

			for _, f := range doc {
				table.Columns = append(table.Columns, f.Name) // how to fix it ?
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
}

func main() {

	mm, err := new(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
