package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap/zapcore"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection      *mongo.Collection
	AuctionDuration time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.endAuction(ctx, auctionEntity)

	return nil
}

func (ar *AuctionRepository) endAuction(ctx context.Context, auction *auction_entity.Auction) {

	select {
	case <-time.After(getAuctionDuration()):
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		filter := bson.M{"_id": auction.Id}
		_, err := ar.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error("Error trying to update auction status to completed", err)
			return
		}
		logger.Info("Auction completed after interval =>", zapcore.Field{Key: "auction_id", Type: zapcore.StringType, String: auction.Id})
	}

}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		logger.Error("Error trying to parse auction duration", err)
		return time.Minute * 5
	}

	return duration
}
