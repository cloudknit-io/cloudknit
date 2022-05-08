package tftmpl

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"

	"github.com/pkg/errors"
	"go.uber.org/zap/buffer"
)

const (
	TmplTFBackend  TemplateName = "terraform_backend"
	TmplTFData     TemplateName = "terraform_data"
	TmplTFModule   TemplateName = "terraform_module"
	TmplTFOutputs  TemplateName = "terraform_outputs"
	TmplTFProvider TemplateName = "terraform_provider"
	TmplTFSecrets  TemplateName = "terraform_secrets"
	TmplTFVersions TemplateName = "terraform_versions"
)

//go:embed terraform_backend.tmpl
var tmplTFBackend string

//go:embed terraform_data.tmpl
var tmplTFData string

//go:embed terraform_module.tmpl
var tmplTFModule string

//go:embed terraform_outputs.tmpl
var tmplTFOutputs string

//go:embed terraform_provider.tmpl
var tmplTFProvider string

//go:embed terraform_secrets.tmpl
var tmplTFSecrets string

//go:embed terraform_versions.tmpl
var tmplTFVersions string

type TerraformTemplates struct {
	templates map[TemplateName]*template.Template
}

var _ Template = (*TerraformTemplates)(nil)

func NewTerraformTemplates() (*TerraformTemplates, error) {
	tpl := TerraformTemplates{}
	if err := tpl.init(); err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (tpl *TerraformTemplates) init() error {
	tpl.templates = make(map[TemplateName]*template.Template, 7)

	funcMap := sprig.TxtFuncMap()
	funcMap["indentMultiline"] = indentMultiline

	t, err := template.New(string(TmplTFBackend)).Parse(tmplTFBackend)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFBackend)
	}
	tpl.templates[TmplTFBackend] = t

	t, err = template.New(string(TmplTFData)).Parse(tmplTFData)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFData)
	}
	tpl.templates[TmplTFData] = t

	t, err = template.New(string(TmplTFModule)).Funcs(funcMap).Parse(tmplTFModule)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFModule)
	}
	tpl.templates[TmplTFModule] = t

	t, err = template.New(string(TmplTFOutputs)).Parse(tmplTFOutputs)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFOutputs)
	}
	tpl.templates[TmplTFOutputs] = t

	t, err = template.New(string(TmplTFProvider)).Parse(tmplTFProvider)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFProvider)
	}
	tpl.templates[TmplTFProvider] = t

	t, err = template.New(string(TmplTFSecrets)).Parse(tmplTFSecrets)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFSecrets)
	}
	tpl.templates[TmplTFSecrets] = t

	t, err = template.New(string(TmplTFVersions)).Parse(tmplTFVersions)
	if err != nil {
		return errors.Errorf("error parsing template: %s", TmplTFVersions)
	}
	tpl.templates[TmplTFVersions] = t

	return nil
}

func indentMultiline(spaces int, v string) (string, error) {
	lines := strings.Split(v, "\n")
	for i := range lines {
		pad := strings.Repeat(" ", spaces)
		lines[i] = pad + strings.ReplaceAll(lines[i], "\n", "\n"+pad)
	}
	return strings.Join(lines, "\n"), nil
}

func (tpl *TerraformTemplates) Execute(vars any, tmplt TemplateName) (output string, err error) {
	tfTemplate := tpl.templates[tmplt]
	if tfTemplate == nil {
		return "", errors.Errorf("invalid template: %s", tmplt)
	}

	var b buffer.Buffer
	if err := tfTemplate.Execute(&b, vars); err != nil {
		return "", errors.Wrapf(err, "error executing template %s", tmplt)
	}

	return b.String(), nil
}

func (tpl *TerraformTemplates) Templates() []TemplateName {
	templates := make([]TemplateName, 0, len(tpl.templates))
	for k := range tpl.templates {
		templates = append(templates, k)
	}
	return templates
}
