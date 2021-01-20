package gomongo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/agustadewa/gomongo/tools"
	"github.com/gin-gonic/gin"
	"gitlab.com/yosiaagustadewa/qsl-service/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Adaptor Type
type Adaptor struct {
	Client mongo.Client
	DBName string
}

// Connect method
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

// QueryUpdateDocument method
func (adaptor *Adaptor) QueryUpdateMany(ctx context.Context, collname string, filterQuery bson.M, updateQuery bson.M) error {
	var err error
	fmt.Println("filterQuery", filterQuery)
	fmt.Println("updateQuery", updateQuery)

	Collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	_, err = Collection.UpdateMany(ctx, filterQuery, updateQuery)

	return err
}

// QueryUpdateOne method
func (adaptor *Adaptor) QueryUpdateOne(ctx context.Context, collname string, updateOpt *options.UpdateOptions, filterQuery bson.M, updateQuery bson.M, result *mongo.UpdateResult) error {
	var err error
	result, err = adaptor.Client.Database(adaptor.DBName).Collection(collname).UpdateOne(ctx, filterQuery, updateQuery, updateOpt)
	if err != nil {
		return err
	}
	return nil
}

// QueryCreateCollection create collection in mongodb
func (adaptor *Adaptor) QueryCreateCollection(ctx context.Context, collname string) error {
	errCreateCollection := adaptor.Client.Database(adaptor.DBName).CreateCollection(ctx, collname)
	return errCreateCollection
}

// QueryInsert Query Insert to mongodb
func (adaptor *Adaptor) QueryInsert(ctx context.Context, collname string, byteQuery []byte) (interface{}, error) {
	var insertResult interface{}
	var errorInserting error

	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	insertResult, errorInserting = collection.InsertOne(ctx, query)

	return insertResult, errorInserting
}

// QueryInsertV2 Query Insert to mongodb
func (adaptor *Adaptor) QueryInsertV2(ctx context.Context, collname string, query interface{}, result interface{}) error {
	result, errorInserting := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		InsertOne(ctx, query)

	if errorInserting != nil {
		log.Println(errorInserting)
	}

	return errorInserting
}

// QueryInsertV2 Query Insert to mongodb
func (adaptor *Adaptor) QueryInsertV3(ctx context.Context, collname string, query interface{}) (*mongo.InsertOneResult, error) {
	result, errorInserting := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		InsertOne(ctx, query)

	if errorInserting != nil {
		log.Println(errorInserting)
	}

	return result, errorInserting
}

// QueryFind query find to mongodb
func (adaptor *Adaptor) QueryFind(ctx context.Context, collname string, byteQuery []byte) ([]byte, error) {
	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)

	var received bson.M
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	errFinding := collection.FindOne(ctx, query).Decode(&received)
	jsonBytes, _ := bson.MarshalJSON(&received)

	return jsonBytes, errFinding
}

// QueryFindV2 query find to mongodb
func (adaptor *Adaptor) QueryFindV2(ctx context.Context, collName string, findOneOptions *options.FindOneOptions, query interface{}, result interface{}) error {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collName)
	return collection.FindOne(ctx, query, findOneOptions).Decode(result)
}

// QueryFindMany query find many to mongodb
func (adaptor *Adaptor) QueryFindMany(ctx context.Context, collname string, byteQuery []byte, findOptions *options.FindOptions) ([]byte, error) {
	var query bson.M
	bson.UnmarshalJSON(byteQuery, &query)

	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	cursor, _ := collection.Find(ctx, query, findOptions)

	var received []bson.M
	var err error

	if err = cursor.All(ctx, &received); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(received)
	var results []byte
	results, err = bson.MarshalJSON(received)

	return results, err
}

// QueryFindManyV2 query find many to mongodb
func (adaptor *Adaptor) QueryFindManyV2(ctx context.Context, collname string, findOptions *options.FindOptions, query interface{}, result interface{}) error {
	collection := adaptor.Client.Database(adaptor.DBName).Collection(collname)
	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, result)
	if err != nil {
		return err
	}

	return err
}

// QueryCount query find to mongodb
func (adaptor *Adaptor) QueryCount(ctx context.Context, collname string, query bson.M) (int64, error) {
	Count, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		CountDocuments(ctx, query)

	return Count, err
}

// QueryFindAndUpdate method
func (adaptor *Adaptor) QueryFindAndUpdate(ctx context.Context, collname string, queryFilter bson.M, setQuery bson.M, setOnInsertQuery bson.M) (int64, error) {
	var err error
	var count int64

	updateQuery := bson.M{
		"$set":         setQuery,
		"$setOnInsert": setOnInsertQuery,
	}

	var updateOptions *options.FindOneAndUpdateOptions
	updateOptions.SetReturnDocument(1)
	updateOptions.SetUpsert(true)

	insertResult := adaptor.Client.Database(adaptor.DBName).Collection(collname).FindOneAndUpdate(ctx, queryFilter, updateQuery, updateOptions)
	fmt.Println(*insertResult)
	return count, err
}

// QueryFindAndUpdateV2 method
// updateQuery := bson.M{
//		"$set":         setQuery,
//		"$setOnInsert": setOnInsertQuery,
//	}
func (adaptor *Adaptor) QueryFindAndUpdateV2(ctx context.Context, collname string, findAndUpdateOpt *options.FindOneAndUpdateOptions, filterQuery interface{}, updateQuery interface{}, result interface{}) error {
	err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		FindOneAndUpdate(ctx, filterQuery, updateQuery, findAndUpdateOpt).
		Decode(result)

	return err
}

