package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/config"
	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrDispatchNotFound = errors.New("dispatch not found")
var ErrDispatchConflict = errors.New("dispatch already exists")

type DispatchRepository interface {
	List(ctx context.Context, status, city string) ([]entity.DispatchEntity, error)
	Create(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error)
	GetByID(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error)
	UpdateByID(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error)
	UpdateStatusByID(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error)
	DeleteUnfinishedByVehicleCallSign(ctx context.Context, callSign string) error
	DeleteByID(ctx context.Context, dispatchID int64) error
}

type repositoryDispatch struct {
	dispatches *mongo.Collection
	counters   *mongo.Collection
}

func NewDispatchRepository(db *mongo.Database, cfg config.MongoConfig) (DispatchRepository, error) {
	repo := &repositoryDispatch{
		dispatches: db.Collection(cfg.DispatchesColl),
		counters:   db.Collection(cfg.CountersColl),
	}

	if err := repo.ensureIndexes(context.Background()); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *repositoryDispatch) ensureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "dispatchId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "incidentNumber", Value: 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err := r.dispatches.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}
	return nil
}

func (r *repositoryDispatch) List(ctx context.Context, status, city string) ([]entity.DispatchEntity, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if city != "" {
		filter["city"] = city
	}

	cursor, err := r.dispatches.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "dispatchId", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()

	items := make([]entity.DispatchEntity, 0)
	for cursor.Next(ctx) {
		var item entity.DispatchEntity
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

func (r *repositoryDispatch) Create(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
	if dispatch.DispatchID == 0 {
		id, err := r.nextDispatchID(ctx)
		if err != nil {
			return entity.DispatchEntity{}, err
		}
		dispatch.DispatchID = id
	}

	_, err := r.dispatches.InsertOne(ctx, dispatch)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.DispatchEntity{}, ErrDispatchConflict
		}
		return entity.DispatchEntity{}, err
	}

	return dispatch, nil
}

func (r *repositoryDispatch) GetByID(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
	var item entity.DispatchEntity
	err := r.dispatches.FindOne(ctx, bson.M{"dispatchId": dispatchID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return entity.DispatchEntity{}, ErrDispatchNotFound
	}
	if err != nil {
		return entity.DispatchEntity{}, err
	}

	return item, nil
}

func (r *repositoryDispatch) UpdateByID(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
	dispatch.DispatchID = dispatchID

	result, err := r.dispatches.ReplaceOne(ctx, bson.M{"dispatchId": dispatchID}, dispatch)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.DispatchEntity{}, ErrDispatchConflict
		}
		return entity.DispatchEntity{}, err
	}
	if result.MatchedCount == 0 {
		return entity.DispatchEntity{}, ErrDispatchNotFound
	}

	return dispatch, nil
}

func (r *repositoryDispatch) UpdateStatusByID(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
	result, err := r.dispatches.UpdateOne(
		ctx,
		bson.M{"dispatchId": dispatchID},
		bson.M{"$set": bson.M{"status": status, "updatedAt": updatedAt}},
	)
	if err != nil {
		return entity.DispatchEntity{}, err
	}
	if result.MatchedCount == 0 {
		return entity.DispatchEntity{}, ErrDispatchNotFound
	}

	return r.GetByID(ctx, dispatchID)
}

func (r *repositoryDispatch) DeleteUnfinishedByVehicleCallSign(ctx context.Context, callSign string) error {
	_, err := r.dispatches.DeleteMany(ctx, bson.M{
		"ambulanceCallSign": callSign,
		"status":            bson.M{"$ne": "COMPLETED"},
	})
	return err
}

func (r *repositoryDispatch) DeleteByID(ctx context.Context, dispatchID int64) error {
	result, err := r.dispatches.DeleteOne(ctx, bson.M{"dispatchId": dispatchID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrDispatchNotFound
	}
	return nil
}

func (r *repositoryDispatch) nextDispatchID(ctx context.Context) (int64, error) {
	type counterDocument struct {
		ID  string `bson:"_id"`
		Seq int64  `bson:"seq"`
	}

	var result counterDocument
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	err := r.counters.FindOneAndUpdate(
		ctx,
		bson.M{"_id": "dispatches"},
		bson.M{"$inc": bson.M{"seq": 1}},
		opts,
	).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Seq, nil
}
