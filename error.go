package frontapp

import "fmt"

// Error is an error returned by api
type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Title   string `json:"title"`
}

// Error implements error interface
func (e Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Status, e.Title, e.Message)
}
