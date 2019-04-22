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

const (
	cRoot = "hafenhause"

	// doc Person name; or
	dDefaults = "Defaults"

	bedtimeField = "bedtime"
)

// ":root/:name/:function" e.g.
// "hafenhause/Cody/bedtime" or
// "hafenhause/Defaults/bedtime"
func getDocPath(docName string) string {
	return fmt.Sprintf("%s/%s", cRoot, docName)
}

func newBedtimedb() bedtimedb {
	db := nosqldb.NewNosqldb()
	return bedtimedb{db}
}

func (db bedtimedb) createDefaultBedtime(name string) (err error) {
	defaultsPath := getDocPath(dDefaults)

	var defaultsContents map[string]interface{}
	if defaultsContents, err = db.Read(defaultsPath); err != nil || len(defaultsContents) != 1 {
		log.Printf("Failed to get defaults with err: %v\n", err)
		return
	}

	var defaultBedtime bedtime
	if err = mapstructure.Decode(defaultsContents[bedtimeField], &defaultBedtime); err != nil {
		log.Printf("Failed to decode defaults with err: %v\n", err)
		return
	}

	now := time.Now()
	toCreate := map[string]bedtime{
		bedtimeField: bedtime{
			Hour:    defaultBedtime.Hour,
			Minute:  defaultBedtime.Minute,
			Updated: &now,
		}}

	newDocPath := getDocPath(name)
	if err = db.Create(newDocPath, toCreate); err != nil {
		log.Printf("Failed to save new %s to db with err: %v\n", name, err)
		return
	}

	return
}

func (db bedtimedb) getBedtimes(name string) (bedtimes []bedtime, err error) {
	if name == "" {
		// TODO: get all
	}

	docPath := getDocPath(name)

	var docContents map[string]interface{}
	if docContents, err = db.Read(docPath); err != nil {
		log.Printf("Failed to get %s from db with err: %v\n", name, err)
		return
	}

	var bedtime bedtime
	if err = mapstructure.Decode(docContents[bedtimeField], &bedtime); err != nil {
		log.Printf("Failed to decode bedtime with err: %v\n", err)
		return
	}

	bedtime.Name = &name
	bedtime.Updated = nil

	bedtimes = append(bedtimes, bedtime)

	return
}

func (db bedtimedb) updateBedtime(name string, b bedtime) (err error) {
	docPath := getDocPath(name)

	now := time.Now()
	b.Updated = &now

	toUpdate := map[string]bedtime{
		bedtimeField: b,
	}

	if err = db.Update(docPath, toUpdate); err != nil {
		log.Printf("Failed to update %s to db with err: %v\n", name, err)
		return
	}

	return
}

func (db bedtimedb) deleteBedtime(name string) (err error) {
	docPath := getDocPath(name)

	if err = db.Delete(docPath, bedtimeField); err != nil {
		log.Printf("Failed to delete %s from db with err: %v\n", name, err)
		return
	}

	return
}
