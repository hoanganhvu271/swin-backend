package models

type WoodPiece struct {
	ID          string   `json:"id" firestore:"id"`
	DatabaseID  string   `json:"database_id" firestore:"database_id"`
	Name        string   `json:"name" firestore:"name"`
	Description string   `json:"description" firestore:"description"`
	ImageUrls   []string `json:"image_urls" firestore:"image_urls"`
}
