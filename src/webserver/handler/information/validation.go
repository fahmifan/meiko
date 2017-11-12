package information

import (
	"fmt"
	"html"
	"strconv"

	"github.com/melodiez14/meiko/src/util/helper"
)

func (params detailInfromationParams) validate() (detailInfromationArgs, error) {
	var args detailInfromationArgs
	// Information ID validation
	if helper.IsEmpty(params.ID) {
		return args, fmt.Errorf("Information ID can not be empty")
	}
	id, err := strconv.ParseInt(params.ID, 10, 64)
	if err != nil {
		return args, fmt.Errorf("Error to convert Information string ID to int64")
	}
	return detailInfromationArgs{
		ID: id,
	}, nil
}
func (params createParams) validate() (createArgs, error) {

	var args createArgs
	params = createParams{
		Title:       html.EscapeString(params.Title),
		Description: html.EscapeString(params.Description),
		ScheduleID:  params.ScheduleID,
	}

	// Title validation
	if helper.IsEmpty(params.Title) {
		return args, fmt.Errorf("Title can not be empty")
	}

	// Description validation
	if helper.IsEmpty(params.Description) {
		return args, fmt.Errorf("Content can not be empty")
	}

	// Schedule ID validation
	var scheduleID int64
	var err error
	if !helper.IsEmpty(params.ScheduleID) {
		scheduleID, err = strconv.ParseInt(params.ScheduleID, 10, 64)
		if err != nil {
			return args, err
		}
	}

	return createArgs{
		Title:       params.Title,
		Description: params.Description,
		ScheduleID:  scheduleID,
	}, nil

}
func (params updateParams) validate() (upadateArgs, error) {

	var args upadateArgs
	var err error
	params = updateParams{
		ID:          params.ID,
		Title:       html.EscapeString(params.Title),
		Description: html.EscapeString(params.Description),
		ScheduleID:  params.ScheduleID,
	}
	// Information ID validation
	if helper.IsEmpty(params.ID) {
		return args, fmt.Errorf("Information ID can not be empty")
	}
	informationID, err := strconv.ParseInt(params.ID, 10, 64)
	if err != nil {
		return args, fmt.Errorf("Error convert information id to int64")
	}

	// Title validation
	if helper.IsEmpty(params.Title) {
		return args, fmt.Errorf("Title can not be empty")
	}

	// Description validation
	if helper.IsEmpty(params.Description) {
		return args, fmt.Errorf("Content can not be empty")
	}

	// Schedule ID validation
	var scheduleID int64
	if !helper.IsEmpty(params.ScheduleID) {
		scheduleID, err = strconv.ParseInt(params.ScheduleID, 10, 64)
		if err != nil {
			return args, err
		}
	}

	return upadateArgs{
		ID:          informationID,
		Title:       params.Title,
		Description: params.Description,
		ScheduleID:  scheduleID,
	}, nil

}