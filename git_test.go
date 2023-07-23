package godeps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePrBranch(t *testing.T) {
	t.Parallel()

	var err = createPrBranch("branchName")
	assert.NoError(t, err)
}
