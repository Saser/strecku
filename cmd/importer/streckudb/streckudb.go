package streckudb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	cache  map[primitive.ObjectID]interface{}
}

type User struct {
	ID    primitive.ObjectID `bson:"_id"`
	Email string             `bson:"email"`
	Name  string             `bson:"name"`
}

type Store struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Range []*StoreRangeItem  `bson:"range"`
}

type StoreRangeItem struct {
	Product     primitive.ObjectID `bson:"product"`
	PriceLevels []float64          `bson:"pricelevels"`
}

type Product struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Purchase struct {
	ID      primitive.ObjectID  `bson:"_id"`
	User    primitive.ObjectID  `bson:"user"`
	Store   primitive.ObjectID  `bson:"store"`
	Price   float64             `bson:"price"`
	Note    *string             `bson:"note"`
	Amount  *int                `bson:"amount"`
	Product *primitive.ObjectID `bson:"product"`
}

func Client(ctx context.Context, username, password string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI("mongodb://localhost:27017").
			SetAuth(options.Credential{
				AuthSource: "admin",
				Username:   username,
				Password:   password,
			}),
	)
	if err != nil {
		return nil, fmt.Errorf("db client: %w", err)
	}
	return client, nil
}

func New(client *mongo.Client) *DB {
	return &DB{
		client: client,
		cache:  make(map[primitive.ObjectID]interface{}),
	}
}

func (db *DB) GetUser(ctx context.Context, id primitive.ObjectID) (*User, error) {
	if user, ok := db.cache[id]; ok {
		return user.(*User), nil
	}
	usersCollection := db.client.Database("strecku").Collection("users")
	user := new(User)
	if err := usersCollection.FindOne(ctx, bson.M{"_id": id}).Decode(user); err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	db.cache[id] = user
	return user, nil
}

func (db *DB) FindUsers(ctx context.Context) ([]*User, error) {
	wrap := func(err error) error { return fmt.Errorf("find users: %w", err) }
	usersCollection := db.client.Database("strecku").Collection("users")
	var allUsers []*User
	cur, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, wrap(err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	for cur.Next(ctx) {
		user := new(User)
		if err := cur.Decode(user); err != nil {
			return nil, wrap(err)
		}
		allUsers = append(allUsers, user)
	}
	for _, user := range allUsers {
		db.cache[user.ID] = user
	}
	return allUsers, nil
}

func (db *DB) GetStore(ctx context.Context, id primitive.ObjectID) (*Store, error) {
	if store, ok := db.cache[id]; ok {
		return store.(*Store), nil
	}
	storesCollection := db.client.Database("strecku").Collection("stores")
	store := new(Store)
	if err := storesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(store); err != nil {
		return nil, fmt.Errorf("get store: %w", err)
	}
	db.cache[id] = store
	return store, nil
}

func (db *DB) FindStores(ctx context.Context) ([]*Store, error) {
	wrap := func(err error) error { return fmt.Errorf("find stores: %w", err) }
	storesCollection := db.client.Database("strecku").Collection("stores")
	var allStores []*Store
	cur, err := storesCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, wrap(err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	for cur.Next(ctx) {
		store := new(Store)
		if err := cur.Decode(store); err != nil {
			return nil, wrap(err)
		}
		allStores = append(allStores, store)
	}
	for _, store := range allStores {
		db.cache[store.ID] = store
	}
	return allStores, nil
}

func (db *DB) GetProduct(ctx context.Context, id primitive.ObjectID) (*Product, error) {
	if product, ok := db.cache[id]; ok {
		return product.(*Product), nil
	}
	productsCollection := db.client.Database("strecku").Collection("products")
	product := new(Product)
	if err := productsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(product); err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	db.cache[id] = product
	return product, nil
}

func (db *DB) FindProducts(ctx context.Context) ([]*Product, error) {
	wrap := func(err error) error { return fmt.Errorf("find products: %w", err) }
	productsCollection := db.client.Database("strecku").Collection("products")
	var allProducts []*Product
	cur, err := productsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, wrap(err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	for cur.Next(ctx) {
		product := new(Product)
		if err := cur.Decode(product); err != nil {
			return nil, wrap(err)
		}
		allProducts = append(allProducts, product)
	}
	for _, product := range allProducts {
		db.cache[product.ID] = product
	}
	return allProducts, nil
}

func (db *DB) GetPurchase(ctx context.Context, id primitive.ObjectID) (*Purchase, error) {
	if purchase, ok := db.cache[id]; ok {
		return purchase.(*Purchase), nil
	}
	purchasesCollection := db.client.Database("strecku").Collection("purchases")
	purchase := new(Purchase)
	if err := purchasesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(purchase); err != nil {
		return nil, fmt.Errorf("get purchase: %w", err)
	}
	db.cache[id] = purchase
	return purchase, nil
}

func (db *DB) FindPurchases(ctx context.Context) ([]*Purchase, error) {
	wrap := func(err error) error { return fmt.Errorf("find purchases: %w", err) }
	purchasesCollection := db.client.Database("strecku").Collection("purchases")
	var allPurchases []*Purchase
	cur, err := purchasesCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, wrap(err)
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Print(err)
		}
	}()
	for cur.Next(ctx) {
		purchase := new(Purchase)
		if err := cur.Decode(purchase); err != nil {
			return nil, wrap(err)
		}
		allPurchases = append(allPurchases, purchase)
	}
	for _, purchase := range allPurchases {
		db.cache[purchase.ID] = purchase
	}
	return allPurchases, nil
}
