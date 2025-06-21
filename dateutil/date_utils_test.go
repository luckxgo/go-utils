package dateutil

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	now := Now()
	if now.After(time.Now()) {
		t.Error("Now() returned a time in the future")
	}
}

func TestFormatDateTime(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: "2023-10-05 15:30:45",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDateTime(tt.t); got != tt.want {
				t.Errorf("FormatDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: "2023-10-05",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.t); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: "15:30:45",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTime(tt.t); got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDateTime(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    time.Time
		wantErr bool
	}{{
		name:    "valid date time",
		s:       "2023-10-05 15:30:45",
		want:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		wantErr: false,
	}, {
		name:    "empty string",
		s:       "",
		want:    time.Time{},
		wantErr: true,
	}, {
		name:    "invalid format",
		s:       "2023/10/05 15:30:45",
		want:    time.Time{},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDateTime(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    time.Time
		wantErr bool
	}{{
		name:    "valid date",
		s:       "2023-10-05",
		want:    time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		wantErr: false,
	}, {
		name:    "empty string",
		s:       "",
		want:    time.Time{},
		wantErr: true,
	}, {
		name:    "invalid format",
		s:       "2023/10/05",
		want:    time.Time{},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYear(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 2023,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Year(tt.t); got != tt.want {
				t.Errorf("Year() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonth(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 10,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Month(tt.t); got != tt.want {
				t.Errorf("Month() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDay(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 5,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Day(tt.t); got != tt.want {
				t.Errorf("Day() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeGenerate(t *testing.T) {
	start := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		start     time.Time
		end       time.Time
		unit      TimeUnit
		wantLen   int
		wantFirst time.Time
		wantLast  time.Time
	}{{
		name:      "daily range",
		start:     start,
		end:       end,
		unit:      DayUnit,
		wantLen:   3,
		wantFirst: start,
		wantLast:  end,
	}, {
		name:      "hourly range",
		start:     time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
		end:       time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
		unit:      HourUnit,
		wantLen:   3,
		wantFirst: time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC),
		wantLast:  time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
	}, {
		name:      "monthly range",
		start:     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		end:       time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
		unit:      MonthUnit,
		wantLen:   3,
		wantFirst: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
		wantLast:  time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
	}, {
		name:      "yearly range",
		start:     time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC),
		end:       time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
		unit:      YearUnit,
		wantLen:   3,
		wantFirst: time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC),
		wantLast:  time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
	}, {
		name:      "minute range",
		start:     time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
		end:       time.Date(2023, 10, 1, 10, 32, 0, 0, time.UTC),
		unit:      MinuteUnit,
		wantLen:   3,
		wantFirst: time.Date(2023, 10, 1, 10, 30, 0, 0, time.UTC),
		wantLast:  time.Date(2023, 10, 1, 10, 32, 0, 0, time.UTC),
	}, {
		name:      "second range",
		start:     time.Date(2023, 10, 1, 10, 30, 45, 0, time.UTC),
		end:       time.Date(2023, 10, 1, 10, 30, 47, 0, time.UTC),
		unit:      SecondUnit,
		wantLen:   3,
		wantFirst: time.Date(2023, 10, 1, 10, 30, 45, 0, time.UTC),
		wantLast:  time.Date(2023, 10, 1, 10, 30, 47, 0, time.UTC),
	}, {
		name:      "month end rollover",
		start:     time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
		end:       time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
		unit:      MonthUnit,
		wantLen:   2,
		wantFirst: time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
		wantLast:  time.Date(2023, 4, 30, 0, 0, 0, 0, time.UTC),
	}, {
		name:      "single day",
		start:     start,
		end:       start,
		unit:      DayUnit,
		wantLen:   1,
		wantFirst: start,
		wantLast:  start,
	}, {
		name:      "empty range",
		start:     end,
		end:       start,
		unit:      DayUnit,
		wantLen:   0,
		wantFirst: time.Time{},
		wantLast:  time.Time{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeGen := Range(tt.start, tt.end, tt.unit)
			result := rangeGen.Generate()

			if len(result) != tt.wantLen {
				t.Errorf("Generate() length = %v, want %v", len(result), tt.wantLen)
				return
			}

			if tt.wantLen > 0 {
				if !result[0].Equal(tt.wantFirst) {
					t.Errorf("Generate() first = %v, want %v", result[0], tt.wantFirst)
				}
				if !result[len(result)-1].Equal(tt.wantLast) {
					t.Errorf("Generate() last = %v, want %v", result[len(result)-1], tt.wantLast)
				}
			}
		})
	}
}

func TestRangeContains(t *testing.T) {
	range1 := Range(
		time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	range2 := Range(
		time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	range3 := Range(
		time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 10, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	range4 := Range(
		time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	range5 := Range(
		time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		HourUnit,
	)

	tests := []struct {
		name      string
		start     *DateRange
		end       *DateRange
		wantLen   int
		wantDates []time.Time
	}{{
		name:    "partial overlap",
		start:   range1,
		end:     range2,
		wantLen: 3,
		wantDates: []time.Time{
			time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		},
	}, {
		name:    "full overlap",
		start:   range1,
		end:     range4,
		wantLen: 2,
		wantDates: []time.Time{
			time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
		},
	}, {
		name:      "no overlap",
		start:     range1,
		end:       range3,
		wantLen:   0,
		wantDates: []time.Time{},
	}, {
		name:      "different units no overlap",
		start:     range1,
		end:       range5,
		wantLen:   0,
		wantDates: []time.Time{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RangeContains(tt.start, tt.end)
			if len(result) != tt.wantLen {
				t.Errorf("RangeContains() length = %v, want %v", len(result), tt.wantLen)
				return
			}

			for i, date := range tt.wantDates {
				if !result[i].Equal(date) {
					t.Errorf("RangeContains() date[%d] = %v, want %v", i, result[i], date)
				}
			}
		})
	}
}

func TestRangeNotContains(t *testing.T) {
	rangeA := Range(
		time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	rangeB := Range(
		time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	rangeC := Range(
		time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC),
		DayUnit,
	)
	rangeD := Range(
		time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		HourUnit,
	)

	tests := []struct {
		name      string
		start     *DateRange
		end       *DateRange
		wantLen   int
		wantDates []time.Time
	}{{
		name:    "partial difference",
		start:   rangeA,
		end:     rangeB,
		wantLen: 2,
		wantDates: []time.Time{
			time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		},
	}, {
		name:    "no overlap difference",
		start:   rangeA,
		end:     rangeC,
		wantLen: 3,
		wantDates: []time.Time{
			time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC),
		},
	}, {
		name:      "different units full difference",
		start:     rangeA,
		end:       rangeD,
		wantLen:   0,
		wantDates: []time.Time{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RangeNotContains(tt.start, tt.end)
			if len(result) != tt.wantLen {
				t.Errorf("RangeNotContains() length = %v, want %v", len(result), tt.wantLen)
				return
			}

			// 只验证非空结果的具体日期
			if tt.wantLen > 0 && len(tt.wantDates) > 0 {
				for i, date := range tt.wantDates {
					if !result[i].Equal(date) {
						t.Errorf("RangeNotContains() date[%d] = %v, want %v", i, result[i], date)
					}
				}
			}
		})
	}
}

func TestHour(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 15,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hour(tt.t); got != tt.want {
				t.Errorf("Hour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinute(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 30,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Minute(tt.t); got != tt.want {
				t.Errorf("Minute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecond(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want int
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		want: 45,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Second(tt.t); got != tt.want {
				t.Errorf("Second() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want bool
	}{{
		name: "saturday",
		t:    time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "sunday",
		t:    time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "monday",
		t:    time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC),
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeekend(tt.t); got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddYears(t *testing.T) {
	tests := []struct {
		name  string
		t     time.Time
		years int
		want  time.Time
	}{{
		name:  "add positive years",
		t:     time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		years: 2,
		want:  time.Date(2025, 10, 5, 15, 30, 45, 0, time.UTC),
	}, {
		name:  "add negative years",
		t:     time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		years: -1,
		want:  time.Date(2022, 10, 5, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddYears(tt.t, tt.years); !got.Equal(tt.want) {
				t.Errorf("AddYears() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddMonths(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		months int
		want   time.Time
	}{{
		name:   "add positive months",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		months: 3,
		want:   time.Date(2024, 1, 5, 15, 30, 45, 0, time.UTC),
	}, {
		name:   "add negative months",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		months: -2,
		want:   time.Date(2023, 8, 5, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddMonths(tt.t, tt.months); !got.Equal(tt.want) {
				t.Errorf("AddMonths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		days int
		want time.Time
	}{{
		name: "add positive days",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		days: 3,
		want: time.Date(2023, 10, 8, 15, 30, 45, 0, time.UTC),
	}, {
		name: "add negative days",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		days: -2,
		want: time.Date(2023, 10, 3, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddDays(tt.t, tt.days); !got.Equal(tt.want) {
				t.Errorf("AddDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddHours(t *testing.T) {
	tests := []struct {
		name  string
		t     time.Time
		hours int
		want  time.Time
	}{{
		name:  "add positive hours",
		t:     time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		hours: 3,
		want:  time.Date(2023, 10, 5, 18, 30, 45, 0, time.UTC),
	}, {
		name:  "add negative hours",
		t:     time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		hours: -2,
		want:  time.Date(2023, 10, 5, 13, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddHours(tt.t, tt.hours); !got.Equal(tt.want) {
				t.Errorf("AddHours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddMinutes(t *testing.T) {
	tests := []struct {
		name    string
		t       time.Time
		minutes int
		want    time.Time
	}{{
		name:    "add positive minutes",
		t:       time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		minutes: 30,
		want:    time.Date(2023, 10, 5, 16, 0, 45, 0, time.UTC),
	}, {
		name:    "add negative minutes",
		t:       time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		minutes: -20,
		want:    time.Date(2023, 10, 5, 15, 10, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddMinutes(tt.t, tt.minutes); !got.Equal(tt.want) {
				t.Errorf("AddMinutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddSeconds(t *testing.T) {
	tests := []struct {
		name    string
		t       time.Time
		seconds int
		want    time.Time
	}{{
		name:    "add positive seconds",
		t:       time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		seconds: 30,
		want:    time.Date(2023, 10, 5, 15, 31, 15, 0, time.UTC),
	}, {
		name:    "add negative seconds",
		t:       time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		seconds: -20,
		want:    time.Date(2023, 10, 5, 15, 30, 25, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddSeconds(tt.t, tt.seconds); !got.Equal(tt.want) {
				t.Errorf("AddSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffYears(t *testing.T) {
	tests := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want int
	}{{
		name: "t1 after t2",
		t1:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		t2:   time.Date(2020, 10, 5, 0, 0, 0, 0, time.UTC),
		want: 3,
	}, {
		name: "t1 before t2",
		t1:   time.Date(2020, 10, 5, 0, 0, 0, 0, time.UTC),
		t2:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		want: -3,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffYears(tt.t1, tt.t2); got != tt.want {
				t.Errorf("DiffYears() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffDays(t *testing.T) {
	tests := []struct {
		name string
		t1   time.Time
		t2   time.Time
		want int
	}{{
		name: "t1 after t2",
		t1:   time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		t2:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		want: 3,
	}, {
		name: "t1 before t2",
		t1:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		t2:   time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		want: -3,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffDays(tt.t1, tt.t2); got != tt.want {
				t.Errorf("DiffDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfSecond(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfSecond(tt.date); !got.Equal(tt.want) {
				t.Errorf("BeginOfSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfSecond(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 30, 45, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfSecond(tt.date); !got.Equal(tt.want) {
				t.Errorf("EndOfSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfMinute(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 30, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfMinute(tt.date); !got.Equal(tt.want) {
				t.Errorf("BeginOfMinute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfMinute(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 30, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfMinute(tt.date); !got.Equal(tt.want) {
				t.Errorf("EndOfMinute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfHour(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfHour(tt.date); !got.Equal(tt.want) {
				t.Errorf("BeginOfHour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfHour(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 15, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfHour(tt.date); !got.Equal(tt.want) {
				t.Errorf("EndOfHour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfDay(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfDay(tt.date); !got.Equal(tt.want) {
				t.Errorf("BeginOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "normal case",
		date: time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 5, 23, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfDay(tt.date); !got.Equal(tt.want) {
				t.Errorf("EndOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfWeek(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "monday",
		date: time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC),
	}, {
		name: "sunday",
		date: time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfWeek(tt.date); !got.Equal(tt.want) {
				t.Errorf("BeginOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfWeekWithMondayStart(t *testing.T) {
	tests := []struct {
		name               string
		date               time.Time
		isMondayAsFirstDay bool
		want               time.Time
	}{{
		name:               "monday as first day",
		date:               time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		isMondayAsFirstDay: true,
		want:               time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
	}, {
		name:               "sunday as first day",
		date:               time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		isMondayAsFirstDay: false,
		want:               time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfWeekWithMondayStart(tt.date, tt.isMondayAsFirstDay); !got.Equal(tt.want) {
				t.Errorf("BeginOfWeekWithMondayStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfWeek(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want time.Time
	}{{
		name: "sunday",
		date: time.Date(2025, 6, 21, 0, 0, 0, 0, time.UTC),
		want: time.Date(2025, 6, 22, 23, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfWeek(tt.date); !got.Equal(tt.want) {
				t.Errorf("EndOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfWeekWithSundayEnd(t *testing.T) {
	tests := []struct {
		name              string
		date              time.Time
		isSundayAsLastDay bool
		want              time.Time
	}{{
		name:              "sunday as last day",
		date:              time.Date(2023, 10, 2, 0, 0, 0, 0, time.Local),
		isSundayAsLastDay: true,
		want:              time.Date(2023, 10, 8, 23, 59, 59, 999000000, time.Local),
	}, {
		name:              "saturday as last day",
		date:              time.Date(2025, 6, 20, 0, 0, 0, 0, time.Local),
		isSundayAsLastDay: false,
		want:              time.Date(2025, 6, 21, 23, 59, 59, 999000000, time.Local),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfWeekWithSundayEnd(tt.date, tt.isSundayAsLastDay); !got.Equal(tt.want) {
				t.Errorf("EndOfWeekWithSundayEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfMonth(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfMonth(tt.t); !got.Equal(tt.want) {
				t.Errorf("BeginOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfMonth(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "october",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 10, 31, 23, 59, 59, 999000000, time.UTC),
	}, {
		name: "february leap year",
		t:    time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2020, 2, 29, 23, 59, 59, 999000000, time.UTC),
	}, {
		name: "february non-leap year",
		t:    time.Date(2021, 2, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2021, 2, 28, 23, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfMonth(tt.t); !got.Equal(tt.want) {
				t.Errorf("EndOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfQuarter(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "q1",
		t:    time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}, {
		name: "q2",
		t:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
	}, {
		name: "q3",
		t:    time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC),
	}, {
		name: "q4",
		t:    time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfQuarter(tt.t); !got.Equal(tt.want) {
				t.Errorf("BeginOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfQuarter(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "q1",
		t:    time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 3, 31, 23, 59, 59, 999000000, time.UTC),
	}, {
		name: "q2",
		t:    time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 6, 30, 23, 59, 59, 999000000, time.UTC),
	}, {
		name: "q3",
		t:    time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 9, 30, 23, 59, 59, 999000000, time.UTC),
	}, {
		name: "q4",
		t:    time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC),
		want: time.Date(2023, 12, 31, 23, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfQuarter(tt.t); !got.Equal(tt.want) {
				t.Errorf("EndOfQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBeginOfYear(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeginOfYear(tt.t); !got.Equal(tt.want) {
				t.Errorf("BeginOfYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfYear(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{{
		name: "normal case",
		t:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: time.Date(2023, 12, 31, 23, 59, 59, 999000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfYear(tt.t); !got.Equal(tt.want) {
				t.Errorf("EndOfYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYesterday(t *testing.T) {
	today := time.Now()
	yesterday := BeginOfDay(today.AddDate(0, 0, -1))
	got := Yesterday()
	if !got.Equal(yesterday) {
		t.Errorf("Yesterday() = %v, want %v", got, yesterday)
	}
}

func TestTomorrow(t *testing.T) {
	today := time.Now()
	tomorrow := BeginOfDay(today.AddDate(0, 0, 1))
	got := Tomorrow()
	if !got.Equal(tomorrow) {
		t.Errorf("Tomorrow() = %v, want %v", got, tomorrow)
	}
}

func TestLastWeek(t *testing.T) {
	expected := time.Now().AddDate(0, 0, -7)
	got := LastWeek()
	if got.After(expected.Add(1*time.Second)) || got.Before(expected.Add(-1*time.Second)) {
		t.Errorf("LastWeek() = %v, want around %v", got, expected)
	}
}

func TestNextWeek(t *testing.T) {
	expected := time.Now().AddDate(0, 0, 7)
	got := NextWeek()
	if got.After(expected.Add(1*time.Second)) || got.Before(expected.Add(-1*time.Second)) {
		t.Errorf("NextWeek() = %v, want around %v", got, expected)
	}
}

func TestLastMonth(t *testing.T) {
	expected := time.Now().AddDate(0, -1, 0)
	got := LastMonth()
	if got.After(expected.Add(1*time.Minute)) || got.Before(expected.Add(-1*time.Minute)) {
		t.Errorf("LastMonth() = %v, want around %v", got, expected)
	}
}

func TestNextMonth(t *testing.T) {
	expected := time.Now().AddDate(0, 1, 0)
	got := NextMonth()
	if got.After(expected.Add(1*time.Minute)) || got.Before(expected.Add(-1*time.Minute)) {
		t.Errorf("NextMonth() = %v, want around %v", got, expected)
	}
}

func TestOffsetMillisecond(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 500,
		want:   time.Date(2023, 10, 5, 15, 30, 45, 500000000, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 500000000, time.UTC),
		offset: -200,
		want:   time.Date(2023, 10, 5, 15, 30, 45, 300000000, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetMillisecond(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetMillisecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetSecond(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 30,
		want:   time.Date(2023, 10, 5, 15, 31, 15, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -20,
		want:   time.Date(2023, 10, 5, 15, 30, 25, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetSecond(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetMinute(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 30,
		want:   time.Date(2023, 10, 5, 16, 0, 45, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -20,
		want:   time.Date(2023, 10, 5, 15, 10, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetMinute(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetMinute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetHour(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 3,
		want:   time.Date(2023, 10, 5, 18, 30, 45, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -2,
		want:   time.Date(2023, 10, 5, 13, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetHour(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetHour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetDay(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 3,
		want:   time.Date(2023, 10, 8, 15, 30, 45, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -2,
		want:   time.Date(2023, 10, 3, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetDay(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetWeek(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 1,
		want:   time.Date(2023, 10, 12, 15, 30, 45, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -1,
		want:   time.Date(2023, 9, 28, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetWeek(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetMonth(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		offset int
		want   time.Time
	}{{
		name:   "positive offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: 3,
		want:   time.Date(2024, 1, 5, 15, 30, 45, 0, time.UTC),
	}, {
		name:   "negative offset",
		t:      time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		offset: -2,
		want:   time.Date(2023, 8, 5, 15, 30, 45, 0, time.UTC),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OffsetMonth(tt.t, tt.offset); !got.Equal(tt.want) {
				t.Errorf("OffsetMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	tests := []struct {
		name  string
		begin time.Time
		end   time.Time
		unit  TimeUnit
		isAbs bool
		want  int64
	}{{
		name:  "milliseconds",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 5, 15, 30, 45, 500000000, time.UTC),
		unit:  Millisecond,
		isAbs: true,
		want:  500,
	}, {
		name:  "seconds",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 5, 15, 31, 15, 0, time.UTC),
		unit:  SecondUnit,
		isAbs: true,
		want:  30,
	}, {
		name:  "minutes",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 5, 16, 0, 45, 0, time.UTC),
		unit:  MinuteUnit,
		isAbs: true,
		want:  30,
	}, {
		name:  "hours",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 5, 18, 30, 45, 0, time.UTC),
		unit:  HourUnit,
		isAbs: true,
		want:  3,
	}, {
		name:  "days",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 8, 15, 30, 45, 0, time.UTC),
		unit:  DayUnit,
		isAbs: true,
		want:  3,
	}, {
		name:  "weeks",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 12, 15, 30, 45, 0, time.UTC),
		unit:  WeekUnit,
		isAbs: true,
		want:  1,
	}, {
		name:  "months",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2024, 1, 5, 15, 30, 45, 0, time.UTC),
		unit:  MonthUnit,
		isAbs: true,
		want:  3,
	}, {
		name:  "years",
		begin: time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2025, 10, 5, 15, 30, 45, 0, time.UTC),
		unit:  YearUnit,
		isAbs: true,
		want:  2,
	}, {
		name:  "negative without abs",
		begin: time.Date(2023, 10, 8, 15, 30, 45, 0, time.UTC),
		end:   time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		unit:  DayUnit,
		isAbs: false,
		want:  -3,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Between(tt.begin, tt.end, tt.unit, tt.isAbs); got != tt.want {
				t.Errorf("Between() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenMs(t *testing.T) {
	begin := time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC)
	end := time.Date(2023, 10, 5, 15, 30, 45, 500000000, time.UTC)
	if got := BetweenMs(begin, end); got != 500 {
		t.Errorf("BetweenMs() = %v, want %v", got, 500)
	}
}

func TestBetweenDay(t *testing.T) {
	tests := []struct {
		name    string
		begin   time.Time
		end     time.Time
		isReset bool
		want    int64
	}{{
		name:    "with reset",
		begin:   time.Date(2023, 10, 5, 23, 59, 59, 0, time.UTC),
		end:     time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    2,
	}, {
		name:    "without reset",
		begin:   time.Date(2023, 10, 5, 23, 59, 59, 0, time.UTC),
		end:     time.Date(2023, 10, 7, 0, 0, 0, 0, time.UTC),
		isReset: false,
		want:    1,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BetweenDay(tt.begin, tt.end, tt.isReset); got != tt.want {
				t.Errorf("BetweenDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenWeek(t *testing.T) {
	tests := []struct {
		name    string
		begin   time.Time
		end     time.Time
		isReset bool
		want    int64
	}{{
		name:    "with reset",
		begin:   time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    2,
	}, {
		name:    "without reset",
		begin:   time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 10, 16, 12, 0, 0, 0, time.UTC),
		isReset: false,
		want:    2,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BetweenWeek(tt.begin, tt.end, tt.isReset); got != tt.want {
				t.Errorf("BetweenWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenMonth(t *testing.T) {
	tests := []struct {
		name    string
		begin   time.Time
		end     time.Time
		isReset bool
		want    int64
	}{{
		name:    "with reset",
		begin:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    2,
	}, {
		name:    "without reset",
		begin:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		isReset: false,
		want:    2,
	}, {
		name:    "end before begin",
		begin:   time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    2,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BetweenMonth(tt.begin, tt.end, tt.isReset); got != tt.want {
				t.Errorf("BetweenMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBetweenYear(t *testing.T) {
	tests := []struct {
		name    string
		begin   time.Time
		end     time.Time
		isReset bool
		want    int64
	}{{
		name:    "with reset",
		begin:   time.Date(2020, 10, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    3,
	}, {
		name:    "without reset",
		begin:   time.Date(2020, 10, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		isReset: false,
		want:    3,
	}, {
		name:    "end before begin",
		begin:   time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		end:     time.Date(2020, 10, 5, 0, 0, 0, 0, time.UTC),
		isReset: true,
		want:    3,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BetweenYear(tt.begin, tt.end, tt.isReset); got != tt.want {
				t.Errorf("BetweenYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameTime(t *testing.T) {
	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{{
		name: "same time",
		a:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		b:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		want: true,
	}, {
		name: "different time",
		a:    time.Date(2023, 10, 5, 15, 30, 45, 123456789, time.UTC),
		b:    time.Date(2023, 10, 5, 15, 30, 45, 987654321, time.UTC),
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameTime(tt.a, tt.b); got != tt.want {
				t.Errorf("IsSameTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{{
		name: "same day",
		a:    time.Date(2023, 10, 5, 15, 30, 45, 0, time.UTC),
		b:    time.Date(2023, 10, 5, 23, 59, 59, 0, time.UTC),
		want: true,
	}, {
		name: "different day",
		a:    time.Date(2023, 10, 5, 23, 59, 59, 0, time.UTC),
		b:    time.Date(2023, 10, 6, 0, 0, 0, 0, time.UTC),
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameDay(tt.a, tt.b); got != tt.want {
				t.Errorf("IsSameDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameWeek(t *testing.T) {
	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{{
		name: "same week",
		a:    time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC),
		b:    time.Date(2023, 10, 13, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "different week",
		a:    time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		b:    time.Date(2023, 10, 9, 0, 0, 0, 0, time.UTC),
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameWeek(tt.a, tt.b); got != tt.want {
				t.Errorf("IsSameWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameMonth(t *testing.T) {
	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{{
		name: "same month",
		a:    time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		b:    time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "different month",
		a:    time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		b:    time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameMonth(tt.a, tt.b); got != tt.want {
				t.Errorf("IsSameMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		name string
		year int
		want bool
	}{{
		name: "leap year divisible by 400",
		year: 2000,
		want: true,
	}, {
		name: "not leap year divisible by 100 but not 400",
		year: 1900,
		want: false,
	}, {
		name: "leap year divisible by 4 but not 100",
		year: 2020,
		want: true,
	}, {
		name: "not leap year not divisible by 4",
		year: 2021,
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLeapYear(tt.year); got != tt.want {
				t.Errorf("IsLeapYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgeOfNow(t *testing.T) {
	// 模拟当前时间为2023-10-05
	now := time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)
	// 使用匿名函数包装AgeOfNow，以便注入当前时间进行测试
	ageOfNow := func(birthDay time.Time) int {
		return age(birthDay, now)
	}

	tests := []struct {
		name     string
		birthDay time.Time
		want     int
	}{{
		name:     "birthday passed",
		birthDay: time.Date(2000, 5, 15, 0, 0, 0, 0, time.UTC),
		want:     23,
	}, {
		name:     "birthday today",
		birthDay: time.Date(2000, 10, 5, 0, 0, 0, 0, time.UTC),
		want:     23,
	}, {
		name:     "birthday not passed",
		birthDay: time.Date(2000, 12, 25, 0, 0, 0, 0, time.UTC),
		want:     22,
	}, {
		name:     "future birthday",
		birthDay: time.Date(2030, 10, 5, 0, 0, 0, 0, time.UTC),
		want:     0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ageOfNow(tt.birthDay); got != tt.want {
				t.Errorf("AgeOfNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgeOfNowString(t *testing.T) {
	tests := []struct {
		name        string
		birthDayStr string
		want        int
		wantErr     bool
	}{{
		name:        "valid date",
		birthDayStr: "2000-05-15",
		want:        25,
		wantErr:     false,
	}, {
		name:        "invalid format",
		birthDayStr: "2000/05/15",
		want:        0,
		wantErr:     true,
	}, {
		name:        "empty string",
		birthDayStr: "",
		want:        0,
		wantErr:     true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AgeOfNowString(tt.birthDayStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("AgeOfNowString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AgeOfNowString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChineseZodiac(t *testing.T) {
	tests := []struct {
		name string
		year int
		want string
	}{{
		name: "鼠",
		year: 1900,
		want: "鼠",
	}, {
		name: "牛",
		year: 1901,
		want: "牛",
	}, {
		name: "虎",
		year: 1902,
		want: "虎",
	}, {
		name: "兔",
		year: 1903,
		want: "兔",
	}, {
		name: "龙",
		year: 1904,
		want: "龙",
	}, {
		name: "蛇",
		year: 1905,
		want: "蛇",
	}, {
		name: "马",
		year: 1906,
		want: "马",
	}, {
		name: "羊",
		year: 1907,
		want: "羊",
	}, {
		name: "猴",
		year: 1908,
		want: "猴",
	}, {
		name: "鸡",
		year: 1909,
		want: "鸡",
	}, {
		name: "狗",
		year: 1910,
		want: "狗",
	}, {
		name: "猪",
		year: 1995,
		want: "猪",
	}, {
		name: "year before 1900",
		year: 1899,
		want: "",
	}, {
		name: "negative year",
		year: -1,
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChineseZodiac(tt.year); got != tt.want {
				t.Errorf("GetChineseZodiac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLengthOfYear(t *testing.T) {
	tests := []struct {
		name string
		year int
		want int
	}{{
		name: "leap year",
		year: 2020,
		want: 366,
	}, {
		name: "non-leap year",
		year: 2021,
		want: 365,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LengthOfYear(tt.year); got != tt.want {
				t.Errorf("LengthOfYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLengthOfMonth(t *testing.T) {
	tests := []struct {
		name       string
		month      int
		isLeapYear bool
		want       int
	}{{
		name:       "january",
		month:      1,
		isLeapYear: false,
		want:       31,
	}, {
		name:       "february leap",
		month:      2,
		isLeapYear: true,
		want:       29,
	}, {
		name:       "february non-leap",
		month:      2,
		isLeapYear: false,
		want:       28,
	}, {
		name:       "april",
		month:      4,
		isLeapYear: false,
		want:       30,
	}, {
		name:       "invalid month",
		month:      13,
		isLeapYear: false,
		want:       0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LengthOfMonth(tt.month, tt.isLeapYear); got != tt.want {
				t.Errorf("LengthOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want bool
	}{{
		name: "last day of month",
		date: time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "not last day",
		date: time.Date(2023, 10, 30, 0, 0, 0, 0, time.UTC),
		want: false,
	}, {
		name: "leap year february",
		date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
		want: true,
	}, {
		name: "non-leap year february",
		date: time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC),
		want: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLastDayOfMonth(tt.date); got != tt.want {
				t.Errorf("IsLastDayOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want int
	}{{
		name: "october",
		date: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC),
		want: 31,
	}, {
		name: "april",
		date: time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
		want: 30,
	}, {
		name: "leap year february",
		date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		want: 29,
	}, {
		name: "non-leap year february",
		date: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
		want: 28,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLastDayOfMonth(tt.date); got != tt.want {
				t.Errorf("GetLastDayOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
