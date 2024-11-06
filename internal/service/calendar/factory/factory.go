package factory

import (
	"context"
	"time"

	"github.com/yanodincov/skyeng-ics/internal/config"
	ics2 "github.com/yanodincov/skyeng-ics/pkg/ics"

	ics "github.com/arran4/golang-ical"
	"github.com/yanodincov/skyeng-ics/internal/repository/skyeng/model"
)

const EventPrefix = "[Skyeng] "

type Factory struct {
	cfg *config.Config
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		cfg: cfg,
	}
}

func (f *Factory) CreateCalendarFromLessons(_ context.Context, lessons []model.Lesson) (*ics.Calendar, error) {
	calendar := ics.NewCalendar()
	calendar.SetMethod(ics.MethodPublish)
	calendar.SetProductId("Skyeng ICS")
	calendar.SetName("Skyeng Calendar")

	now := time.Now()
	createdAt := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

	for _, lesson := range lessons {
		event := calendar.AddEvent(lesson.LessonID)
		event.SetSummary(EventPrefix + lesson.Teacher.Name + " lesson")
		event.SetCreatedTime(createdAt)
		event.SetDtStampTime(createdAt)
		event.SetModifiedAt(createdAt)
		event.SetStartAt(lesson.StartAt)
		event.SetEndAt(lesson.EndAt)
		event.SetLocation("skyeng.com")
		event.SetDescription(lesson.EducationService.Title + " - " + lesson.Teacher.Name)
		event.SetURL("https://skyeng.ru")
		event.SetOrganizer("Skyeng", ics.WithCN("Skyeng ICS Maker"))

		for _, notifyTimeDur := range f.cfg.Calendar.NotifyTimeList.Vals() {
			alarm := event.AddAlarm()
			alarm.SetAction(ics.ActionDisplay)
			alarm.SetDescription("Skyeng lesson reminder")
			alarm.SetSummary("Skyeng lesson starts in " + notifyTimeDur.String())
			alarm.SetTrigger(ics2.ConvertDurationToICS(notifyTimeDur))
		}
	}

	return calendar, nil
}
