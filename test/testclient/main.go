package main

import (
	"context"
	"log"
	"time"

	"bitbucket.org/latonaio/calendar-module-kube/cmd/calendar-server/proto/calendarpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerEnv struct {
	Host string
	Port string
}

func printSchedule(schedule *calendarpb.Schedule) {
	log.Printf(
		"schedule_id: %d\nstart_at: %s\n end_at: %s\n title: %s\n description: %s\n user_id: %d\n user_name: %s\n",
		schedule.ScheduleId,
		ptypes.TimestampString(schedule.Date.Start),
		ptypes.TimestampString(schedule.Date.End),
		schedule.Title,
		schedule.Description,
		schedule.User.UserId,
		schedule.User.UserName)
	for _, v := range schedule.TagList {
		printTag(v)
	}
}

func printTag(tag *calendarpb.Tag) {
	log.Printf(
		"tag_id: %d\ntag_name: %s\n",
		tag.TagId,
		tag.TagName)
}

func generateSchedule(title string, description string, start *timestamp.Timestamp, end *timestamp.Timestamp, user *calendarpb.User, tagList []*calendarpb.Tag) *calendarpb.Schedule {
	schedule := &calendarpb.Schedule{
		Date: &calendarpb.Date{
			Start: start,
			End:   end,
		},
		Title:       title,
		User:        user,
		Description: description,
		TagList:     tagList,
	}
	return schedule
}

func main() {
	var serverEnv ServerEnv
	envconfig.Process("SERVER", &serverEnv)
	conn, err := grpc.Dial(serverEnv.Host+":"+serverEnv.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close()
	client := calendarpb.NewCalendarClient(conn)

	const location = "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}

	// 2020/01/01 10:00
	baseDate := time.Date(2020, 1, 1, 10, 0, 0, 0, loc)
	date20200101, _ := ptypes.TimestampProto(baseDate)
	date20200102, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 1))
	date20200103, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 2))
	date20200104, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 3))
	date20200105, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 4))
	date20200106, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 5))
	date20200107, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 6))
	date20200108, _ := ptypes.TimestampProto(baseDate.AddDate(0, 0, 7))

	userAlice := &calendarpb.User{
		UserId:   1,
		UserName: "Alice",
	}
	userBob := &calendarpb.User{
		UserId:   2,
		UserName: "Bob",
	}

	// Test create tag
	log.Println("create tag")
	tagCleaning, err := client.CreateTag(context.TODO(), &calendarpb.Tag{TagName: "cleaning"})
	if err != nil {
		log.Println(err)
	} else {
		printTag(tagCleaning.Tag)
	}

	log.Println("create tag")
	tagCheckedIn, err := client.CreateTag(context.TODO(), &calendarpb.Tag{TagName: "checked in"})
	if err != nil {
		log.Println(err)
	} else {
		printTag(tagCheckedIn.Tag)
	}

	// Test update tag
	log.Println("Update tag")
	tagCheckedOut, err := client.UpdateTag(context.TODO(), &calendarpb.Tag{TagId: tagCheckedIn.Tag.TagId, TagName: "checked out"})
	if err != nil {
		log.Println(err)
	} else {
		printTag(tagCheckedOut.Tag)
	}

	// Test get tag list
	log.Println("Get tag list")
	tagListRes, err := client.GetTagList(context.TODO(), &emptypb.Empty{})
	if err != nil {
		log.Println(err)
	} else {
		for _, v := range tagListRes.TagList {
			printTag(v)
		}
	}

	// Test create schedule
	log.Println("create schedule: ")
	res, err := client.CreateSchedule(context.TODO(), generateSchedule(
		"test title",
		"description",
		date20200101,
		date20200102,
		userAlice,
		tagListRes.TagList))
	if err != nil {
		log.Println(err)
	} else {
		printSchedule(res.Schedule)
	}

	log.Println("create schedule: ")
	res, err = client.CreateSchedule(context.TODO(), generateSchedule(
		"test title",
		"description",
		date20200103,
		date20200104,
		userAlice,
		nil))
	if err != nil {
		log.Println(err)
	} else {
		printSchedule(res.Schedule)
	}

	log.Println("create schedule: ")
	res, err = client.CreateSchedule(context.TODO(), generateSchedule(
		"test title",
		"description",
		date20200105,
		date20200107,
		userBob,
		nil))
	if err != nil {
		log.Println(err)
	} else {
		printSchedule(res.Schedule)
	}

	// Test esarch schedule by user id
	resultList, err := client.SearchScheduleByUserId(context.TODO(), userAlice)
	if err != nil {
		log.Println(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}

	// Test search schedule by user name
	log.Println("Search ScheduleByUserName: Bob")
	resultList, err = client.SearchScheduleByUserName(context.TODO(), userBob)
	if err != nil {
		log.Println(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}
	// Test search schedule by tag name
	log.Println("Search ScheduleByTagName: cleaning")
	resultList, err = client.SearchScheduleByTagName(context.TODO(), &calendarpb.Tag{TagName: "cleaning"})
	if err != nil {
		log.Fatalln(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}

	// Test search schedule by date
	log.Println("Search ScheduleByDate: since 2020/01/04")
	resultList, err = client.SearchScheduleByDate(context.TODO(), &calendarpb.Date{Start: date20200104})
	if err != nil {
		log.Fatalln(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}

	log.Println("Search ScheduleByDate: untile 2020/01/04")
	resultList, err = client.SearchScheduleByDate(context.TODO(), &calendarpb.Date{End: date20200104})
	if err != nil {
		log.Fatalln(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}

	log.Println("Search ScheduleByDate: while 2020/01/04 to 2020/01/07")
	resultList, err = client.SearchScheduleByDate(context.TODO(), &calendarpb.Date{Start: date20200104, End: date20200107})
	if err != nil {
		log.Fatalln(err)
	} else {
		for _, v := range resultList.ScheduleList {
			printSchedule(v)
		}
	}

	// Test update schedule
	log.Println("create schedule: ")
	tags := make([]*calendarpb.Tag, 1)
	tags[0] = tagCheckedOut.Tag
	res, err = client.CreateSchedule(context.TODO(), generateSchedule(
		"test title",
		"description",
		date20200106,
		date20200107,
		userAlice,
		tags))
	if err != nil {
		log.Println(err)
	} else {
		printSchedule(res.Schedule)
	}

	res.Schedule.TagList[0] = tagCleaning.Tag
	log.Println("update schedule: ")
	res.Schedule.Date.Start = date20200107
	res.Schedule.Date.End = date20200108
	res, err = client.UpdateSchedule(context.TODO(), res.Schedule)
	if err != nil {
		log.Println(err)
	} else {
		printSchedule(res.Schedule)
	}

	// Test delete schedule
	log.Println("delete schedule: ")
	_, err = client.DeleteSchedule(context.TODO(), res.Schedule)
	if err != nil {
		log.Println(err)
	}
}
