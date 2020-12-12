package gomongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//Identity type
type Identity struct {
	CertificateID     int32  `bson:"certificate_id,omitempty"`
	CertificateNumber string `bson:"certificate_number,omitempty"`
	Date              string `bson:"date,omitempty"`
	Time              string `bson:"time,omitempty"`
	BandUnit          string `bson:"band_unit,omitempty"`
	Band              int32  `bson:"band,omitempty"`
	FrequencyUnit     string `bson:"frequency_unit,omitempty"`
	Frequency         int32  `bson:"frequency,omitempty"`
	CallSign          string `bson:"call_sign,omitempty"`
	Name              string `bson:"name,omitempty"`
	EventID           int32  `bson:"event_id,omitempty"`
	DateCreated       string `bson:"date_created,omitempty"`
	CreatedBy         string `bson:"created_by,omitempty"`
	DateModified      string `bson:"date_modified,omitempty"`
	ModifiedBy        string `bson:"modified_by,omitempty"`
	DownloadCount     int32  `bson:"download_count,omitempty"`
	CityID            int32  `bson:"city_id,omitempty"`
}

// Adaptor Type
type Adaptor struct {
	Client   mongo.Client
	DBName   string
	CollName string
}

//Connect method
func (adaptor *Adaptor) Connect(ctx context.Context, uri string) {
	Client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	adaptor.Client = *Client

	err = adaptor.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

//Insert method
func (adaptor *Adaptor) Insert(ctx context.Context, identity Identity) {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	_, err := collection.InsertOne(ctx, identity)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(identity.Name, "inserted.")
}

//GetQuery method
func (adaptor *Adaptor) GetQuery(ctx context.Context, queryName bson.M, isStringReturned bool) interface{} {
	var received bson.M

	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)

	if errFinding := collection.FindOne(ctx, queryName).Decode(&received); errFinding != nil {
		fmt.Println(errFinding)
	}

	if isStringReturned {
		// JSON
		jsonBytes, _ := bson.MarshalJSON(&received)

		return string(jsonBytes)

	} else {
		// STRUCT
		bsonBytes, _ := bson.Marshal(&received)
		var result Identity
		bson.Unmarshal(bsonBytes, &result)

		return result
	}
}

//GetQueries method
func (adaptor *Adaptor) GetQueries(ctx context.Context, queryName bson.M, optionsFilter bson.M, isStringReturned bool) interface{} {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)

	findOptions := options.Find()

	adaptor.filterHandler(optionsFilter, findOptions)

	cursor, _ := collection.Find(ctx, queryName, findOptions)
	var result []interface{}
	for cursor.Next(ctx) {
		received := bson.M{}

		if err := cursor.Decode(&received); err != nil {
			fmt.Println(err)
		}

		bsonBytes, _ := bson.Marshal(&received)

		var subIdentity bson.M
		bson.Unmarshal(bsonBytes, &subIdentity)

		result = append(result, subIdentity)
	}

	if isStringReturned {
		// JSON
		jsonBytes, _ := bson.MarshalJSON(&result)

		// return string(jsonBytes)
		return string(jsonBytes)

	} else {
		// STRUCT
		bsonBytes, _ := bson.Marshal(&result)
		var result Identity
		bson.Unmarshal(bsonBytes, &result)

		return result
	}
}

//GetIdentities method
func (adaptor *Adaptor) GetIdentities(ctx context.Context, queryName bson.M, name string) []Identity {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	cursor, err := collection.Find(ctx, queryName)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result []Identity

	for cursor.Next(ctx) {
		received := bson.M{}

		if err := cursor.Decode(&received); err != nil {
			fmt.Println(err)
		}

		bsonBytes, _ := bson.Marshal(&received)
		var subIdentity Identity
		bson.Unmarshal(bsonBytes, &subIdentity)

		result = append(result, subIdentity)
	}
	return result
}

//DeleteIdentity method
func (adaptor *Adaptor) DeleteIdentity(ctx context.Context, name string) {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	delResult, err := collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delResult)
}

//DeleteCollection method
func (adaptor *Adaptor) DeleteCollection(ctx context.Context, name string, deleteCode string) {
	if deleteCode == "AGREE TO DELETE "+adaptor.CollName {
		collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
		delResult, err := collection.DeleteMany(ctx, bson.M{"name": name})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(delResult)
	} else {
		fmt.Println("ACCESS DENIED")
	}
}

///////////// FILTER HANDLDER /////////////
func (adaptor *Adaptor) filterHandler(optionsFilter bson.M, mongoFindOptions *options.FindOptions) {
	SetLimitHandler := func(limitValue int64) {
		mongoFindOptions.SetLimit(limitValue)
	}

	for key, value := range optionsFilter["options"].(map[string]interface{}) {
		// fmt.Printf("%v %T", key, value)
		if key == "limit" {
			SetLimitHandler(int64(value.(float64)))
		}
	}
}
