package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/olafser/ambulance-management-api-app/internal/config"
	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrVehicleNotFound = errors.New("vehicle not found")
var ErrVehicleConflict = errors.New("vehicle already exists")

type VehicleRepository interface {
	List(ctx context.Context, status, station string) ([]entity.VehicleEntity, error)
	Create(ctx context.Context, vehicle entity.VehicleEntity) (entity.VehicleEntity, error)
	GetByID(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error)
	UpdateByID(ctx context.Context, vehicleID int64, vehicle entity.VehicleEntity) (entity.VehicleEntity, error)
	UpdateStatusByID(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error)
	DeleteByID(ctx context.Context, vehicleID int64) error
}

type repositoryVehicle struct {
	vehicles *mongo.Collection
	counters *mongo.Collection
}

func NewVehicleRepository(db *mongo.Database, cfg config.MongoConfig) (VehicleRepository, error) {
	repo := &repositoryVehicle{
		vehicles: db.Collection(cfg.VehiclesColl),
		counters: db.Collection(cfg.CountersColl),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *repositoryVehicle) ensureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "vehicleId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "callSign", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "plateNumber", Value: 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := r.vehicles.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}
	return nil
}

func (r *repositoryVehicle) List(ctx context.Context, status, station string) ([]entity.VehicleEntity, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if station != "" {
		filter["station"] = station
	}

	cursor, err := r.vehicles.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "vehicleId", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]entity.VehicleEntity, 0)
	for cursor.Next(ctx) {
		var item entity.VehicleEntity
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repositoryVehicle) Create(ctx context.Context, vehicle entity.VehicleEntity) (entity.VehicleEntity, error) {
	if vehicle.VehicleID == 0 {
		id, err := r.nextVehicleID(ctx)
		if err != nil {
			return entity.VehicleEntity{}, err
		}
		vehicle.VehicleID = id
	}

	_, err := r.vehicles.InsertOne(ctx, vehicle)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.VehicleEntity{}, ErrVehicleConflict
		}
		return entity.VehicleEntity{}, err
	}

	return vehicle, nil
}

func (r *repositoryVehicle) GetByID(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
	var item entity.VehicleEntity
	err := r.vehicles.FindOne(ctx, bson.M{"vehicleId": vehicleID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return entity.VehicleEntity{}, ErrVehicleNotFound
	}
	if err != nil {
		return entity.VehicleEntity{}, err
	}

	return item, nil
}

func (r *repositoryVehicle) UpdateByID(ctx context.Context, vehicleID int64, vehicle entity.VehicleEntity) (entity.VehicleEntity, error) {
	vehicle.VehicleID = vehicleID

	result, err := r.vehicles.ReplaceOne(ctx, bson.M{"vehicleId": vehicleID}, vehicle)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.VehicleEntity{}, ErrVehicleConflict
		}
		return entity.VehicleEntity{}, err
	}
	if result.MatchedCount == 0 {
		return entity.VehicleEntity{}, ErrVehicleNotFound
	}

	return vehicle, nil
}

func (r *repositoryVehicle) UpdateStatusByID(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
	result, err := r.vehicles.UpdateOne(ctx, bson.M{"vehicleId": vehicleID}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return entity.VehicleEntity{}, err
	}
	if result.MatchedCount == 0 {
		return entity.VehicleEntity{}, ErrVehicleNotFound
	}

	return r.GetByID(ctx, vehicleID)
}

func (r *repositoryVehicle) DeleteByID(ctx context.Context, vehicleID int64) error {
	result, err := r.vehicles.DeleteOne(ctx, bson.M{"vehicleId": vehicleID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrVehicleNotFound
	}
	return nil
}

func (r *repositoryVehicle) nextVehicleID(ctx context.Context) (int64, error) {
	type counterDocument struct {
		ID  string `bson:"_id"`
		Seq int64  `bson:"seq"`
	}

	var result counterDocument
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := r.counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "vehicles"},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}
