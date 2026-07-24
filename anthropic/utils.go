package anthropic

import (
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// getDateRangeFromQuals derives an inclusive [start, end] date range (UTC,
// truncated to day) from timestamp quals on the given column. If no quals are
// present, it returns the trailing defaultDays window ending yesterday (most
// analytics data lags by up to a day).
func getDateRangeFromQuals(d *plugin.QueryData, column string, defaultDays int) (time.Time, time.Time) {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	start := now.AddDate(0, 0, -defaultDays)
	end := now

	if quals, ok := d.Quals[column]; ok {
		for _, q := range quals.Quals {
			ts := q.Value.GetTimestampValue()
			if ts == nil {
				continue
			}
			t := ts.AsTime().UTC().Truncate(24 * time.Hour)
			switch q.Operator {
			case "=":
				start, end = t, t
			case ">":
				start = t.AddDate(0, 0, 1)
			case ">=":
				start = t
			case "<":
				end = t.AddDate(0, 0, -1)
			case "<=":
				end = t
			}
		}
	}
	if end.Before(start) {
		end = start
	}
	return start, end
}
