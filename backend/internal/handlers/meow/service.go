
package Meow

import (
	"context"
	"log/slog"

	"github.com/abhikaboy/GERM-template/xutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// newService receives the map of collections and picks out Jobs
func newService(collections map[string]*mongo.Collection) *Service {
	return &Service{
		Meows: collections["meows"],
	}
}

// GetAllMeows fetches all Meow documents from MongoDB
func (s *Service) GetAllMeows() ([]MeowDocument, error) {
	ctx := context.Background()
	cursor, err := s.Meows.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []MeowDocument
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetMeowByID returns a single Meow document by its ObjectID
func (s *Service) GetMeowByID(id primitive.ObjectID) (*MeowDocument, error) {
	ctx := context.Background()
	filter := bson.M{"_id": id}

	var Meow MeowDocument
	err := s.Meows.FindOne(ctx, filter).Decode(&Meow)

	if err == mongo.ErrNoDocuments {
		// No matching Meow found
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		// Different error occurred
		return nil, err
	}

	return &Meow, nil
}

// InsertMeow adds a new Meow document
func (s *Service) CreateMeow(r *MeowDocument) (*MeowDocument, error) {
	ctx := context.Background()
	// Insert the document into the collection

	result, err := s.Meows.InsertOne(ctx, r)
	if err != nil {
		return nil, err
	}

	// Cast the inserted ID to ObjectID
	id := result.InsertedID.(primitive.ObjectID)
	r.ID = id
	slog.LogAttrs(ctx, slog.LevelInfo, "Meow inserted", slog.String("id", id.Hex()))

	return r, nil
}

// UpdatePartialMeow updates only specified fields of a Meow document by ObjectID.
func (s *Service) UpdatePartialMeow(id primitive.ObjectID, updated UpdateMeowDocument) error {
	ctx := context.Background()
	filter := bson.M{"_id": id}

	updateFields, err := xutils.ToDoc(updated)
	if err != nil {
		return err
	}

	update := bson.M{"$set": updateFields}

	_, err = s.Meows.UpdateOne(ctx, filter, update)
	return err
}

// DeleteMeow removes a Meow document by ObjectID.
func (s *Service) DeleteMeow(id primitive.ObjectID) error {
	ctx := context.Background()

	filter := bson.M{"_id": id}

	_, err := s.Meows.DeleteOne(ctx, filter)
	return err
}

