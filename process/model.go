package process

import "time"

type DataEntry struct {
	ID    int
	Title string
	Tags  []string
	Time  time.Time
}

type Article struct {
	ID      int
	Title   string
	Content string
	Time    time.Time
}
