package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"

	"backend/config"
)

func getClient(ctx context.Context) *firestore.Client {
	client, err := config.FirebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	return client
}

// -------------------- WRITE --------------------

// SetDocument ghi dữ liệu vào document (overwrite nếu tồn tại)
func SetDocument(collection, docID string, data interface{}) error {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	_, err := client.Collection(collection).Doc(docID).Set(ctx, data)
	return err
}

// AddDocument thêm document mới với ID tự sinh
func AddDocument(collection string, data interface{}) (string, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	docRef, _, err := client.Collection(collection).Add(ctx, data)
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// -------------------- READ --------------------

// GetDocument đọc document theo ID
func GetDocument(collection, docID string) (map[string]interface{}, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	doc, err := client.Collection(collection).Doc(docID).Get(ctx)
	if err != nil {
		return nil, err
	}
	return doc.Data(), nil
}

// GetCollection đọc toàn bộ document trong collection
func GetCollection(collection string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	iter := client.Collection(collection).Documents(ctx)
	defer iter.Stop()

	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		results = append(results, doc.Data())
	}
	return results, nil
}

// DeleteDocument xóa document theo collection và docID
func DeleteDocument(collection, docID string) error {
	ctx := context.Background()
	client := getClient(ctx)
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}(client)

	_, err := client.Collection(collection).Doc(docID).Delete(ctx)
	return err
}
