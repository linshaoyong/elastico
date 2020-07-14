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

func TestIsOldIndex(t *testing.T) {
	now := time.Now().Unix()
	old := isOldIndex("hello", []string{"20060102", "2006.01.02"}, 7, now)
	assert.False(t, old)

	old = isOldIndex("hello-20200102", []string{"2006.01.02"}, 7, now)
	assert.False(t, old)

	old = isOldIndex("hello-20200102", []string{"2006.01.02", "20060102"}, 7, now)
	assert.True(t, old)

	old = isOldIndex("hello-20200102", []string{"2006.01.02"}, 7, now)
	assert.False(t, old)

	old = isOldIndex("hello-2020.01.02", []string{"20060102", "2006.01.02"}, 7, now)
	assert.True(t, old)
}

func TestGetOldIndexs(t *testing.T) {
	names := []string{
		"hello",
		"k8s_cluster_xm-2020.05.06",
		"fusion_cdn-2020.06.29",
		"hexv2-2030.07.09",
		"fruit_dc_coverregion_compare-20200609",
		"fruit_dc_coverregion_compare-20300606",
	}
	olds := getOldIndexNames(names, 7)
	assert.Equal(t, 3, len(olds))
	assert.True(t, contains(olds, "k8s_cluster_xm-2020.05.06"))
	assert.True(t, contains(olds, "fusion_cdn-2020.06.29"))
	assert.True(t, contains(olds, "fruit_dc_coverregion_compare-20200609"))
	assert.False(t, contains(olds, "hexv2-2030.07.09"))
}
