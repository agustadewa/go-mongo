package gomongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Payload type
type Payload struct {
	Kind    string      `bson:"kind"`
	Values  interface{} `bson:"values"`
	Options FindOptions `bson:"options",omitempy`
}

//FindOptions type
type FindOptions struct {
	Limit int64 `bson:"limit"`
}

//FullIdentity type
type FullIdentity struct {
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

//Identity type
type Identity struct {
	Name     string `bson:"name"`
	Band     int32  `bson:"band"`
	CallSign string `bson:"call_sign"`
}

// EventCallSign type
type EventCallSign struct {
	Description         string `bson:"description"`
	Date                int64  `bson:"date"`
	CertificateFormat   string `bson:"format_certificate"`
	CertificateTemplate string `bson:"certificate_template"`
	Name                string `bson:"name"`
	IsActive            bool   `bson:"is_active"`
	IsHidden            bool   `bson:"is_hidden"`
}

// Adaptor Type
type Adaptor struct {
	Client   mongo.Client
	DBName   string
	CollName string
	Temp     interface{}
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

//QueryFind query find to mongodb
func (adaptor *Adaptor) QueryFind(ctx context.Context, query bson.M) ([]byte, error) {
	var received bson.M

	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	errFinding := collection.FindOne(ctx, query).Decode(&received)
	jsonBytes, _ := bson.MarshalJSON(&received)

	return jsonBytes, errFinding
}

//QueryInsert Query Insert to mongodb
func (adaptor *Adaptor) QueryInsert(ctx context.Context, byteQuery []byte) (interface{}, error) {
	var insertResult interface{}
	var errorInserting error

	var query bson.M
	fmt.Println(string(byteQuery))
	bson.UnmarshalJSON(byteQuery, &query)

	collection := adaptor.Client.Database(adaptor.DBName).Collection("event")
	insertResult, errorInserting = collection.InsertOne(ctx, query)

	return insertResult, errorInserting
}

// //Insert method
// func (adaptor *Adaptor) Insert(ctx context.Context, identity FullIdentity) {
// 	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 	_, err := collection.InsertOne(ctx, identity)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Println(identity.Name, "inserted.")
// }

// //GetQueries method
// func (adaptor *Adaptor) GetQueries(ctx context.Context, queryName bson.M, optionsFilter bson.M, isStringReturned bool) interface{} {
// 	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
//
// 	findOptions := options.Find()
//
// 	adaptor.filterHandler(optionsFilter, findOptions)
//
// 	cursor, _ := collection.Find(ctx, queryName, findOptions)
// 	var result []bson.M
// 	for cursor.Next(ctx) {
// 		received := bson.M{}
//
// 		if err := cursor.Decode(&received); err != nil {
// 			fmt.Println(err)
// 		}
// 		result = append(result, received)
// 	}
//
// 	if isStringReturned {
// 		// JSON
// 		// jsonBytes, _ := bson.MarshalJSON(&result)
//
// 		// return string(jsonBytes)
//
// 		var testvar interface{}
// 		Parser{}.Parse(&result, &testvar, true)
//
// 		return testvar.(string)
//
// 	} else {
// 		// STRUCT
// 		bsonBytes, _ := bson.Marshal(&result)
// 		var result FullIdentity
// 		bson.Unmarshal(bsonBytes, &result)
//
// 		return result
// 	}
// }

// //GetIdentities method
// func (adaptor *Adaptor) GetIdentities(ctx context.Context, queryName bson.M, name string) []FullIdentity {
// 	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 	cursor, err := collection.Find(ctx, queryName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(ctx)
//
// 	var result []FullIdentity
// 	for cursor.Next(ctx) {
// 		received := bson.M{}
//
// 		if err := cursor.Decode(&received); err != nil {
// 			fmt.Println(err)
// 		}
//
// 		bsonBytes, _ := bson.Marshal(&received)
// 		var subIdentity FullIdentity
// 		bson.Unmarshal(bsonBytes, &subIdentity)
//
// 		result = append(result, subIdentity)
// 	}
// 	return result
// }

// //DeleteIdentity method
// func (adaptor *Adaptor) DeleteIdentity(ctx context.Context, name string) {
// 	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 	delResult, err := collection.DeleteOne(ctx, bson.M{"name": name})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Println(delResult)
// }

// //DeleteCollection method
// func (adaptor *Adaptor) DeleteCollection(ctx context.Context, name string, deleteCode string) {
// 	if deleteCode == "AGREE TO DELETE "+adaptor.CollName {
// 		collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
// 		delResult, err := collection.DeleteMany(ctx, bson.M{"name": name})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		fmt.Println(delResult)
// 	} else {
// 		fmt.Println("ACCESS DENIED")
// 	}
// }

///////////// FILTER HANDLDER /////////////
// func (adaptor *Adaptor) filterHandler(optionsFilter bson.M, mongoFindOptions *options.FindOptions) {
// 	SetLimitHandler := func(limitValue int64) {
// 		mongoFindOptions.SetLimit(limitValue)
// 	}
//
// 	for key, value := range optionsFilter["options"].(map[string]interface{}) {
// 		// fmt.Printf("%v %T", key, value)
// 		if key == "limit" {
// 			SetLimitHandler(int64(value.(float64)))
// 		}
// 	}
// }

///////////// PAYLOAD FILTER /////////////

// ParsePayload method
func (adaptor *Adaptor) ParsePayload(jsonByte []byte, out interface{}) {
	if isErr := bson.UnmarshalJSON(jsonByte, out); isErr != nil {
		fmt.Println(isErr)
	}

}

// LegalizePayload method
func (adaptor *Adaptor) LegalizePayload(bsonPayload bson.M, out interface{}) {
	for key, value := range bsonPayload {
		fmt.Printf("Key: %v Value: %v\n", key, value)
	}
}

///////////// FILTER HANDLDER /////////////
func (adaptor *Adaptor) parserFilter(payload Payload, mongoFindOptions *options.FindOptions) {
	var payloadOptions FindOptions = payload.Options

	if payloadOptions.Limit != 0 {
		fmt.Println("limit: ", payloadOptions.Limit)
	}

	// SetLimitHandler := func(limitValue int64) {
	// 	mongoFindOptions.SetLimit(limitValue)
	// }

	// for key, value := range optionsFilter["options"].(map[string]interface{}) {
	// 	// fmt.Printf("%v %T", key, value)
	// 	if key == "limit" {
	// 		SetLimitHandler(int64(value.(float64)))
	// 	}
	// }
}
