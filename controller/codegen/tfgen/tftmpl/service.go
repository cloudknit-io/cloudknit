package tftmpl

type TemplateName string

type Template interface {
	Execute(vars any, tmplt TemplateName) (output string, err error)
	Templates() []TemplateName
}
