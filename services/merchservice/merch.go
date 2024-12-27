package merchservice

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ut-sama-art-studio/art-market-backend/database"
	"github.com/ut-sama-art-studio/art-market-backend/graph/model"
)

// MerchItem represents the database object
type MerchItem struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       float64   `json:"price"`
	Inventory   *int      `json:"inventory"`
	Type        string    `json:"type"`
	Height      *float64  `json:"height"`
	Width       *float64  `json:"width"`
	Unit        *string   `json:"unit"`
	ImageURLs   []*string `json:"images"` // Stores up to 5 image URLs, order matters
	Timestamp   string    `json:"timestamp"`
}

var (
	ALL_COLUMNS = `id, owner_id, name, description, price, inventory, type, height, width, unit, image_url1, image_url2, image_url3, image_url4, image_url5, timestamp`
)

func (merch *MerchItem) ToGraphqlMerchItem() *model.MerchItem {
	return &model.MerchItem{
		ID:          merch.ID,
		OwnerID:     merch.OwnerID,
		Name:        merch.Name,
		Description: merch.Description,
		Price:       merch.Price,
		Inventory:   merch.Inventory,
		Type:        merch.Type,
		Width:       merch.Width,
		Height:      merch.Height,
		Images:      filterListNil(merch.ImageURLs),
		Unit:        merch.Unit,
	}
}

func filterListNil(list []*string) []*string {
	i := 0
	for _, item := range list {
		if item != nil {
			list[i] = item
			i++
		}
	}
	// Slice to only include non-nil elements
	return list[:i]
}

