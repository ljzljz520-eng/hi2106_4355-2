package grades

type SubmissionResult struct {
	Student Student
	Average string
}

type SubmitCallback func(SubmissionResult)

type Extensions struct {
	AfterSubmit SubmitCallback
}

type ExtensionExecutor struct {
	extensions Extensions
}

func NewExtensionExecutor(extensions Extensions) ExtensionExecutor {
	return ExtensionExecutor{extensions: extensions}
}

func (e ExtensionExecutor) ExecuteAfterSubmit(result SubmissionResult) {
	e.extensions.AfterSubmit(result)
}
