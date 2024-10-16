// Package historyclean is a watcher implement used to delete expired record from the database.
package historyclean

import (
	"context"
	"github.com/rosas99/streaming/internal/nightwatch/watcher"
	"github.com/rosas99/streaming/internal/pkg/client/store"
	"github.com/rosas99/streaming/pkg/log"
	"time"
)

var _ watcher.Watcher = (*historiesCleanWatcher)(nil)

// watcher implement.
type historiesCleanWatcher struct {
	store store.Interface
}

func (w *historiesCleanWatcher) Spec() string {
	return watcher.EveryDay
}

// Run runs the watcher.
func (w *historiesCleanWatcher) Run() {
	_, histories, err := w.store.Sms().Histories().List(context.Background())
	if err != nil {
		log.Errorw(err, "Failed to list secrets")
		return
	}

	for _, history := range histories {
		// Check if the history is older than one year
		// deletes all records that are older than one year.
		if history.CreatedAt.Unix() < time.Now().AddDate(-1, 0, 0).Unix() {
			filter := map[string]any{"id": history.ID}
			err := w.store.Sms().Histories().Delete(context.TODO(), filter)
			if err != nil {
				log.Warnw("Failed to delete secret from database", "userID", history.ID, "name", "secret.Name")
				continue
			}
			log.Infow("Successfully deleted secret from database", "userID", history.ID, "name", "secret.Name")
		}
	}
}

// SetAggregateConfig initializes the watcher for later execution.
func (w *historiesCleanWatcher) SetAggregateConfig(ctx context.Context, config *watcher.Config) error {
	w.store = config.Store
	return nil
}

func init() {
	watcher.Register("historyclean", &historiesCleanWatcher{})
}
