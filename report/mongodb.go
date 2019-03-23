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
    	"log"

    	"github.com/mongodb/mongo-go-driver/bson"
    	"github.com/mongodb/mongo-go-driver/mongo"
    	"github.com/mongodb/mongo-go-driver/mongo/options"

	c "github.com/future-architect/vuls/config"
	"github.com/future-architect/vuls/models"
	"github.com/future-architect/vuls/util"
)

// MongodbWriter writes results to Mongodb
type MongodbWriter struct{}

func getConn() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), c.Conf.Mongodb.URI)

	if err != nil {
            log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
	     log.Fatal(err)
	}

	util.Log.Debugf("Connected to MongoDB!")

	return client
}

// Write results to Mongodb
func (w MongodbWriter) Write(rs ...models.ScanResult) (err error) {
	if len(rs) == 0 {
		return nil
	}

	client := getConn()
	defer client.Disconnect(context.TODO())
	defer util.Log.Debugf("Close Connection.")

	for _, r := range rs {
		key := r.ReportFileName()

		util.Log.Debugf("Insert %s Server Scan Result into MongoDB.",key)

		if err := putData(client, r , key); err != nil {
			util.Log.Errorf("Insert %s Server into MongoDB fail : %s ",key,err)
			return err
		}

	}
	return nil
}

func putData(client *mongo.Client, data models.ScanResult, key string) error {

	collection := client.Database(c.Conf.Mongodb.DB).Collection(c.Conf.Mongodb.Collection)

	filter := bson.M{ "serverName" : key }
	util.Log.Debugf("Update Key: %s ", filter)
	
	var upsert = true 

	rep, err := collection.ReplaceOne( context.TODO(), filter , &data , &options.ReplaceOptions{ Upsert: &upsert })

	if err != nil {
	    return err
	}

	util.Log.Infof("Insert Result to Mongodb finished")

	util.Log.Debugf("DB Insert Result : %+v ", rep)

	return nil
}

