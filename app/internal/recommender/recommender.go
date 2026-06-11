package recommender

import (
	"context"
	"errors"
)

type Stub struct {
}

func (s *Stub) Recommend(ctx context.Context, userID string) ([]string, error) {
	if userID == "boom" {
		return nil, errors.New("boom")
	}
	return []string{"item1", "item2", "item3"}, nil
}
