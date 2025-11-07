package db

import (
	"backend-challenge/models"
	"context"
	"database/sql"
	"fmt"
)

func (db *DB) GetAllProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products`

	// Add pagination if limit is specified
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, rows.Err()
}

func (db *DB) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products WHERE id = ?`

	row := db.QueryRowContext(ctx, query, id)
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (db *DB) IsCouponValid(ctx context.Context, code string) (bool, error) {
	// Coupons are preprocessed but best to defensively check length
	if len(code) < 8 || len(code) > 10 {
		return false, nil
	}

	query := `SELECT COUNT(*) FROM valid_coupons WHERE code = ?`
	var count int
	err := db.QueryRowContext(ctx, query, code).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to validate coupon: %w", err)
	}

	return count > 0, nil
}

// scanProduct scans a row into a Product, handling nullable image fields
func scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Product, error) {
	var p models.Product
	var thumbnail, mobile, tablet, desktop sql.NullString

	err := scanner.Scan(&p.ID, &p.Name, &p.Category, &p.Price,
		&thumbnail, &mobile, &tablet, &desktop)
	if err != nil {
		return nil, err
	}

	// Only create Image object if at least one field has content
	if thumbnail.String != "" || mobile.String != "" ||
		tablet.String != "" || desktop.String != "" {
		p.Image = &models.ProductImage{
			Thumbnail: thumbnail.String,
			Mobile:    mobile.String,
			Tablet:    tablet.String,
			Desktop:   desktop.String,
		}
	}

	return &p, nil
}
