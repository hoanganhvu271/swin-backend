package models

type WoodDatabase struct {
	ID          string `firestore:"id"`
	Title       string `firestore:"title"`
	Size        int    `firestore:"size"`
	Description string `firestore:"description"`
	Image       string `firestore:"image"`
}
