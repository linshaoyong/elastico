package cmd

import "testing"

func TestGetOldIndexs(t *testing.T) {
	names := []string{
		"k8s_cluster_xm-2020.05.06",
		"fusion_cdn-2020.06.29",
		"hexv2-2030.07.09",
	}

	olds := getOldIndexNames(names)
	if len(olds) != 1 {
		t.Errorf("fail...")
	}
}
