package terraform

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Wrapper struct {
	ctx context.Context
	tf  *tfexec.Terraform
}

type StateWrapper struct {
	state *tfjson.State
}
