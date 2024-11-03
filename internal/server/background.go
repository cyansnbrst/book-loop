package server

import (
	"fmt"

	"bookloop.net/pkg/sl"
)

func (s *Server) background(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				s.logger.Error("an error occured", sl.Err(fmt.Errorf("%s", err)))
			}
		}()

		fn()
	}()
}
