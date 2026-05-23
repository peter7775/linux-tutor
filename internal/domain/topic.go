package domain

type Topic struct {
	ID          string
	Track       string
	Code        string
	Title       string
	Description string
	Tags        []string
	Difficulty  int
	Weight      int
	ParentID    string
	Active      bool
}
