package hafenhause

import (
	"fmt"
	"log"
	"time"

	"github.com/chodyo/hafenhause/nosqldb"

	"github.com/mitchellh/mapstructure"
)

type bedtimedb struct {
	nosqldb.NosqldbContract
}

const collectionName = "hafenhause"

func newBedtimedb() bedtimedb {
	db := nosqldb.NewNosqldb(collectionName)
	return bedtimedb{db}
}

func (db bedtimedb) createDefaultBedtime(name string) (err error) {
	var bedtimeDefaults []interface{}

	defaultsKey := []string{"defaults.bedtime"}
	if bedtimeDefaults, err = db.Read("static", defaultsKey); err != nil || len(bedtimeDefaults) != 1 {
		log.Printf("Failed to get defaults with err: %v\n", err)
		return
	}

	var defaults bedtime
	if err = mapstructure.Decode(bedtimeDefaults[0], &defaults); err != nil {
		log.Printf("Failed to decode defaults with err: %v\n", err)
		return
	}

	now := time.Now()
	bedtime := bedtime{
		Hour:    defaults.Hour,
		Minute:  defaults.Minute,
		Updated: &now,
	}

	newFieldKey := fmt.Sprintf("%s.bedtime", name)
	if err = db.Create("state", newFieldKey, bedtime); err != nil {
		log.Printf("Failed to save new %s to db with err: %v\n", name, err)
		return
	}

	return
}

func (db bedtimedb) getBedtimes(name string) (bedtimes []bedtime, err error) {
	fieldKey := fmt.Sprintf("%s.bedtime", name)

	var fields []interface{}
	if fields, err = db.Read("state", []string{fieldKey}); err != nil {
		log.Printf("Failed to get %s from db with err: %v\n", name, err)
		return
	}

	for _, field := range fields {
		var bedtime bedtime
		if err = mapstructure.Decode(field, &bedtime); err != nil {
			log.Printf("Failed to decode bedtime with err: %v\n", err)
			return
		}

		bedtimes = append(bedtimes, bedtime)
	}

	return
}

func (db bedtimedb) updateBedtime(name string, bedtime bedtime) (err error) {
	now := time.Now()
	bedtime.Updated = &now

	fieldKey := fmt.Sprintf("%s.bedtime", name)
	if err = db.Update("state", fieldKey, bedtime); err != nil {
		log.Printf("Failed to update %s to db with err: %v\n", name, err)
		return
	}

	return
}

func (db bedtimedb) deleteBedtime(name string) (err error) {
	fieldKey := fmt.Sprintf("%s.bedtime", name)
	if err = db.Delete("state", fieldKey); err != nil {
		log.Printf("Failed to delete %s from db with err: %v\n", name, err)
		return
	}

	return
}
