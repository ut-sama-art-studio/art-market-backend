package merchitems

import (
	"database/sql"
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
		Images:      merch.ImageURLs,
		Unit:        merch.Unit,
	}
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
		SELECT id, owner_id, name, description, price, inventory, type, height, width, unit, image_url1, image_url2, image_url3, image_url4, image_url5, timestamp
		FROM "MerchItem"
		WHERE id = $1
	`
	var item MerchItem
	err := database.Db.QueryRow(query, id).Scan(
		&item.ID, &item.OwnerID, &item.Name, &item.Description, &item.Price, &item.Inventory,
		&item.Type, &item.Height, &item.Width, &item.Unit, &item.ImageURLs[0], &item.ImageURLs[1], &item.ImageURLs[2], &item.ImageURLs[3], &item.ImageURLs[4],
		&item.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Print("Error fetching MerchItem by ID: ", err)
		return nil, err
	}
	return &item, nil
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
		SELECT id, owner_id, name, description, price, inventory, type, height, width, unit, image_url1, image_url2, image_url3, image_url4, image_url5, timestamp
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

	var items []MerchItem
	for rows.Next() {
		var item MerchItem
		item.ImageURLs = make([]*string, 5)
		err := rows.Scan(
			&item.ID, &item.OwnerID, &item.Name, &item.Description, &item.Price, &item.Inventory,
			&item.Type, &item.Height, &item.Width, &item.Unit, &item.ImageURLs[0], &item.ImageURLs[1], &item.ImageURLs[2], &item.ImageURLs[3], &item.ImageURLs[4],
			&item.Timestamp,
		)
		if err != nil {
			log.Printf("Error scanning MerchItem for owner %s: %v", ownerID, err)
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error with MerchItems rows for owner %s: %v", ownerID, err)
		return nil, err
	}

	return items, nil
}
