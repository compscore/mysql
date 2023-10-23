package check_template

import (
	"context"
)

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	return true, ""
}
