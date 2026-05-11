package cryptoutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvUtils(t *testing.T) {
	t.Run("ParseEnv", func(t *testing.T) {
		raw := []byte("KEY1=value1\nKEY2=value2\n")
		env, err := ParseEnv(raw)
		require.NoError(t, err)
		assert.Equal(t, "value1", env["KEY1"])
		assert.Equal(t, "value2", env["KEY2"])

		_, err = ParseEnv([]byte("KEY1=value1\nMALFORMED_ENTRY\n"))
		require.Error(t, err)
	})

	t.Run("EncodeAndNormalizeEnv", func(t *testing.T) {
		env := map[string]string{"KEY1": "VAL1", "KEY2": "VAL2"}
		encoded, err := EncodeEnv(env)
		require.NoError(t, err)
		assert.Contains(t, string(encoded), "KEY2=\"VAL2\"")

		normalized := NormalizeEnv(env)
		assert.Equal(t, "KEY1=VAL1\nKEY2=VAL2\n", string(normalized))
	})

	t.Run("Compression", func(t *testing.T) {
		data := []byte("SOME_DATA_TO_COMPRESS")
		compressed, _ := CompressEnv(data)
		decompressed, err := DecompressEnv(compressed)
		require.NoError(t, err)
		assert.Equal(t, data, decompressed)
	})
}

func TestDiffEnvVersions(t *testing.T) {
	oldVersion := map[string]string{"REMOVED": "1", "MODIFIED": "old", "SAME": "val"}
	newVersion := map[string]string{"ADDED": "1", "MODIFIED": "new", "SAME": "val"}

	result := DiffEnvVersions(oldVersion, newVersion)
	assert.ElementsMatch(t, []string{"ADDED"}, result.Added)
	assert.ElementsMatch(t, []string{"REMOVED"}, result.Removed)
	assert.ElementsMatch(t, []string{"MODIFIED"}, result.Modified)

	resultEmpty := DiffEnvVersions(nil, nil)
	assert.Empty(t, resultEmpty.Added)

	resultRemoveAll := DiffEnvVersions(oldVersion, nil)
	assert.ElementsMatch(t, []string{"REMOVED", "MODIFIED", "SAME"}, resultRemoveAll.Removed)
}
