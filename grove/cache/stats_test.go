package cache

import (
	"testing"
)

func StatsInc(s Stats, num int) {
	for i := 0; i < num; i++ {
		s.IncConnections()
	}
}

func StatsDec(s Stats, num int) {
	for i := 0; i < num; i++ {
		s.DecConnections()
	}
}

func TestStatsCount(t *testing.T) {
	{
		r := RemapRule{RemapRuleBase: RemapRuleBase{Name: "foo"}}
		stats := NewStats([]RemapRule{r})
		expected := 10
		StatsInc(stats, expected)
		if actual := stats.Connections(); actual != uint64(expected) {
			t.Errorf("Stats.Connections() expected %v actual %v", expected, actual)
		}
	}

	{
		r := RemapRule{RemapRuleBase: RemapRuleBase{Name: "foo"}}
		stats := NewStats([]RemapRule{r})
		count := 10
		StatsInc(stats, count)
		StatsDec(stats, count)
		if actual := stats.Connections(); actual != 0 {
			t.Errorf("Stats.Connections() expected %v actual %v", 0, actual)
		}
	}

	{
		r := RemapRule{RemapRuleBase: RemapRuleBase{Name: "foo"}}
		stats := NewStats([]RemapRule{r})
		count := 10
		StatsInc(stats, count)
		StatsDec(stats, 1)
		if actual := stats.Connections(); actual != uint64(count-1) {
			t.Errorf("Stats.Connections() expected %v actual %v", count-1, actual)
		}
	}

}
