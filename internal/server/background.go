package server

import (
	"fmt"

	"bookloop.net/pkg/sl"
)

func (s *Server) background(fn func()) {
	s.WG.Add(1)

	go func() {
		defer s.WG.Done()

		defer func() {
			if err := recover(); err != nil {
				s.Logger.Error("an error occured", sl.Err(fmt.Errorf("%s", err)))
			}
		}()

		fn()
	}()
}
