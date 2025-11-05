package db

import (
	"backend-challenge/models"
	"database/sql"
	"fmt"
)

func (db *DB) InsertProduct(p *models.Product) error {
	query := `
		INSERT OR REPLACE INTO products (id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	var thumbnail, mobile, tablet, desktop sql.NullString
	if p.Image != nil {
		thumbnail = sql.NullString{String: p.Image.Thumbnail, Valid: p.Image.Thumbnail != ""}
		mobile = sql.NullString{String: p.Image.Mobile, Valid: p.Image.Mobile != ""}
		tablet = sql.NullString{String: p.Image.Tablet, Valid: p.Image.Tablet != ""}
		desktop = sql.NullString{String: p.Image.Desktop, Valid: p.Image.Desktop != ""}
	}

	_, err := db.Exec(query, p.ID, p.Name, p.Category, p.Price, thumbnail, mobile, tablet, desktop)
	return err
}

func (db *DB) GetAllProducts() ([]models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products`

	rows, err := db.Query(query)
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

		if thumbnail.Valid || mobile.Valid || tablet.Valid || desktop.Valid {
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

func (db *DB) GetProductByID(id string) (*models.Product, error) {
	query := `SELECT id, name, category, price, image_thumbnail, image_mobile, image_tablet, image_desktop FROM products WHERE id = ?`

	var p models.Product
	var thumbnail, mobile, tablet, desktop sql.NullString

	err := db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Category, &p.Price, &thumbnail, &mobile, &tablet, &desktop)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if thumbnail.Valid || mobile.Valid || tablet.Valid || desktop.Valid {
		p.Image = &models.ProductImage{
			Thumbnail: thumbnail.String,
			Mobile:    mobile.String,
			Tablet:    tablet.String,
			Desktop:   desktop.String,
		}
	}

	return &p, nil
}

func (db *DB) InsertCoupon(code string, count int) error {
	query := `INSERT OR REPLACE INTO valid_coupons (code, occurrence_count) VALUES (?, ?)`
	_, err := db.Exec(query, code, count)
	return err
}

func (db *DB) IsCouponValid(code string) (bool, error) {
	if code == "" {
		return true, nil
	}

	if len(code) < 8 || len(code) > 10 {
		return false, nil
	}

	query := `SELECT COUNT(*) FROM valid_coupons WHERE code = ?`
	var count int
	err := db.QueryRow(query, code).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to validate coupon: %w", err)
	}

	return count > 0, nil
}