// QueryRemoveOne method
func (adaptor *Adaptor) QueryRemoveOne(ctx context.Context, collname string, queryFilter interface{}) (int64, error) {
	delResult, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		DeleteOne(ctx, queryFilter)
	if err != nil {
		return 0, err
	}

	return delResult.DeletedCount, err
}

// QueryRemoveMany method
func (adaptor *Adaptor) QueryRemoveMany(ctx context.Context, collname string, queryFilter interface{}) (int64, error) {
	delResult, err := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		DeleteMany(ctx, queryFilter)
	if err != nil {
		return 0, err

	}
	return delResult.DeletedCount, err
}

// QueryConfirm method
func (adaptor *Adaptor) QueryConfirm(ctx context.Context, collname, key, value string) bool {
	queryResult := bson.M{}
	errFindKey := adaptor.Client.
		Database(adaptor.DBName).
		Collection(collname).
		FindOne(ctx, bson.M{"key": key}).
		Decode(&queryResult)
	if errFindKey != nil {
		panic(errFindKey)
	}

	if queryResult["value"].(string) == value {
		return true
	} else {
		return false
	}
}

// /////////// PAYLOAD FILTER /////////////

// ParsePayload method
func (adaptor *Adaptor) ParsePayload(jsonByte []byte, out interface{}, c *gin.Context) {
	if isErr := bson.UnmarshalJSON(jsonByte, out); isErr != nil {
		c.JSON(400, c.Error(isErr))
		return
	}
}

// Modeling filler
func (adaptor *Adaptor) Modeling(jsonByte *[]byte, collname string) error {
	var err error

	if collname == "identity" {
		identity := models.Identity{}
		err = bson.UnmarshalJSON(*jsonByte, &identity)
		*jsonByte, err = bson.MarshalJSON(&identity)

	} else if collname == "event" {
		event := models.EventCallSign{}
		err = bson.UnmarshalJSON(*jsonByte, &event)
		*jsonByte, err = bson.MarshalJSON(&event)
	}
	return err
}

// ParseOptions method
func (adaptor *Adaptor) ParseOptions(payload models.Payload, options *options.FindOptions) {
	// LIMIT
	limitVal := payload.Options.Limit
	if limitVal > 0 {
		if limitVal >= 100 {
			options.SetLimit(100)
		} else {
			options.SetLimit(limitVal)
		}
	} else {
		options.SetLimit(100)
	}

	// SORT
	if payload.Options.Sort != nil {
		options.SetSort(payload.Options.Sort)
	}

	// SKIP
	skipVal := payload.Options.Skip
	if skipVal >= 0 {
		options.SetSkip(skipVal)
	} else {
		options.SetSkip(0)
	}

	// PROJECTION
	if payload.Options.Projection != nil {
		options.SetProjection(payload.Options.Projection)
	}
}

// ParseCertificateFormat
func (adaptor *Adaptor) ParseCertificateFormat(certificateAttribute *models.CertificateAttribute) {
	var certNumberGenerator tools.StringNumber
	certNumberGenerator.SetNDigit(4)
	certNumberGenerator.SetCounter(certificateAttribute.Number)

	certificateAttribute.Frequency = strings.ReplaceAll(certificateAttribute.Frequency, " MHz", "")
	certificateAttribute.Band = strings.ReplaceAll(certificateAttribute.Band, " M", "")
	certificateAttribute.Format = strings.ReplaceAll(certificateAttribute.Format, "#NO#", certNumberGenerator.ValString(certificateAttribute.Number))
	certificateAttribute.Format = strings.ReplaceAll(certificateAttribute.Format, "#FREQUENCY#", certificateAttribute.Frequency)
	certificateAttribute.Format = strings.ReplaceAll(certificateAttribute.Format, "#MODE#", certificateAttribute.Mode)
	certificateAttribute.Format = strings.ReplaceAll(certificateAttribute.Format, "#STATION#", certificateAttribute.Station)
	certificateAttribute.Format = strings.ReplaceAll(certificateAttribute.Format, "#BAND#", certificateAttribute.Band)
}

// SetDownloadLog
func (adaptor *Adaptor) SetDownloadLog(ctx context.Context, downloadLogData models.DownloadLog) error {
	_, errSetLog := adaptor.Client.Database(adaptor.DBName).
		Collection(models.CollCertificateDownloadLog).
		InsertOne(ctx, &downloadLogData)
	if errSetLog != nil {
		return errors.New("error inserting log: " + errSetLog.Error())
	}
	return nil
}

// // GetDownloadLog
// func (adaptor *Adaptor)GetDownloadLog(ctx context.Context, downloadLogQuery models.GetDownloadLog) error {
//	_, errGetLog := adaptor.Client.Database(adaptor.DBName).
//		Collection(models.CollCertificateDownloadLog).
//		InsertOne(ctx, &downloadLogQuery)
//	if errGetLog != nil {
//		return errors.New("error inserting log: "+errGetLog.Error())
//	}
//	return nil
// }

// GetDate method
func (adaptor *Adaptor) GetDate() string {
	Time := time.Now().UnixNano()
	dateRune := []rune(strconv.Itoa(int(Time)))
	parsedDate := string(dateRune[0:13])

	return parsedDate
}
