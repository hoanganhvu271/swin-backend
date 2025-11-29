package models

type WoodPiece struct {
	ID          string   `firestore:"id"`
	DatabaseID  string   `firestore:"database_id"`
	Name        string   `firestore:"name"`
	Description string   `firestore:"description"`
	ImageUrls   []string `firestore:"image_urls"`
}
