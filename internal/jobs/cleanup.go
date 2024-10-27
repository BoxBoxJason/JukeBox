package jobs

import (
	"time"

	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	"github.com/boxboxjason/jukebox/pkg/logger"
)

// TokenCleanup starts a background job that deletes expired tokens every hour
func TokenCleanup() {
	token_cleanup_ticker := time.NewTicker(1 * time.Hour)
	defer token_cleanup_ticker.Stop()

	err := db_controller.DeleteExpiredTokens(nil)
	if err != nil {
		logger.Error("Error deleting expired tokens", err)
	} else {
		logger.Info("Deleted expired tokens")
	}

	go func() {
		for {
			select {
			case <-token_cleanup_ticker.C:
				if err := db_controller.DeleteExpiredTokens(nil); err != nil {
					logger.Error("Error deleting expired tokens", err)
				} else {
					logger.Info("Deleted expired tokens")
				}
			}
		}
	}()
}
