package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	gormDB *gorm.DB
}

type Schedule struct {
	ScheduleId  uint32 `gorm:"primary_key"`
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	UserId      uint32
	UserName    string
}

func (Schedule) TableName() string {
	return "schedule"
}

type TagMaster struct {
	TagId   uint32 `gorm:"primary_key"`
	TagName string
}

func (TagMaster) TableName() string {
	return "tag_master"
}

// TODO need rename
type ScheduleIdToTagId struct {
	ScheduleId uint32
	TagId      uint32
}

func (ScheduleIdToTagId) TableName() string {
	return "schedule_to_tag"
}

func NewDB(env *DatabaseEnv) (*DB, error) {
	connect := env.User + ":" + env.Password + "@tcp(" + env.Addr + ":" + env.Port + ")/" + env.Name + "?charset=utf8&collation=utf8_general_ci&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(connect), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can not connect mysql user=%s, server=%s:%s, database=%s\n", env.User, env.Addr, env.Port, env.Name)
	}
	return &DB{gormDB}, nil
}

func (db *DB) CreateSchedule(schedule *Schedule, tagList []TagMaster) (*Schedule, error) {
	tx := db.gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, fmt.Errorf("can not begin transactions\n")
	}

	err := tx.Create(&schedule).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create schedule: %+v\n", schedule)
	}

	if tagList != nil {
		for _, tag := range tagList {
			err = db.mapScheduleTag(tx, schedule.ScheduleId, tag.TagId)
			if err != nil {
				return nil, fmt.Errorf("Failed to create mapping schedule=%+v to tag=%+v\n", schedule, tag)
			}
		}
	}

	tx.Commit()
	return schedule, nil
}

func (db *DB) CreateTag(tag *TagMaster) (*TagMaster, error) {
	exists, err := db.findTagIdByTagName(tag.TagName)
	if err != nil {
		return nil, err
	}

	if exists != nil {
		return nil, fmt.Errorf("tag=%s is already exists: %+v", tag.TagName, exists)
	}
	tx := db.gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, fmt.Errorf("can not begin transactions\n")
	}

	err = tx.Create(&tag).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create tag: %+v\n", tag)
	}

	tx.Commit()
	return tag, nil
}

func (db *DB) UpdateTag(tag *TagMaster) (*TagMaster, error) {
	tx := db.gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, fmt.Errorf("can not begin transactions\n")
	}

	err := tx.Save(&tag).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update tag: %+v\n", tag)
	}

	tx.Commit()

	return tag, nil
}

func (db *DB) GetScheduleList() ([]Schedule, error) {
	scheduleList := []Schedule{}

	err := db.gormDB.Find(&scheduleList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all schedule\n")
	}
	return scheduleList, nil
}

func (db *DB) GetTagList() ([]TagMaster, error) {
	tagList := []TagMaster{}

	err := db.gormDB.Find(&tagList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all tag\n")
	}
	return tagList, nil
}

func (db *DB) SearchTagListByScheduleId(scheduleId uint32) ([]TagMaster, error) {
	scheduleTags := []ScheduleIdToTagId{}

	err := db.gormDB.Where("schedule_id = ?", scheduleId).Find(&scheduleTags).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tag mapped to the schedule_id=%d\n", scheduleId)
	}

	tagList := make([]TagMaster, len(scheduleTags))
	for i, v := range scheduleTags {
		tag := TagMaster{}
		err = db.gormDB.Where("tag_id = ?", v.TagId).First(&tag).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get tag_id=%d\n", v.TagId)
		}
		tagList[i] = tag
	}
	return tagList, nil
}

func (db *DB) DeleteSchedule(schedule *Schedule) error {
	tx := db.gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return fmt.Errorf("can not begin transactions\n")
	}

	err := tx.Delete(schedule).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete schedule=%+v\n", schedule)
	}
	err = db.unmapTagbyScheduleId(tx, schedule.ScheduleId)
	if err != nil {
		return fmt.Errorf("failed to unmap tag by schedule=%+v\n", schedule)
	}
	tx.Commit()
	return nil
}

