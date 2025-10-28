package repository

import (
	"context"
	"errors"
	"time"

	"loterias-api-golang/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ResultadoRepository struct {
	collection *mongo.Collection
}

func NewResultadoRepository(db *mongo.Database) *ResultadoRepository {
	return &ResultadoRepository{
		collection: db.Collection("resultados"),
	}
}

func (r *ResultadoRepository) FindByLoteria(loteria string) ([]model.Resultado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id.loteria": loteria}
	opts := options.Find().SetSort(bson.D{{Key: "_id.concurso", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resultados []model.Resultado
	if err = cursor.All(ctx, &resultados); err != nil {
		return nil, err
	}

	for i := range resultados {
		resultados[i].AfterFind()
	}

	return resultados, nil
}

func (r *ResultadoRepository) FindByID(loteria string, concurso int) (*model.Resultado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"_id.loteria":  loteria,
		"_id.concurso": concurso,
	}

	var resultado model.Resultado
	err := r.collection.FindOne(ctx, filter).Decode(&resultado)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	resultado.AfterFind()
	return &resultado, nil
}

func (r *ResultadoRepository) FindLatest(loteria string) (*model.Resultado, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id.loteria": loteria}
	opts := options.FindOne().SetSort(bson.D{{Key: "_id.concurso", Value: -1}})

	var resultado model.Resultado
	err := r.collection.FindOne(ctx, filter, opts).Decode(&resultado)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.Resultado{}, nil
		}
		return nil, err
	}

	resultado.AfterFind()
	return &resultado, nil
}

func (r *ResultadoRepository) Save(resultado *model.Resultado) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultado.BeforeSave()

	filter := bson.M{
		"_id": resultado.ID,
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, filter, resultado, opts)
	return err
}

func (r *ResultadoRepository) SaveAll(resultados []model.Resultado) error {
	if len(resultados) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var operations []mongo.WriteModel
	for i := range resultados {
		resultados[i].BeforeSave()

		filter := bson.M{
			"_id": resultados[i].ID,
		}

		operation := mongo.NewReplaceOneModel()
		operation.SetFilter(filter)
		operation.SetReplacement(resultados[i])
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	_, err := r.collection.BulkWrite(ctx, operations)
	return err
}
