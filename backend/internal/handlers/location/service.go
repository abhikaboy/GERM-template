
package location

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
        locations: collections["locations"],
    }
}

// GetAlllocations fetches all location documents from MongoDB
func (s *Service) GetAlllocations() ([]locationDocument, error) {
    ctx := context.Background()
    cursor, err := s.locations.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var results []locationDocument
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}

// GetlocationByID returns a single location document by its ObjectID
func (s *Service) GetlocationByID(id primitive.ObjectID) (*locationDocument, error) {
    ctx := context.Background()
    filter := bson.M{"_id": id}

    var location locationDocument
    err := s.locations.FindOne(ctx, filter).Decode(&location)

    if err == mongo.ErrNoDocuments {
        return nil, mongo.ErrNoDocuments
    } else if err != nil {
        return nil, err
    }

    return &location, nil
}

// Insertlocation adds a new location document
func (s *Service) Insertlocation(r locationDocument) (*locationDocument, error) {
    ctx := context.Background()
    result, err := s.locations.InsertOne(ctx, r)
    if err != nil {
        return nil, err
    }

    id := result.InsertedID.(primitive.ObjectID)
    r.ID = id
    slog.LogAttrs(ctx, slog.LevelInfo, "location inserted", slog.String("id", id.Hex()))

    return &r, nil
}

// UpdatePartiallocation updates only specified fields of a location document by ObjectID.
func (s *Service) UpdatePartiallocation(id primitive.ObjectID, updated UpdatelocationDocument) error {
    ctx := context.Background()
    filter := bson.M{"_id": id}

    updateFields, err := xutils.ToDoc(updated)
    if err != nil {
        return err
    }

    update := bson.M{"$set": updateFields}

    _, err = s.locations.UpdateOne(ctx, filter, update)
    return err
}

// Deletelocation removes a location document by ObjectID.
func (s *Service) Deletelocation(id primitive.ObjectID) error {
    ctx := context.Background()
    filter := bson.M{"_id": id}
    _, err := s.locations.DeleteOne(ctx, filter)
    return err
}

// GetNearbylocations fetches location documents within a radius of a location
func (s *Service) GetNearbylocations(location []float64, radius float64) ([]locationDocument, error) {
    ctx := context.Background()
    filter := bson.M{
        "location": bson.M{
            "$near": location,
            "$maxDistance": radius,
        },
    }
    cursor, err := s.locations.Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    defer cursor.Close(ctx)

    var results []locationDocument
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }

    return results, nil
}
