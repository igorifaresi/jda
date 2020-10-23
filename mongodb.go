package jda

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func MongoConnectToDatabase(uri, database string) (*mongo.Database, error) {
	l := GetLogger()

	client, err := mongo.NewClient(
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to create mongodb client")
		return nil, l.ErrorQueue
	}

	err = client.Connect(context.TODO())
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to connect to mongodb client")
		return nil, l.ErrorQueue
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to ping mongodb database")
		return nil, l.ErrorQueue
	}

	return client.Database(database), nil
}

func MongoInsert(
	database *mongo.Database,
	s interface{},
	collectionName string,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	data, err := bson.Marshal(s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert struct interface to bson")
		return l.ErrorQueue
	}

	_, err = collection.InsertOne(context.TODO(), data)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to insert element in mongodb database")
		return l.ErrorQueue
	}

	return nil
}

func MongoGetAll(
	database *mongo.Database,
	outputArray interface{},
	collectionName string,
	optionsArg ...*options.FindOptions,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)
	
	var findOptions *options.FindOptions
	if len(optionsArg) == 0 || optionsArg[0] == nil {
		findOptions = options.Find()
	} else {
		findOptions = optionsArg[0]
	}
	findOptions.SetLimit(max)
	
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in get all elements in collection")
		return l.ErrorQueue
	}

	err = cursor.All(context.TODO(), outputArray)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to decode the data to interface")
		return l.ErrorQueue
	}

	return nil
}

func MongoFind(
	database *mongo.Database,
	s interface{},
	outputArray interface{},
	collectionName string,
	max int64,
	optionsArg ...*options.FindOptions,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	data, err := bson.Marshal(s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert struct interface to bson")
		return l.ErrorQueue
	}

	var findOptions *options.FindOptions
	if len(optionsArg) == 0 || optionsArg[0] == nil {
		findOptions = options.Find()
	} else {
		findOptions = optionsArg[0]
	}
	findOptions.SetLimit(max)

	cursor, err := collection.Find(context.TODO(), data, findOptions)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in find elements in database")
		return l.ErrorQueue
	}

	err = cursor.All(context.TODO(), outputArray)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to decode the data to interface")
		return l.ErrorQueue
	}

	return nil
}

func MongoFindOne(
	database *mongo.Database,
	s interface{},
	output interface{},
	collectionName string,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	data, err := bson.Marshal(s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert struct interface to bson")
		return l.ErrorQueue
	}

	err = collection.FindOne(context.TODO(), data).Decode(output)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in find and decode one element in database")
		return l.ErrorQueue
	}

	return nil
}

func MongoUpdateOne(
	database *mongo.Database,
	query interface{},
	toUpdate bson.M,
	collectionName string,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	queryData, err := bson.Marshal(query)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert query struct interface to bson")
		return l.ErrorQueue
	}

	toUpdateData, err := bson.Marshal(toUpdate)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert to update struct interface to bson")
		return l.ErrorQueue
	}

	_, err = collection.UpdateOne(
		context.TODO(),
		queryData,
		toUpdateData,
	)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to update struct")
		return l.ErrorQueue
	}

	return nil
}

func MongoAggregate(
	database *mongo.Database,
	query []bson.M,
	output interface{},
	collectionName string,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	cursor, err := collection.Aggregate(context.TODO(), query)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in find elements in database")
		return l.ErrorQueue
	}

	err = cursor.All(context.TODO(), output)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to decode the data to interface")
		return l.ErrorQueue
	}

	return nil
}

func MongoCount(
	database *mongo.Database,
	s interface{},
	collectionName string,
) (int64, error) {
	l := GetLogger()

	collection := database.Collection(collectionName)

	data, err := bson.Marshal(s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert struct interface to bson")
		return 0, l.ErrorQueue
	}

	qnt, err := collection.CountDocuments(context.TODO(), data)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in find and decode one element in database")
		return 0, l.ErrorQueue
	}

	return qnt, nil
}

func MongoDeleteOne(
	database *mongo.Database,
	s interface{},
	collectionName string,
) error {
	l := GetLogger()

	collection := database.Collection(collectionName)

	data, err := bson.Marshal(s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to convert struct interface to bson")
		return l.ErrorQueue
	}

	resp, err := collection.DeleteOne(context.TODO(), data)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in delete one element in database")
		return l.ErrorQueue
	}
	
	if resp.DeletedCount == 0 {
		l.Error("No one was deleted")
		return l.ErrorQueue
	}

	return nil
}
