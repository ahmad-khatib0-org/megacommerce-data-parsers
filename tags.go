package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func getTagsFromJsonFile(file string) ([]string, error) {
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

func removeDuplicateTags() error {
	tags, err := getTagsFromJsonFile("tags_v1.json")
	if err != nil {
		return err
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(tags))
	for _, t := range tags {
		if _, ok := seen[t]; !ok {
			seen[t] = struct{}{}
			result = append(result, t)
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	if err := os.WriteFile("data/tags_v1_unique.json", data, os.ModeAppend|os.ModePerm); err != nil {
		return err
	}

	return nil
}

func InsertTags(ctx context.Context, tx pgx.Tx) error {
	tags, err := getTagsFromJsonFile("tags_v1_unique.json")
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
