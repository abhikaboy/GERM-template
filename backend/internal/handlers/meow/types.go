
package Meow

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateMeowParams struct {
	Field1   string      `validate:"required" json:"field1"`
	Field2   Enumeration `validate:"required" json:"field2"`
	Picture  *string     `validate:"required" json:"picture"`
}

type MeowDocument struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Field1    string            `bson:"field1" json:"field1"`
	Field2    Enumeration       `bson:"field2" json:"field2"`
	Picture   *string           `bson:"picture" json:"picture"`
	Timestamp time.Time         `bson:"timestamp" json:"timestamp"`
}

type UpdateMeowDocument struct {
	Field1   string      `bson:"field1,omitempty" json:"field1,omitempty"`
	Field2   Enumeration `bson:"field2,omitempty" json:"field2,omitempty"`
	Picture  *string     `bson:"picture,omitempty" json:"picture,omitempty"`
}

type Enumeration string

const (
	Option1 Enumeration = "Option1"
	Option2 Enumeration = "Option2"
	Option3 Enumeration = "Option3"
)


/*
Meow Service to be used by Meow Handler to interact with the
Database layer of the application
*/

type Service struct {
	Meows *mongo.Collection
}