func (db *DB) UpdateSchedule(schedule *Schedule, tagList []TagMaster) (*Schedule, error) {
	tx := db.gormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, fmt.Errorf("can not begin transactions\n")
	}

	err := tx.Save(schedule).Error
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Failed to Update schedule=%+v\n", schedule)
	}

	err = db.unmapTagbyScheduleId(tx, schedule.ScheduleId)
	if err != nil {
		return nil, fmt.Errorf("failed to unmap tag by schedule=%+v\n", schedule)
	}
	if tagList != nil {
		for _, tag := range tagList {
			err = db.mapScheduleTag(tx, schedule.ScheduleId, tag.TagId)
			if err != nil {
				return nil, fmt.Errorf("Failed to create mapping schedule=%+v to tag=%+v\n", schedule, tag)
			}
		}
	}
	tx.Commit()
	return schedule, nil
}

func (db *DB) SearchScheduleByUserId(userId uint32) ([]Schedule, error) {
	record := []Schedule{}
	err := db.gormDB.Where("user_id  = ?", userId).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by user_id=%d\n", userId)
	}
	return record, nil
}

func (db *DB) SearchScheduleByUserName(userName string) ([]Schedule, error) {
	record := []Schedule{}
	err := db.gormDB.Where("user_name  = ?", userName).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by user_name=%s\n", userName)
	}
	return record, nil
}

func (db *DB) SearchScheduleByTagName(tagName string) ([]Schedule, error) {
	tagRecord, err := db.findTagIdByTagName(tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag by tag_name=%s\n", tagName)
	}

	scheduleIdList := []ScheduleIdToTagId{}
	err = db.gormDB.Where("tag_id = ?", tagRecord.TagId).Find(&scheduleIdList).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule_id by tag=%+v\n", tagRecord)
	}
	record := make([]Schedule, len(scheduleIdList))
	for i, v := range scheduleIdList {
		var schedule Schedule
		err = db.gormDB.Where("schedule_id = ?", v.ScheduleId).First(&schedule).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get schedule by schedule_id=%d\n", v.ScheduleId)
		}
		record[i] = schedule
	}
	return record, nil
}

func (db *DB) SearchScheduleDurationByDate(start *time.Time, end *time.Time) ([]Schedule, error) {
	record := []Schedule{}
	err := db.gormDB.Where("start_date BETWEEN ? AND ? OR end_date BETWEEN ? AND ?", start, end, start, end).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by Date start=%v, end=%v\n", start, end)
	}

	return record, nil
}

func (db *DB) SearchScheduleSinceByDate(start *time.Time) ([]Schedule, error) {
	record := []Schedule{}
	err := db.gormDB.Where("start_date >= ? OR end_date >= ?", start, start).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by Date start=%v\n", start)
	}
	return record, nil
}

func (db *DB) SearchScheduleUntileByDate(end *time.Time) ([]Schedule, error) {
	record := []Schedule{}
	err := db.gormDB.Where("start_date <= ? OR end_date <= ?", end, end).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by Date end=%v\n", end)
	}
	return record, nil
}

func (db *DB) findTagIdByTagName(tagName string) (*TagMaster, error) {
	record := []TagMaster{}
	err := db.gormDB.Where("tag_name  = ?", tagName).Find(&record).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find tag: %s", tagName)
	}

	if len(record) > 1 {
		return nil, fmt.Errorf("duplicate tags exist: %s", tagName)
	}

	if len(record) == 0 {
		return nil, nil
	}
	return &record[0], nil
}

func (db *DB) unmapTagbyScheduleId(tx *gorm.DB, scheduleId uint32) error {
	err := tx.Where("schedule_id = ?", scheduleId).Delete(&ScheduleIdToTagId{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (db *DB) mapScheduleTag(tx *gorm.DB, scheduleId uint32, tagId uint32) error {
	err := tx.Create(&ScheduleIdToTagId{ScheduleId: scheduleId, TagId: tagId}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
