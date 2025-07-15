package rbac

import (
	"path/filepath"
	"runtime"

	"github.com/casbin/casbin/v2"
)

// NewEnforcer creates a new Casbin enforcer
func NewEnforcer() (*casbin.Enforcer, error) {
	// Get the current file's directory
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Load the model and policy files
	modelPath := filepath.Join(dir, "model.conf")
	policyPath := filepath.Join(dir, "policy.csv")

	// Create the enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
