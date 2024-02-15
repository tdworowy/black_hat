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

func New(host string) (*MongoMiner, error) {
	m := MongoMiner{Host: host}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MongoMiner) connect() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Host))

	if err != nil {
		panic(err)
	}
	m.Context = ctx
	m.Client = client
	return nil

}

func (m *MongoMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema)

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
			cur, err := m.Client.Database(dbname).Collection(collection).Find(m.Context, bson.D{})
			if err != nil {
				return nil, err
			}

			defer cur.Close(m.Context)
			for cur.Next(m.Context) {
				var f dbminer.Table
				err := cur.Decode(&f)
				if err != nil {
					return nil, err
				}
				table.Columns = append(table.Columns, f.Name)
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
} // TODO somethin is wrong

func main() {

	var addr string
	if len(os.Args) > 1 {
		addr = os.Args[1]
	} else {
		addr = "mongodb://localhost:27017"
	}

	mm, err := New(addr)
	if err != nil {
		panic(err)
	}
	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
