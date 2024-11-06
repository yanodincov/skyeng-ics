package model

import "time"

type Schedule struct {
	Data []Lesson `json:"data"`
}

type Lesson struct {
	LessonID         string           `json:"lessonId"`
	StartAt          time.Time        `json:"startAt"`
	EndAt            time.Time        `json:"endAt"`
	EducationService EducationService `json:"educationService"`
	Teacher          Teacher          `json:"teacher"`
}

type Teacher struct {
	Name string `json:"name"`
}

type EducationService struct {
	ServiceTypeKey string `json:"serviceTypeKey"`
	Title          string `json:"title"`
}
