package gomongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Identity type
type Identity struct {
	CertificateID     int32  `json:"certificate_id,omitempty"`
	CertificateNumber string `json:"certificate_number,omitempty"`
	Date              string `json:"date,omitempty"`
	Time              string `json:"time,omitempty"`
	BandUnit          string `json:"band_unit,omitempty"`
	Band              int32  `json:"band,omitempty"`
	FrequencyUnit     string `json:"frequency_unit,omitempty"`
	Frequency         int32  `json:"frequency,omitempty"`
	CallSign          string `json:"call_sign,omitempty"`
	Name              string `json:"name,omitempty"`
	EventID           int32  `json:"event_id,omitempty"`
	DateCreated       string `json:"date_created,omitempty"`
	CreatedBy         string `json:"created_by,omitempty"`
	DateModified      string `json:"date_modified,omitempty"`
	ModifiedBy        string `json:"modified_by,omitempty"`
	DownloadCount     int32  `json:"download_count,omitempty"`
	CityID            int32  `json:"city_id,omitempty"`
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

//GetIdentity method
func (adaptor *Adaptor) GetIdentity(ctx context.Context, name string) bson.M {
	identity := bson.M{}

	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(identity)

	if err != nil {
		fmt.Println(err)
	}

	return identity
}

//GetIdentities method
func (adaptor *Adaptor) GetIdentities(ctx context.Context, name string) []bson.M {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(adaptor.CollName)
	cursor, err := collection.Find(ctx, bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	result := []bson.M{}

	for cursor.Next(ctx) {
		received := bson.M{}
		err := cursor.Decode(&received)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, received)
	}
	return result
}

//DeletePost method
func (adaptor *Adaptor) DeletePost(ctx context.Context, name string) {
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
