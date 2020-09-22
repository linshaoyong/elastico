package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func TestisIndexEarlyThan(t *testing.T) {
	now := time.Now().Unix()
	old, _ := isIndexEarlyThan("hello", []string{"20060102", "2006.01.02"}, now-7*86400)
	assert.False(t, old)

	old, _ = isIndexEarlyThan("hello-20200102", []string{"2006.01.02"}, now-7*86400)
	assert.False(t, old)

	old, _ = isIndexEarlyThan("hello-20200102", []string{"2006.01.02", "20060102"}, now-7*86400)
	assert.True(t, old)

	old, _ = isIndexEarlyThan("hello-20200102", []string{"2006.01.02"}, now-7*86400)
	assert.False(t, old)

	old, _ = isIndexEarlyThan("hello-2020.01.02", []string{"20060102", "2006.01.02"}, now-7*86400)
	assert.True(t, old)

	old, _ = isIndexEarlyThan("hello-2030.01", []string{"2006.01"}, now-7*86400)
	assert.False(t, old)

	old, _ = isIndexEarlyThan("hello-202001", []string{"200601"}, now-50*86400)
	assert.True(t, old)
}

func TestFilterIndexsEarlyThan(t *testing.T) {
	names := []string{
		"hello",
		"index-a-2020.05.06",
		"index-e-2020.06.29",
		"index-c-2030.07.09",
		"index-b-20200609",
		"index-b-20300606",
	}
	olds := filterIndexsEarlyThan(names, 7, map[string]int64{})
	assert.Equal(t, 3, len(olds))
	assert.True(t, contains(olds, "index-a-2020.05.06"))
	assert.True(t, contains(olds, "index-e-2020.06.29"))
	assert.True(t, contains(olds, "index-b-20200609"))
	assert.False(t, contains(olds, "index-c-2030.07.09"))
}

func TestFilterIndexsEarlyThanWithCustomDays(t *testing.T) {
	var customDays = map[string]int64{
		"index-e": 10000000,
		"index-c": 7,
	}

	names := []string{
		"hello",
		"index-a-2020.05.06",
		"index-e-2020.06.29",
		"index-c-2030.07.09",
		"index-b-20200609",
		"index-b-20300606",
	}
	olds := filterIndexsEarlyThan(names, 7, customDays)
	assert.Equal(t, 2, len(olds))
	assert.True(t, contains(olds, "index-a-2020.05.06"))
	assert.True(t, contains(olds, "index-b-20200609"))
	assert.False(t, contains(olds, "index-e-2020.06.29"))
	assert.False(t, contains(olds, "index-c-2030.07.09"))
}
