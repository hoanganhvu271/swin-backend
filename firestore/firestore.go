package firestore

import (
	"backend/models"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

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

// CreateDocument tạo document mới, trả lỗi nếu đã tồn tại
func CreateDocument(collection, docID string, data interface{}) error {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	docRef := client.Collection(collection).Doc(docID)

	// Kiểm tra document đã tồn tại chưa
	_, err := docRef.Get(ctx)
	if err == nil {
		return fmt.Errorf("document with ID '%s' already exists", docID)
	}

	_, err = docRef.Set(ctx, data)
	return err
}

// UpdateDocument cập nhật document, trả lỗi nếu không tồn tại
func UpdateDocument(collection, docID string, data interface{}) error {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	docRef := client.Collection(collection).Doc(docID)

	// Kiểm tra document có tồn tại không
	_, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("document with ID '%s' not found", docID)
	}

	_, err = docRef.Set(ctx, data)
	return err
}

// SetDocument ghi dữ liệu vào document (overwrite nếu tồn tại) - giữ lại cho backward compatibility
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

// DocumentExists kiểm tra document có tồn tại không
func DocumentExists(collection, docID string) (bool, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	_, err := client.Collection(collection).Doc(docID).Get(ctx)
	if err != nil {
		if err == iterator.Done || isNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, doc.Data())
	}
	return results, nil
}

// PaginationParams chứa các tham số phân trang
type PaginationParams struct {
	Limit      int
	Offset     int
	OrderBy    string
	Descending bool
}

// PaginatedResult chứa kết quả phân trang
type PaginatedResult struct {
	Data    []map[string]interface{} `json:"data"`
	Total   int                      `json:"total"`
	Limit   int                      `json:"limit"`
	Offset  int                      `json:"offset"`
	HasMore bool                     `json:"has_more"`
}

// GetCollectionPaginated đọc collection với phân trang
func GetCollectionPaginated(collection string, params PaginationParams) (*PaginatedResult, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	collRef := client.Collection(collection)

	// Đếm tổng số documents
	allDocs := collRef.Documents(ctx)
	total := 0
	for {
		_, err := allDocs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		total++
	}
	allDocs.Stop()

	// Query với phân trang
	var query firestore.Query
	if params.OrderBy != "" {
		if params.Descending {
			query = collRef.OrderBy(params.OrderBy, firestore.Desc)
		} else {
			query = collRef.OrderBy(params.OrderBy, firestore.Asc)
		}
	} else {
		query = collRef.Query
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, doc.Data())
	}

	return &PaginatedResult{
		Data:    results,
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(results) < total,
	}, nil
}

// GetCollectionWithFilter đọc collection với filter và phân trang
func GetCollectionWithFilter(collection, filterField, filterValue string, params PaginationParams) (*PaginatedResult, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	collRef := client.Collection(collection)

	// Query với filter
	baseQuery := collRef.Where(filterField, "==", filterValue)

	// Đếm tổng số documents với filter
	allDocs := baseQuery.Documents(ctx)
	total := 0
	for {
		_, err := allDocs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		total++
	}
	allDocs.Stop()

	// Query với phân trang
	var query firestore.Query = baseQuery
	if params.OrderBy != "" {
		if params.Descending {
			query = query.OrderBy(params.OrderBy, firestore.Desc)
		} else {
			query = query.OrderBy(params.OrderBy, firestore.Asc)
		}
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var results []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, doc.Data())
	}

	return &PaginatedResult{
		Data:    results,
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(results) < total,
	}, nil
}

// DeleteDocument xóa document theo collection và docID
func DeleteDocument(collection, docID string) error {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	_, err := client.Collection(collection).Doc(docID).Delete(ctx)
	return err
}

func isNotFoundError(err error) bool {
	return err != nil && (err.Error() == "rpc error: code = NotFound desc = Document not found" ||
		err.Error() == "rpc error: code = NotFound desc = No document to update")
}

// GetDocumentsByField lấy tất cả documents của collection theo field = value
func GetDocumentsByField(collection, field string, value interface{}) ([]models.WoodPiece, error) {
	ctx := context.Background()
	client := getClient(ctx)
	defer client.Close()

	iter := client.Collection(collection).Where(field, "==", value).Documents(ctx)
	defer iter.Stop()

	var results []models.WoodPiece
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		var p models.WoodPiece
		if err := doc.DataTo(&p); err != nil {
			return nil, err
		}
		results = append(results, p)
	}

	return results, nil
}
