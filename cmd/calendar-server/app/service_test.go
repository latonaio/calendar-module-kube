package app

import (
	"context"
	"testing"
	"time"

	"bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server/proto/calendarpb"
	"github.com/golang/protobuf/ptypes"
)

func TestCreateScheduleError(t *testing.T) {
	service := &CalendarService{}
	ctx := context.Background()

	t.Run("without date", func(t *testing.T) {
		schedule := &calendarpb.Schedule{}
		_, err := service.CreateSchedule(ctx, schedule)
		if err == nil {
			t.Fatalf("failed test: %v", err)
		}
		if err.Error() != "not exist Date" {
			t.Fatalf("failed test: %v", err)
		}

	})

	t.Run("without StartDate", func(t *testing.T) {
		end, err := time.Parse(time.RFC3339, "2020-01-01T00:00:00+09:00")
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		pend, err := ptypes.TimestampProto(end)
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}

		schedule := &calendarpb.Schedule{Date: &calendarpb.Date{End: pend}}
		_, err = service.CreateSchedule(ctx, schedule)
		if err == nil {
			t.Fatalf("failed test: %v", err)
		}

		if err.Error() != "not exist Start in Date" {
			t.Fatalf("failed test: %v", err)
		}
	})

	t.Run("without EndDate", func(t *testing.T) {
		start, err := time.Parse(time.RFC3339, "2020-01-02T00:00:00+09:00")
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		pstart, err := ptypes.TimestampProto(start)
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		schedule := &calendarpb.Schedule{Date: &calendarpb.Date{Start: pstart}}
		_, err = service.CreateSchedule(ctx, schedule)
		if err == nil {
			t.Fatalf("failed test: %v", err)
		}
		if err.Error() != "not exist End in Date" {
			t.Fatalf("failed test: %v", err)
		}
	})

	t.Run("StartDate later than EndDate", func(t *testing.T) {
		start, err := time.Parse(time.RFC3339, "2020-01-03T00:00:00+09:00")
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		pstart, err := ptypes.TimestampProto(start)
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		end, err := time.Parse(time.RFC3339, "2019-01-03T00:00:00+09:00")
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		pend, err := ptypes.TimestampProto(end)
		if err != nil {
			t.Fatalf("test time invalid: %v", err)
		}
		schedule := &calendarpb.Schedule{
			Date: &calendarpb.Date{Start: pstart, End: pend}}
		_, err = service.CreateSchedule(ctx, schedule)
		if err == nil {
			t.Fatalf("failed test: %v", err)
		}
		if err.Error() != "StartDate later than EndDate" {
			t.Fatalf("failed test: %v", err)
		}
	})
}
