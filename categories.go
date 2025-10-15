package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Category struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Subcategories []*Subcategory `json:"subcategories"`
}

type Subcategory struct {
	ID         string                           `json:"id"`
	Name       string                           `json:"name"`
	CreatedAt  string                           `json:"created_at"`
	Version    int                              `json:"version"`
	Attributes map[string]*SubcategoryAttribute `json:"attributes"`
}

type SubcategoryAttribute struct {
	Required    bool        `json:"required"`
	Type        string      `json:"type"`
	Reference   string      `json:"reference"`
	IsMultiple  bool        `json:"is_multiple"`
	StringArray []string    `json:"string_array"`
	Validation  *Validation `json:"validation"`
}

type Validation struct {
	Numeric *struct {
		Rules []struct {
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		} `json:"rules"`
	} `json:"numeric"`
	Str *struct {
		Rules []struct {
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		} `json:"rules"`
	} `json:"str"`
}

type CategoryEdit struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	Subcategories []*Subcategory `json:"subcategories"`
	Version       int            `json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
}

func InsertCategories(ctx context.Context, tx pgx.Tx) error {
	var array []*Category

	data, err := os.ReadFile("data/categories_v1.json")
	if err != nil {
		return fmt.Errorf("failed to load categories %v ", err)
	}

	if err = json.Unmarshal(data, &array); err != nil {
		return fmt.Errorf("failed to load categories %v ", err)
	}

	// stmt := "INSERT INTO categories(id, name, subcategories, edits, version, created_at) VALUES($1, $2, $3, $4, $5, $6)"
	payload := make([][]any, 0, len(array))

	for _, c := range array {

		subData, err := json.Marshal(c.Subcategories)
		if err != nil {
			return err
		}

		createdAt := time.Now().UTC()
		edits := []*CategoryEdit{{
			ID:            c.ID,
			Name:          c.Name,
			Version:       1,
			Subcategories: c.Subcategories,
			CreatedAt:     createdAt,
		}}

		editsData, err := json.Marshal(edits)
		if err != nil {
			return err
		}

		args := []any{c.ID, c.Name, c.Image, subData, editsData, 1, createdAt}
		payload = append(payload, args)
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"categories"}, []string{"id", "name", "image", "subcategories", "edits", "version", "created_at"}, pgx.CopyFromRows(payload))
	if err != nil {
		return err
	}

	return nil
}
