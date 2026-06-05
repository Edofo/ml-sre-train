package recommender

import "context"

type Recommender struct {
}

func (r *Recommender) Recommend(ctx context.Context, userID string) ([]string, error) {
	return []string{"item1", "item2", "item3"}, nil
}
