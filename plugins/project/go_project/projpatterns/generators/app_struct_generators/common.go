package app_struct_generators

type InitDepFuncGenArgs struct {
	InitFunctionName string
	Imports          map[string]string
	Functions        []InitFuncCall
}

type InitFuncCall struct {
	FuncName string
	Args     string

	ResultName string
	ResultType string

	ErrorMessage string
}
