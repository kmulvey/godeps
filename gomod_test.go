package godeps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackupOriginalGoMod(t *testing.T) {
	t.Parallel()

	var err = backupOriginalGoMod()
	assert.NoError(t, err)
}
