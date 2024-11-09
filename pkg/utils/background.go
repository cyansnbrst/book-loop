package utils

import (
	"fmt"
	"log/slog"
	"sync"

	"bookloop.net/pkg/sl"
)

func Background(wg *sync.WaitGroup, logger *slog.Logger, fn func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				logger.Error("an error occured", sl.Err(fmt.Errorf("%s", err)))
			}
		}()

		fn()
	}()
}