func scanMerchItemRow(scanner interface {
	Scan(dest ...interface{}) error
}) (*MerchItem, error) {
	var item MerchItem
	item.ImageURLs = make([]*string, 5)

	err := scanner.Scan(
		&item.ID, &item.OwnerID, &item.Name, &item.Description, &item.Price, &item.Inventory,
		&item.Type, &item.Height, &item.Width, &item.Unit,
		&item.ImageURLs[0], &item.ImageURLs[1], &item.ImageURLs[2], &item.ImageURLs[3], &item.ImageURLs[4],
		&item.Timestamp,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Scans a single row
func scanMerchItem(row *sql.Row) (*MerchItem, error) {
	item, err := scanMerchItemRow(row)
	if err == sql.ErrNoRows {
		log.Println("No merch item found:", err)
		return nil, err
	} else if err != nil {
		log.Println("Error scanning single merch item:", err)
		return nil, err
	}
	return item, nil
}

// Scans multiple rows
func scanMerchItems(rows *sql.Rows) ([]MerchItem, error) {
	var items []MerchItem

	defer rows.Close()
	for rows.Next() {
		item, err := scanMerchItemRow(rows)
		if err != nil {
			log.Println("Error scanning merch item row:", err)
			return nil, err
		}
		items = append(items, *item)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through rows:", err)
		return nil, err
	}

	return items, nil
}

// Create inserts a new MerchItem into the database
func (item *MerchItem) Create() (string, error) {
	query := `
		INSERT INTO "MerchItem" (id, owner_id, name, description, price, inventory, type, height, width, unit, image_url1, image_url2, image_url3, image_url4, image_url5)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`
	var id string
	err := database.Db.QueryRow(
		query, item.ID, item.OwnerID, item.Name, item.Description, item.Price, item.Inventory,
		item.Type, item.Height, item.Width, item.Unit,
		item.ImageURLs[0], item.ImageURLs[1], item.ImageURLs[2], item.ImageURLs[3], item.ImageURLs[4],
	).Scan(&id)
	if err != nil {
		log.Println("Error creating MerchItem: ", err)
		return "", err
	}

	// log.Println("MerchItem created: ", id)
	return id, nil
}

// GetByID retrieves a MerchItem by its ID
func GetByID(id string) (*MerchItem, error) {
	query := `
		SELECT ` + ALL_COLUMNS + `
		FROM "MerchItem"
		WHERE id = $1
	`
	row := database.Db.QueryRow(query, id)
	merch, err := scanMerchItem(row)
	if err != nil {
		return nil, err
	}

	return merch, nil
}

// Update modifies an existing MerchItem in the database
func (item *MerchItem) Update() error {
	query := `
		UPDATE "MerchItem"
		SET name = $1, description = $2, price = $3, inventory = $4, type = $5, height = $6, width = $7, unit = $8, image_url1 = $9, image_url2 = $10, image_url3 = $11, image_url4 = $12, image_url5 = $13
		WHERE id = $14
	`
	_, err := database.Db.Exec(
		query, item.Name, item.Description, item.Price, item.Inventory, item.Type, item.Height, item.Width, item.Unit,
		item.ImageURLs[0], item.ImageURLs[1], item.ImageURLs[2], item.ImageURLs[3], item.ImageURLs[4],
		item.ID,
	)
	if err != nil {
		log.Print("Error updating MerchItem: ", err)
		return err
	}
	return nil
}

// DeleteByID removes a MerchItem by ID
func DeleteByID(id string) error {
	query := `DELETE FROM "MerchItem" WHERE id = $1`
	_, err := database.Db.Exec(query, id)
	if err != nil {
		log.Print("Error deleting MerchItem: ", err)
		return err
	}
	return nil
}

// GetByOwnerID retrieves all MerchItems belonging to a specific owner by their owner ID.
func GetByOwnerID(ownerID string) ([]MerchItem, error) {
	query := `
		SELECT ` + ALL_COLUMNS + `
		FROM "MerchItem"
		WHERE owner_id = $1
		ORDER BY timestamp DESC
	`
	rows, err := database.Db.Query(query, ownerID)
	if err != nil {
		log.Printf("Error fetching MerchItems for owner %s: %v", ownerID, err)
		return nil, err
	}
	defer rows.Close()

	items, err := scanMerchItems(rows)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func SearchMerch(keyword *string, typeArg *string, minPrice *float64, maxPrice *float64, page *int, pageSize *int, sortBy *string, sortOrder *string) (*model.MerchSearchResult, error) {
	var items []MerchItem

	// Build the query dynamically
	query := `FROM "MerchItem" WHERE 1=1` // 1=1 to start the where chain
	args := []interface{}{}

	// Add default values to optional parameters
	if keyword != nil && *keyword != "" {
		query += fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d)`, len(args)+1, len(args)+2)
		args = append(args, "%"+*keyword+"%", "%"+*keyword+"%")
	}
	if typeArg != nil && *typeArg != "" {
		query += fmt.Sprintf(` AND type = $%d`, len(args)+1)
		args = append(args, *typeArg)
	}
	if minPrice != nil {
		query += fmt.Sprintf(` AND price >= $%d`, len(args)+1)
		args = append(args, *minPrice)
	}
	if maxPrice != nil {
		query += fmt.Sprintf(` AND price <= $%d`, len(args)+1)
		args = append(args, *maxPrice)
	}

	if page == nil || *page < 1 {
		// Set defaults for pagination
		p := 1
		page = &p
	}
	if pageSize == nil || *pageSize < 1 {
		ps := 10
		pageSize = &ps
	}

	// Count total items for pagination
	totalQuery := `SELECT COUNT(*) ` + query
	var totalItems int
	err := database.Db.QueryRow(totalQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, err
	}
	totalPages := (totalItems + *pageSize - 1) / *pageSize // ceil

	// Add sorting
	sortField := "name"
	if sortBy != nil && *sortBy != "" {
		sortField = *sortBy
	}
	order := "asc"
	if sortOrder != nil && *sortOrder == "desc" {
		order = "desc"
	}
	query = `SELECT ` + ALL_COLUMNS + ` ` + query + fmt.Sprintf(` ORDER BY %s %s`, sortField, order)

	// Add pagination
	offset := (*page - 1) * *pageSize
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, *pageSize, offset)

	// Execute the query
	rows, err := database.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err = scanMerchItems(rows)
	if err != nil {
		return nil, err
	}

	// convert to model objects
	var merchItems []*model.MerchItem
	for _, item := range items {
		merchItems = append(merchItems, item.ToGraphqlMerchItem())
	}

	return &model.MerchSearchResult{
		Items:       merchItems,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: *page,
		PageSize:    *pageSize,
	}, nil
}
