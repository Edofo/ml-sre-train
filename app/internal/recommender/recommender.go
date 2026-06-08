package recommender

import "context"

type Stub struct {
}

func (s *Stub) Recommend(ctx context.Context, userID string) ([]string, error) {
	return []string{"item1", "item2", "item3"}, nil
}
