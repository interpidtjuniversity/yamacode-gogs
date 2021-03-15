package invoker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PipeLIneService(t *testing.T) {
	response := InvokePipeLineService()
	assert.Equal(t, response.Success, true)
}
