package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/cmd/importer/mapper"
	"github.com/Saser/strecku/cmd/importer/streckudb"
	"github.com/Saser/strecku/resources/products"
	"github.com/Saser/strecku/resources/purchases"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
)

var (
	username = flag.String("username", "", "Username in Mongo database.")
	password = flag.String("password", "", "Password in Mongo database.")
)

func priceCents(price float64) int64 {
	return int64(math.Round(price * -100))
}

func imain() int {
	flag.Parse()
	if *username == "" {
		log.Print("flag -username is missing")
		return 1
	}
	if *password == "" {
		log.Print("flag -password is missing")
		return 1
	}

	ctx := context.Background()
	client, err := streckudb.Client(ctx, *username, *password)
	if err != nil {
		log.Print(err)
		return 1
	}
	db := streckudb.New(client)
	if err != nil {
		log.Print(err)
		return 1
	}
	m := mapper.New()

	log.Print("finding all users...")
	allUsers, err := db.FindUsers(ctx)
	if err != nil {
		log.Print(err)
		return 1
	}
	log.Print("found all users.")
	log.Print("finding all stores...")
	allStores, err := db.FindStores(ctx)
	if err != nil {
		log.Print(err)
		return 1
	}
	log.Print("found all stores.")
	log.Print("finding all purchases...")
	allPurchases, err := db.FindPurchases(ctx)
	if err != nil {
		log.Print(err)
		return 1
	}
	log.Print("found all purchases.")

	for _, user := range allUsers {
		m.User(user.ID, users.GenerateName())
	}
	for _, store := range allStores {
		m.Store(store.ID, stores.GenerateName())
	}
	for _, purchase := range allPurchases {
		m.Purchase(purchase.ID, purchases.GenerateName())
	}
	// We have to gather combinations of stores and products. There are two
	// kinds of products we care about:
	// 1. Products which are currently available in a store; and
	// 2. Products which are _not_ currently available in a store, but have been
	//    purchased in that store earlier.
	log.Print("mapping all products to names...")
	storeProducts := make(map[mapper.StoreProduct]bool)
	for _, store := range allStores {
		for _, item := range store.Range {
			sp := mapper.StoreProduct{
				Store:   store.ID,
				Product: item.Product,
			}
			storeProducts[sp] = true
		}
	}
	for _, purchase := range allPurchases {
		if purchase.Product == nil {
			continue
		}
		sp := mapper.StoreProduct{
			Store:   purchase.Store,
			Product: *purchase.Product,
		}
		storeProducts[sp] = true
	}
	for sp := range storeProducts {
		storeName, ok := m.StoreName(sp.Store)
		if !ok {
			log.Printf("store not found in mapper: %v", sp.Store)
			return 1
		}
		productName := products.GenerateName(storeName)
		m.Product(sp, productName)
	}
	log.Print("mapped all products to names.")

	// Now all mappings ObjectID <-> resource name exist, so we can start
	// converting the database entities into API resources.
	log.Print("converting all users to API resources...")
	apiUsers := make([]*pb.User, len(allUsers))
	for i, user := range allUsers {
		name, ok := m.UserName(user.ID)
		if !ok {
			log.Printf("user not found in mapper: %v", user.ID)
			return 1
		}
		apiUser := &pb.User{
			Name:         name,
			EmailAddress: user.Email,
			DisplayName:  user.Name,
		}
		if err := users.Validate(apiUser); err != nil {
			log.Print(err)
			return 1
		}
		apiUsers[i] = apiUser
	}
	log.Print("converted all users to API resources.")
	log.Print("converting all stores to API resources...")
	apiStores := make([]*pb.Store, len(allStores))
	for i, store := range allStores {
		name, ok := m.StoreName(store.ID)
		if !ok {
			log.Printf("store not found in mapper: %v", store.ID)
			return 1
		}
		apiStore := &pb.Store{
			Name:        name,
			DisplayName: store.Name,
		}
		if err := stores.Validate(apiStore); err != nil {
			log.Print(err)
			return 1
		}
		apiStores[i] = apiStore
	}
	log.Print("converted all stores to API resources.")
	log.Print("converting all products to API resources...")
	apiProducts := make(map[string]*pb.Product) // name -> product
	// First, loop through all stores and create API resources for the existing
	// products, using the information in the range items to set price levels.
	for _, store := range allStores {
		for _, item := range store.Range {
			var fullPrice, discountPrice float64
			switch n := len(item.PriceLevels); n {
			case 1:
				fullPrice = item.PriceLevels[0]
				discountPrice = fullPrice
			case 2:
				fullPrice = item.PriceLevels[1]
				discountPrice = item.PriceLevels[0]
			default:
				log.Printf("store %v, product %v has %v price levels, expected 1 or 2", store, item.Product, n)
				return 1
			}
			if fullPrice < discountPrice {
				log.Printf("store %v, product %v has fullPrice (%v) < discountPrice (%v)", store, item.Product, fullPrice, discountPrice)
				return 1
			}
			product, err := db.GetProduct(ctx, item.Product)
			if err != nil {
				log.Print(err)
				return 1
			}
			sp := mapper.StoreProduct{
				Store:   store.ID,
				Product: item.Product,
			}
			productName, ok := m.ProductName(sp)
			if !ok {
				log.Printf("product not found in mapper: %v", sp)
				return 1
			}
			storeName, ok := m.StoreName(store.ID)
			if !ok {
				log.Printf("store not found in mapper: %v", store.ID)
				return 1
			}
			apiProduct := &pb.Product{
				Name:               productName,
				Parent:             storeName,
				DisplayName:        product.Name,
				FullPriceCents:     priceCents(fullPrice),
				DiscountPriceCents: priceCents(discountPrice),
			}
			if err := products.Validate(apiProduct); err != nil {
				log.Print(err)
				return 1
			}
			apiProducts[productName] = apiProduct
		}
	}
	// Then, loop through all purchases and create products for the products
	// that were not created in the previous loop. For these, we cannot know
	// their price levels, so we just set both their full price and discount
	// price to the first price we encounter for them.
	for _, purchase := range allPurchases {
		if purchase.Product == nil {
			continue
		}
		if purchase.Price < 0 {
			continue
		}
		sp := mapper.StoreProduct{
			Store:   purchase.Store,
			Product: *purchase.Product,
		}
		productName, ok := m.ProductName(sp)
		if !ok {
			log.Printf("product not found in mapper: %v", sp)
			return 1
		}
		if _, ok := apiProducts[productName]; ok {
			continue
		}
		storeName, ok := m.StoreName(purchase.Store)
		if !ok {
			log.Printf("store not found in mapper: %v", purchase.Store)
			return 1
		}
		product, err := db.GetProduct(ctx, *purchase.Product)
		if err != nil {
			log.Print(err)
			return 1
		}
		fullPriceCents := priceCents(purchase.Price)
		apiProduct := &pb.Product{
			Name:               productName,
			Parent:             storeName,
			DisplayName:        product.Name,
			FullPriceCents:     fullPriceCents,
			DiscountPriceCents: fullPriceCents,
		}
		if err := products.Validate(apiProduct); err != nil {
			log.Print(err)
			return 1
		}
		apiProducts[productName] = apiProduct
	}
	log.Print("converted all products to API resources.")
	log.Print("converting all purchases to API resources...")
	apiPurchases := make([]*pb.Purchase, 0, len(allPurchases))
	for _, purchase := range allPurchases {
		if purchase.Note == nil && purchase.Product == nil {
			log.Printf("purchase missing both note and product: %+v", purchase)
			return 1
		}
		if purchase.Amount != nil && *purchase.Amount != 1 {
			log.Printf("purchase with existing amount not equal to 1: %+v", purchase)
			return 1
		}
		if purchase.Price < 0 {
			// Negative price, which is a payment, which the API does not
			// currently support.
			continue
		}
		line := new(pb.Purchase_Line)
		if purchase.Note != nil {
			line.Description = *purchase.Note
		} else {
			product, err := db.GetProduct(ctx, *purchase.Product)
			if err != nil {
				log.Print(err)
				return 1
			}
			line.Description = product.Name
		}
		line.Quantity = 1
		line.PriceCents = priceCents(purchase.Price)
		if purchase.Product != nil {
			sp := mapper.StoreProduct{
				Store:   purchase.Store,
				Product: *purchase.Product,
			}
			productName, ok := m.ProductName(sp)
			if !ok {
				log.Printf("product not found in mapper: %v", sp)
				return 1
			}
			line.Product = productName
		}
		purchaseName, ok := m.PurchaseName(purchase.ID)
		if !ok {
			log.Printf("purchase not found in mapper: %v", purchase.ID)
			return 1
		}
		userName, ok := m.UserName(purchase.User)
		if !ok {
			log.Printf("user not found in mapper: %v", purchase.User)
			return 1
		}
		storeName, ok := m.StoreName(purchase.Store)
		if !ok {
			log.Printf("store not found in mapper: %v", purchase.Store)
			return 1
		}
		apiPurchase := &pb.Purchase{
			Name:  purchaseName,
			User:  userName,
			Store: storeName,
			Lines: []*pb.Purchase_Line{line},
		}
		if err := purchases.Validate(apiPurchase); err != nil {
			log.Print(err)
			return 1
		}
		apiPurchases = append(apiPurchases, apiPurchase)
	}
	log.Print("converted all purchases to API resources.")

	userRepo := users.NewRepository()
	for _, user := range apiUsers {
		password := fmt.Sprintf("placeholder password for user %q", user.EmailAddress)
		if err := userRepo.CreateUser(ctx, user, password); err != nil {
			log.Print(err)
			return 1
		}
	}
	storeRepo := stores.NewRepository()
	for _, store := range apiStores {
		if err := storeRepo.CreateStore(ctx, store); err != nil {
			log.Print(err)
			return 1
		}
	}
	productRepo := products.NewRepository()
	for _, product := range apiProducts {
		if err := productRepo.CreateProduct(ctx, product); err != nil {
			log.Print(err)
			return 1
		}
	}
	purchaseRepo := purchases.NewRepository()
	for i, purchase := range apiPurchases {
		if purchase == nil {
			log.Printf("purchase %v = nil", i)
			return 1
		}
		if err := purchaseRepo.CreatePurchase(ctx, purchase); err != nil {
			log.Print(err)
			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(imain())
}
