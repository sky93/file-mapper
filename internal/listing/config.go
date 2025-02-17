package listing

// Config stores all user options from CLI flags
type Config struct {
	RootPath       string
	Include        string
	Exclude        string
	GitTrackedOnly bool

	// Output style
	ShowTree        bool // tree or flat
	ShowContent     bool // whether to include file content at all
	SeparateContent bool // if true, print the tree/flat list first, then content after

	Output string // optional file path for output

	// Content details
	ShowLineNumbers   bool
	ShowHeaderFooters bool
}
