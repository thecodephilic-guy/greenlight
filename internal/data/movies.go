package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"` //The version number starts with 1 and be in incremented each time the information is updated
}

/*
Note:
It’s crucial to point out here that all the fields in our Movie struct are exported
(i.e. start with a capital letter), which is necessary for them to be visible to Go’s
encoding/json package. Any fields which aren’t exported won’t be included when
encoding a struct to JSON.
*/
