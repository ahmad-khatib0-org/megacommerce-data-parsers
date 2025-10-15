package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func getTagsFromJSONFile(file string) ([]string, error) {
	data, err := os.ReadFile(fmt.Sprintf("data/%s", file))
	if err != nil {
		return nil, err
	}

	var tags []string
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func InsertTags(ctx context.Context, tx pgx.Tx) error {
	tags, err := getTagsFromJSONFile("tags_v1.json")
	if err != nil {
		return err
	}

	payload := make([][]any, 0, len(tags))
	for _, t := range tags {
		payload = append(payload, []any{t})
	}

	if _, err := tx.CopyFrom(ctx, pgx.Identifier{"tags"}, []string{"name"}, pgx.CopyFromRows(payload)); err != nil {
		return err
	}

	return nil
}
