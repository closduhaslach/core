package driveq

type Field string

const (
	Name         Field = "name"
	MimeType     Field = "mimeType"
	FullText     Field = "fullText"
	Trashed      Field = "trashed"
	Starred      Field = "starred"
	ModifiedTime Field = "modifiedTime"
	ViewedByMe   Field = "viewedByMe"
	SharedWithMe Field = "sharedWithMeTime"
	Owners       Field = "owners"
	Writers      Field = "writers"
	Readers      Field = "readers"
	Parents      Field = "parents"
)
