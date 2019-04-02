package bedtimes

import "context"

// Save persists data in FireBase
func Save(report *Report) error {
	ctx := context.Background()

	_, err := client.Collection(collection).Doc(doc).Set(ctx, report)

	return err
}
