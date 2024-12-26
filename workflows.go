package coze

type workflows struct {
	Runs *workflowRuns
}

func newWorkflows(core *core) *workflows {
	return &workflows{
		Runs: newWorkflowRun(core),
	}
}
