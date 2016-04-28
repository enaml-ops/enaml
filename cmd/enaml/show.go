package main

import "io"

type show struct {
	release string
}

func (s *show) All(w io.Writer) error {
	return nil
}
