package cmd

import (
	"strings"
	"testing"
)

func Test_formatterOf(t *testing.T) {
	got := formatterOf([]string{"/s"})
	for i := 0; i<1000; i++ {
		if strings.ToUpper(got.Format("asd")) != "ASD" {
			t.Errorf("formatter acted unexpectedly") //bad place for this test!
		}
	}
}
