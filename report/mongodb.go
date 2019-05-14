/* Vuls - Vulnerability Scanner
Copyright (C) 2016  Future Corporation , Japan.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package report

import (

    	"context"
	"time"

        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
        "go.mongodb.org/mongo-driver/bson"

	c "github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/models"
	"github.com/future-architect/vuls/util"
)

// MongodbWriter writes results to Mongodb
type MongodbWriter struct{}

func NewMongoClient() (DBClient *mongo.Client,err error) {

        client, err := mongo.NewClient(options.Client().ApplyURI( c.Conf.Mongodb.URI ))
        if err != nil { return nil,err }

        ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
            util.Log.Fatalf("Could not connect to MongoDB : %+v ", err)
	    return nil, err
	}

	util.Log.Debugf("Connected to MongoDB!")

	return client, nil
}

// Write results to Mongodb
func (w MongodbWriter) Write(rs ...models.ScanResult) (err error) {
	if len(rs) == 0 {
		return nil
	}

        client, err := mongo.NewClient(options.Client().ApplyURI( c.Conf.Mongodb.URI ))
        if err != nil { return err }

        ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
        err = client.Connect(ctx)
        defer cancel()
        if err != nil { return err }

        defer client.Disconnect(context.Background())


	for _, r := range rs {

	  key := r.ReportFileName()

	  filter := bson.M{ "serverName" : key }

	  util.Log.Debugf("Insert %s Server Scan Result into MongoDB.",key)

	  opt := options.Replace().SetUpsert(true)

	  if _, err := client.Database( c.Conf.Mongodb.DB ).Collection( c.Conf.Mongodb.Collection ).ReplaceOne( context.Background(), filter , r , opt ); err!= nil {
	    util.Log.Errorf("Insert %s Server into MongoDB fail : %v ",key, err)
	    return err
	  }

	}
	return nil
}
/*
func putData(client *mongo.Client, data models.ScanResult, key string) error {

	collection := client.Database(c.Conf.Mongodb.DB).Collection(c.Conf.Mongodb.Collection)

	filter := bson.M{ "serverName" : key }
	util.Log.Debugf("Update Key: %s ", filter)
	
	var upsert = true 

	rep, err := collection.ReplaceOne( context.Background(), filter , &data , &options.ReplaceOptions{ Upsert: &upsert })

	if err != nil {
	    return err
	}

	util.Log.Infof("Insert Result to Mongodb finished")

	return nil
}
*/
