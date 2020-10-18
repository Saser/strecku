package mapper

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mapper struct {
	userIDs   map[string]primitive.ObjectID
	userNames map[primitive.ObjectID]string

	storeIDs   map[string]primitive.ObjectID
	storeNames map[primitive.ObjectID]string

	productSPs   map[string]StoreProduct
	productNames map[StoreProduct]string

	purchaseIDs   map[string]primitive.ObjectID
	purchaseNames map[primitive.ObjectID]string

	paymentIDs   map[string]primitive.ObjectID
	paymentNames map[primitive.ObjectID]string
}

type StoreProduct struct {
	Store, Product primitive.ObjectID
}

func New() *Mapper {
	return &Mapper{
		userIDs:   make(map[string]primitive.ObjectID),
		userNames: make(map[primitive.ObjectID]string),

		storeIDs:   make(map[string]primitive.ObjectID),
		storeNames: make(map[primitive.ObjectID]string),

		productSPs:   make(map[string]StoreProduct),
		productNames: make(map[StoreProduct]string),

		purchaseIDs:   make(map[string]primitive.ObjectID),
		purchaseNames: make(map[primitive.ObjectID]string),

		paymentIDs:   make(map[string]primitive.ObjectID),
		paymentNames: make(map[primitive.ObjectID]string),
	}
}

func (m *Mapper) User(id primitive.ObjectID, name string) {
	m.userIDs[name] = id
	m.userNames[id] = name
}

func (m *Mapper) UserID(userName string) (primitive.ObjectID, bool) {
	id, ok := m.userIDs[userName]
	return id, ok
}

func (m *Mapper) UserName(userID primitive.ObjectID) (string, bool) {
	name, ok := m.userNames[userID]
	return name, ok
}

func (m *Mapper) Store(id primitive.ObjectID, name string) {
	m.storeIDs[name] = id
	m.storeNames[id] = name
}

func (m *Mapper) StoreID(storeName string) (primitive.ObjectID, bool) {
	id, ok := m.storeIDs[storeName]
	return id, ok
}

func (m *Mapper) StoreName(storeID primitive.ObjectID) (string, bool) {
	name, ok := m.storeNames[storeID]
	return name, ok
}

func (m *Mapper) Product(sp StoreProduct, name string) {
	m.productSPs[name] = sp
	m.productNames[sp] = name
}

func (m *Mapper) ProductSP(productName string) (StoreProduct, bool) {
	sp, ok := m.productSPs[productName]
	return sp, ok
}

func (m *Mapper) ProductName(sp StoreProduct) (string, bool) {
	name, ok := m.productNames[sp]
	return name, ok
}

func (m *Mapper) Purchase(id primitive.ObjectID, name string) {
	m.purchaseIDs[name] = id
	m.purchaseNames[id] = name
}

func (m *Mapper) PurchaseID(purchaseName string) (primitive.ObjectID, bool) {
	id, ok := m.purchaseIDs[purchaseName]
	return id, ok
}

func (m *Mapper) PurchaseName(purchaseID primitive.ObjectID) (string, bool) {
	name, ok := m.purchaseNames[purchaseID]
	return name, ok
}

func (m *Mapper) Payment(id primitive.ObjectID, name string) {
	m.paymentIDs[name] = id
	m.paymentNames[id] = name
}

func (m *Mapper) PaymentID(paymentName string) (primitive.ObjectID, bool) {
	id, ok := m.paymentIDs[paymentName]
	return id, ok
}

func (m *Mapper) PaymentName(paymentID primitive.ObjectID) (string, bool) {
	name, ok := m.paymentNames[paymentID]
	return name, ok
}
