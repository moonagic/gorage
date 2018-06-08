package data

// UploadItem Upload item template
type UploadItem struct {
	UUID       string
	FileName   string
	Directory  string
	TagTime    string
	UploadTime string
}

// KeyMap keycache map
type KeyMap struct {
	Index   int
	UUID    string
	TagTime string
}
