package tariff

import (
	"k071123/internal/services/parking_service/domain/models"
	"testing"
	"time"
)

func Test_calculateSessionCost(t *testing.T) {
	min15 := 15
	hasFree := true
	type args struct {
		duration time.Duration
		tariff   models.Tariff
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// 1. Бесплатная парковка (FreeTime)
		{
			name: "free_within_free_time",
			args: args{
				duration: 15 * time.Minute,
				tariff: models.Tariff{
					HasFree:     &hasFree,
					FreeTime:    &min15,
					HourlyPrice: 100,
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "free_after_free_time", // done
			args: args{
				duration: 45 * time.Minute,
				tariff: models.Tariff{
					HasFree:     new(bool),
					FreeTime:    func() *int { i := 30; return &i }(),
					HourlyPrice: 100,
				},
			},
			want:    100, // час
			wantErr: false,
		},

		// 2. Менее часа (базовый тариф)
		{
			name: "less_than_one_hour",
			args: args{
				duration: 40 * time.Minute,
				tariff: models.Tariff{
					HourlyPrice: 150,
				},
			},
			want:    150,
			wantErr: false,
		},

		// 3. От 1 до LongHourlyStart (пропорционально часам)
		{
			name: "between_1h_and_long_start",
			args: args{
				duration: 2 * time.Hour,
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180, // 3 часа
					LongHourlyPrice: 200,
					DailyPrice:      1000,
				},
			},
			want:    200, // 2 ч × 100
			wantErr: false,
		},

		// 4. От LongHourlyStart до LongHourlyEnd (фиксированная цена)
		{
			name: "between_long_start_and_long_end",
			args: args{
				duration: 4 * time.Hour,
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180, // 3 ч
					LongHourlyEnd:   360, // 6 ч
					LongHourlyPrice: 500,
					DailyPrice:      1000,
				},
			},
			want:    500, // фиксированная цена
			wantErr: false,
		},

		// 5. После LongHourlyEnd, но менее 48 часов (суточная цена)
		{
			name: "after_long_end_less_than_48h",
			args: args{
				duration: 12 * time.Hour,
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180,
					LongHourlyEnd:   360,
					LongHourlyPrice: 500,
					DailyPrice:      800,
				},
			},
			want:    800, // суточная цена
			wantErr: false,
		},

		// 6. 48 часов и более (пропорционально суткам)
		{
			name: "48_hours_or_more",
			args: args{
				duration: 72 * time.Hour, // 3 суток
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180,
					LongHourlyEnd:   360,
					LongHourlyPrice: 500,
					DailyPrice:      1000,
				},
			},
			want:    3000, // 3 × 1000
			wantErr: false,
		},

		// 7. Граничный случай: ровно LongHourlyStart
		{
			name: "exactly_long_hourly_start",
			args: args{
				duration: 3 * time.Hour, // 180 мин
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180,
					LongHourlyEnd:   360,
					LongHourlyPrice: 600,
					DailyPrice:      1200,
				},
			},
			want:    600, // LongHourlyPrice
			wantErr: false,
		},

		// 8. Граничный случай: ровно LongHourlyEnd
		{
			name: "exactly_long_hourly_end",
			args: args{
				duration: 6 * time.Hour, // 360 мин
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180,
					LongHourlyEnd:   360,
					LongHourlyPrice: 700,
					DailyPrice:      1400,
				},
			},
			want:    700, // LongHourlyPrice
			wantErr: false,
		},

		// 9. Граничный случай: ровно 48 часов
		{
			name: "exactly_48_hours",
			args: args{
				duration: 48 * time.Hour,
				tariff: models.Tariff{
					HourlyPrice:     100,
					LongHourlyStart: 180,
					LongHourlyEnd:   360,
					LongHourlyPrice: 800,
					DailyPrice:      1600,
				},
			},
			want:    3200, // 2 × 1600
			wantErr: false,
		},

		// 10. Крайний случай: нулевая длительность
		{
			name: "zero_duration",
			args: args{
				duration: 0,
				tariff: models.Tariff{
					HourlyPrice: 100,
				},
			},
			want:    0,
			wantErr: false,
		},

		// 11. Крайний случай: tariff.HasFree == nil (нет бесплатной парковки)
		{
			name: "no_free_time_defined",
			args: args{
				duration: 20 * time.Minute,
				tariff: models.Tariff{
					HasFree:     nil, // не задано
					FreeTime:    nil,
					HourlyPrice: 120,
				},
			},
			want:    120, // сразу берём часовой тариф
			wantErr: false,
		},

		// 12. Крайний случай: FreeTime == 0 (бесплатно 0 минут)
		{
			name: "free_time_zero",
			args: args{
				duration: 10 * time.Minute,
				tariff: models.Tariff{
					HasFree:     new(bool),
					FreeTime:    func() *int { i := 0; return &i }(),
					HourlyPrice: 150,
				},
			},
			want:    150, // сразу часовой тариф
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateSessionCost(tt.args.duration, tt.args.tariff)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateSessionCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateSessionCost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
