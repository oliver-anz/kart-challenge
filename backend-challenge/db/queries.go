package db

import (
	"backend-challenge/models"
	"context"
	"database/sql"
	"fmt"
)

func (db *DB) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var thumbnail, mobile, tablet, desktop sql.NullString

		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &thumbnail, &mobile, &tablet, &desktop)
		if err != nil {
			return nil, err
		}

		// Only create Image object if at least one field has actual content
		if thumbnail.String != "" || mobile.String != "" || tablet.String != "" || desktop.String != "" {
			p.Image = &models.ProductImage{
				Thumbnail: thumbnail.String,
				Mobile:    mobile.String,
				Tablet:    tablet.String,
				Desktop:   desktop.String,
			}
		}

		products = append(products, p)
	}

	return products, rows.Err()
}

func (db *DB) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products WHERE id = ?`

	var p models.Product
	var thumbnail, mobile, tablet, desktop sql.NullString

	err := db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Category, &p.Price, &thumbnail, &mobile, &tablet, &desktop)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Only create Image object if at least one field has actual content
	if thumbnail.String != "" || mobile.String != "" || tablet.String != "" || desktop.String != "" {
		p.Image = &models.ProductImage{
			Thumbnail: thumbnail.String,
			Mobile:    mobile.String,
			Tablet:    tablet.String,
			Desktop:   desktop.String,
		}
	}

	return &p, nil
}

func (db *DB) IsCouponValid(ctx context.Context, code string) (bool, error) {
	// Empty coupon code is valid (no coupon provided)
	if code == "" {
		return true, nil
	}

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
