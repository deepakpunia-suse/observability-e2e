package terraform

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	ginkgo "github.com/onsi/ginkgo/v2"
)

type TerraformContext struct {
	Options *terraform.Options
}

// TerraformOptions is your input struct to control dir, vars, env vars
type TerraformOptions struct {
	TerraformDir string
	Vars         map[string]interface{}
	EnvVars      map[string]string
}

// NewTerraformContext initializes a new Terraform context with options.
func NewTerraformContext(opts TerraformOptions) (*TerraformContext, error) {
	tfOpts := &terraform.Options{
		TerraformDir: opts.TerraformDir,
	}

	// Only assign Vars if it is not nil
	if opts.Vars != nil {
		tfOpts.Vars = opts.Vars
	} else {
		tfOpts.Vars = map[string]interface{}{} // or leave it nil if terraform supports that
	}

	// Only assign EnvVars if it is not nil
	if opts.EnvVars != nil {
		tfOpts.EnvVars = opts.EnvVars
	} else {
		tfOpts.EnvVars = map[string]string{} // or leave nil if supported
	}

	return &TerraformContext{
		Options: tfOpts,
	}, nil
}

func (ctx *TerraformContext) InitAndApply() (string, error) {
	return terraform.InitAndApplyE(ginkgo.GinkgoT(), ctx.Options)
}

func (ctx *TerraformContext) OutputAll() map[string]interface{} {
	return terraform.OutputAll(ginkgo.GinkgoT(), ctx.Options)
}

// Destroy tears down the infrastructure.
func (ctx *TerraformContext) Destroy() (string, error) {
	return terraform.DestroyE(ginkgo.GinkgoT(), ctx.Options)
}

// DestroyTarget destroys specific Terraform resources using the -target flag,
// and resets the Targets field so the context can be reused safely.
func (ctx *TerraformContext) DestroyTarget(targets ...string) (string, error) {
	// Backup original targets
	originalTargets := ctx.Options.Targets

	// Set the target(s) temporarily
	ctx.Options.Targets = targets

	// Perform targeted destroy
	output, err := terraform.DestroyE(ginkgo.GinkgoT(), ctx.Options)

	// Reset targets to their original state
	ctx.Options.Targets = originalTargets

	return output, err
}
