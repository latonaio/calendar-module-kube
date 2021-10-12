package app

import (
	"context"
	"errors"
	"log"
	"time"

	"bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server/proto/calendarpb"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CalendarService struct {
	Db *DB
}

func (c *CalendarService) CreateSchedule(ctx context.Context, schedule *calendarpb.Schedule) (*calendarpb.ResponseSchedule, error) {

	log.Printf("[info] create schedule %+v\n", schedule)

	if schedule.Date == nil {
		log.Println("[error] not exist Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Date")
	}

	start, err := ptypes.Timestamp(schedule.Date.Start)
	if err != nil {
		log.Println("[error] not exist Start in Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Start in Date")
	}

	end, err := ptypes.Timestamp(schedule.Date.End)
	if err != nil {
		log.Println("[error] not exist End in Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist End in Date")
	}

	if start.Unix() > end.Unix() {
		log.Println("[error] StartDate later than EndDate")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("StartDate later than EndDate")
	}

	if schedule.Title == "" {
		log.Println("[error] not exist Title")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Title")
	}

	tagList := make([]TagMaster, len(schedule.TagList))
	for i, v := range schedule.TagList {
		tagList[i].TagId = v.TagId
		tagList[i].TagName = v.TagName
	}

	res, err := c.Db.CreateSchedule(&Schedule{
		ScheduleId:  schedule.ScheduleId,
		StartDate:   start,
		EndDate:     end,
		Title:       schedule.Title,
		Description: schedule.Description,
		UserId:      schedule.User.UserId,
		UserName:    schedule.User.UserName,
	}, tagList)

	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeSchedule(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ResponseSchedule{Schedule: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) CreateTag(ctx context.Context, tag *calendarpb.Tag) (*calendarpb.ResponseTag, error) {
	log.Printf("[info] create tag %+v\n", tag)
	if tag == nil {
		log.Println("[error] not exist Tag")
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Tag")
	}

	if tag.TagName == "" {
		log.Println("[error] tag name is empty")
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("tag name is empty")
	}
	res, err := c.Db.CreateTag(&TagMaster{TagName: tag.TagName})
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeTag(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ResponseTag{Tag: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) DeleteSchedule(ctx context.Context, schedule *calendarpb.Schedule) (*calendarpb.ResponseSchedule, error) {
	log.Printf("[info] delete schedule %+v\n", schedule)
	err := c.Db.DeleteSchedule(&Schedule{ScheduleId: schedule.ScheduleId})
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseSchedule{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ResponseSchedule{StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) UpdateSchedule(ctx context.Context, schedule *calendarpb.Schedule) (*calendarpb.ResponseSchedule, error) {

	log.Printf("[info] Update schedule %+v\n", schedule)

	if schedule.Date == nil {
		log.Println("[error] not exist Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Date")
	}

	start, err := ptypes.Timestamp(schedule.Date.Start)
	if err != nil {
		log.Println("[error] not exist Start in Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Start in Date")
	}

	end, err := ptypes.Timestamp(schedule.Date.End)
	if err != nil {
		log.Println("[error] not exist End in Date")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist End in Date")
	}

	if start.Unix() > end.Unix() {
		log.Println("[error] StartDate later than EndDate")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("StartDate later than EndDate")
	}

	if schedule.Title == "" {
		log.Println("[error] not exist Title")
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Title")
	}
	tagList := make([]TagMaster, len(schedule.TagList))
	for i, v := range schedule.TagList {
		tagList[i].TagId = v.TagId
		tagList[i].TagName = v.TagName
	}

	res, err := c.Db.UpdateSchedule(&Schedule{
		ScheduleId:  schedule.ScheduleId,
		StartDate:   start,
		EndDate:     end,
		Title:       schedule.Title,
		Description: schedule.Description,
		UserId:      schedule.User.UserId,
		UserName:    schedule.User.UserName,
	}, tagList)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeSchedule(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseSchedule{Schedule: schedule, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ResponseSchedule{Schedule: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) GetScheduleList(ctx context.Context, e *emptypb.Empty) (*calendarpb.ScheduleList, error) {
	log.Printf("[info] get schedule list\n")
	res, err := c.Db.GetScheduleList()
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeScheduleList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ScheduleList{ScheduleList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) SearchScheduleByUserId(ctx context.Context, user *calendarpb.User) (*calendarpb.ScheduleList, error) {
	log.Printf("[info] search schedule by user id %+v\n", user)
	if user == nil {
		log.Println("[error] not exist User")
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist User")
	}
	res, err := c.Db.SearchScheduleByUserId(user.UserId)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeScheduleList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ScheduleList{ScheduleList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) SearchScheduleByUserName(ctx context.Context, user *calendarpb.User) (*calendarpb.ScheduleList, error) {
	log.Printf("[info] search schedule by user name %+v\n", user)
	if user == nil {
		log.Println("[error] not exist User")
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist User")
	}
	res, err := c.Db.SearchScheduleByUserName(user.UserName)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeScheduleList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ScheduleList{ScheduleList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) SearchScheduleByTagName(ctx context.Context, tag *calendarpb.Tag) (*calendarpb.ScheduleList, error) {
	log.Printf("[info] search schedule by tag name %+v\n", tag)
	if tag == nil {
		log.Println("[error] not exist Tag")
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Tag")
	}
	res, err := c.Db.SearchScheduleByTagName(tag.TagName)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeScheduleList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ScheduleList{ScheduleList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) SearchScheduleByDate(ctx context.Context, date *calendarpb.Date) (*calendarpb.ScheduleList, error) {
	log.Printf("[info] search schedule by date %+v\n", date)

	if date == nil {
		log.Println("[error] not exist Date")
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Date")
	}

	res := []Schedule{}
	var start time.Time
	var end time.Time
	var err error
	if date.Start != nil && date.End != nil {
		if start, err = ptypes.Timestamp(date.Start); err != nil {
			log.Println("[error] invalid timestamp: start")
			return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("invalid timestamp: start")
		}
		if end, err = ptypes.Timestamp(date.End); err != nil {
			log.Println("[error] invalid timestamp: end")
			return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("invalid timestamp: end")
		}
		res, err = c.Db.SearchScheduleDurationByDate(&start, &end)
	} else if date.Start != nil {
		if start, err = ptypes.Timestamp(date.Start); err != nil {
			log.Println("[error] invalid timestamp: start")
			return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("invalid timestamp: start")
		}
		res, err = c.Db.SearchScheduleSinceByDate(&start)
	} else if date.End != nil {
		if end, err = ptypes.Timestamp(date.End); err != nil {
			log.Println("[error] invalid timestamp: end")
			return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("invalid timestamp: end")
		}
		res, err = c.Db.SearchScheduleUntileByDate(&end)
	} else {
		log.Println("[error] not found StartDate or EndDate")
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not found StartDate or EndDate")
	}

	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeScheduleList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ScheduleList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ScheduleList{ScheduleList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) GetTagList(ctx context.Context, e *emptypb.Empty) (*calendarpb.TagList, error) {
	log.Printf("[info] get tag list\n")
	res, err := c.Db.GetTagList()
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.TagList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeTagList(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.TagList{StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.TagList{TagList: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) UpdateTag(ctx context.Context, tag *calendarpb.Tag) (*calendarpb.ResponseTag, error) {
	log.Printf("[info] update tag %+v\n", tag)
	if tag == nil {
		log.Println("[error] not exist Tag")
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("not exist Tag")
	}

	if tag.TagName == "" {
		log.Println("[error] tag name is empty")
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, errors.New("tag name is empty")
	}
	res, err := c.Db.UpdateTag(&TagMaster{TagId: tag.TagId, TagName: tag.TagName})
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}

	result, err := c.makeTag(res)
	if err != nil {
		log.Printf("[error] %+v\n", err)
		return &calendarpb.ResponseTag{Tag: tag, StatusCode: calendarpb.ResponseStatusCode_Failed}, err
	}
	return &calendarpb.ResponseTag{Tag: result, StatusCode: calendarpb.ResponseStatusCode_Success}, nil
}

func (c *CalendarService) makeSchedule(schedule *Schedule) (*calendarpb.Schedule, error) {

	start, err := ptypes.TimestampProto(schedule.StartDate)
	if err != nil {
		return nil, err
	}

	end, err := ptypes.TimestampProto(schedule.EndDate)
	if err != nil {
		return nil, err
	}

	tagList, err := c.Db.SearchTagListByScheduleId(schedule.ScheduleId)
	if err != nil {
		return nil, err
	}

	scheduleTags, err := c.makeTagList(tagList)
	if err != nil {
		return nil, err
	}

	return &calendarpb.Schedule{
		ScheduleId: schedule.ScheduleId,
		Date: &calendarpb.Date{
			Start: start,
			End:   end,
		},
		Title:       schedule.Title,
		Description: schedule.Description,
		User: &calendarpb.User{
			UserId:   schedule.UserId,
			UserName: schedule.UserName,
		},
		TagList: scheduleTags,
	}, nil
}

func (c *CalendarService) makeScheduleList(schedules []Schedule) ([]*calendarpb.Schedule, error) {
	result := make([]*calendarpb.Schedule, len(schedules))
	var err error
	for i, v := range schedules {
		if result[i], err = c.makeSchedule(&v); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *CalendarService) makeTag(tag *TagMaster) (*calendarpb.Tag, error) {
	return &calendarpb.Tag{
		TagId:   tag.TagId,
		TagName: tag.TagName,
	}, nil
}

func (c *CalendarService) makeTagList(tags []TagMaster) ([]*calendarpb.Tag, error) {
	result := make([]*calendarpb.Tag, len(tags))
	var err error
	for i, v := range tags {
		if result[i], err = c.makeTag(&v); err != nil {
			return nil, err
		}
	}

	return result, nil
}
