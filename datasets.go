package coze

type datasets struct {
	Documents *datasetsDocuments
}

func newDatasets(core *core) *datasets {
	return &datasets{
		Documents: newDocuments(core),
	}
}
