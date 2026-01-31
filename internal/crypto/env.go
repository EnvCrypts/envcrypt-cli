package cryptoutils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

func ParseEnv(raw []byte) (map[string]string, error) {
	env, err := godotenv.UnmarshalBytes(raw)
	if err != nil {
		return nil, err
	}

	return env, nil
}
func EncodeEnv(env map[string]string) ([]byte, error) {
	s, err := godotenv.Marshal(env)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func NormalizeEnv(env map[string]string) []byte {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(env[k])
		b.WriteString("\n")
	}

	return []byte(b.String())
}

func CompressEnv(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(data); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecompressEnv(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	return io.ReadAll(gr)
}

func PrepareEnvForStorage(parsed map[string]string) ([]byte, error) {
	normalized := NormalizeEnv(parsed)

	compressed, err := CompressEnv(normalized)
	if err != nil {
		return nil, errors.New("could not compress env")
	}

	return compressed, nil
}

func PrepareEnvForRollback(env map[string]string) ([]byte, error) {
	normalized := NormalizeEnv(env)

	compressed, err := CompressEnv(normalized)
	if err != nil {
		return nil, err
	}

	return compressed, nil
}

func ReadCompressedEnv(data []byte) (map[string]string, error) {
	decompressed, err := DecompressEnv(data)
	if err != nil {
		return nil, errors.New("could not decompress env")
	}

	return ParseEnv(decompressed)
}

type DiffingResult struct {
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

func DiffEnvVersions(oldVersion, newVersion map[string]string) DiffingResult {

	var Added, Removed, Modified []string

	for key, val := range newVersion {
		if _, exists := oldVersion[key]; !exists {
			Added = append(Added, key)
		} else {
			if val != oldVersion[key] {
				Modified = append(Modified, key)
			}
		}
	}
	for key := range oldVersion {
		if _, exists := newVersion[key]; !exists {
			Removed = append(Removed, key)
		}
	}

	return DiffingResult{Added: Added, Removed: Removed, Modified: Modified}
}
