package rbac

import (
	"path/filepath"
	"runtime"

	"github.com/bgaurav7/gin-microservice-boilerplate/config"
	"github.com/casbin/casbin/v2"
)

// NewEnforcer creates a new Casbin enforcer
func NewEnforcer() (*casbin.Enforcer, error) {
	// Load config
	cfg, err := config.Load()
	if err == nil && cfg.RBAC.ModelPath != "" && cfg.RBAC.PolicyPath != "" {
		// Use paths from config
		enforcer, err := casbin.NewEnforcer(cfg.RBAC.ModelPath, cfg.RBAC.PolicyPath)
		if err == nil {
			return enforcer, nil
		}
	}

	// Fallback to default paths if config loading failed or paths are not set
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	modelPath := filepath.Join(dir, "model.conf")
	policyPath := filepath.Join(dir, "policy.csv")

	// Create the enforcer with default paths
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	return enforcer, nil
}
