package errs

import (
	"github.com/fwhezfwhez/errorx"
	"testing"
	"zonst/qipai/api/configapisrv/config"
)

func TestSave(t *testing.T) {
	config.Node.Mode = "pro"
	config.RegisterNode("configapisrv")

	SaveError(errorx.NewServiceError("测试上报", 1))
}
