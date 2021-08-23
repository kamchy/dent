package data

import (
	"testing"
	"time"
)

func TestDate_test(t *testing.T) {

	d, _ := time.ParseInLocation(VISIT_DATE_LAYOUT, "2021-08-30 12:00", time.Local)
	v := Visit{3, d, 5, 34, "asdasdasd", "Ala", "Chyla"}

	t.Run("Date", func(t *testing.T) {
		exp, got := "2021-08-30", v.VisitDateString()
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})

	t.Run("Date", func(t *testing.T) {
		exp, got := "12:00", v.VisitTimeString()
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	})
}
