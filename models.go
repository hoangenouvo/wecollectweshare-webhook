package main

type Transactions struct {
	Description     string   `json:"description" firestore:"description"`
	Email           string   `json:"email" firestore:"email"`
	GiverName       string   `json:"giverName" firestore:"giverName"`
	PhoneNumber     string   `json:"phoneNumber" firestore:"phoneNumber"`
	VolunteerId     string   `json:"volunteer" firestore:"volunteerId"`
	Address         string   `json:"address" firestore:"address"`
	Long            float64  `json:"lng" firestore:"lng"`
	Lat             float64  `json:"lat" firestore:"lat"`
	CreatedDate     int64    `json:"created_date" firestore:"createdDate"`
	ImageURL        []string `json:"imageURL" firestore:"imageURL"`
	Status          string   `json:"status" firestore:"status"`
	TransactionTime int64    `json:"transactionTime" firestore:"transactionTime"`
}

type Person struct {
	Name string `json:"name"`
}
