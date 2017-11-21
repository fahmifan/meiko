package information

import (
	"fmt"
	"html"
	"strconv"
	"strings"

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
		FilesID:     params.FilesID,
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
	var filesID []string
	// FilesID validation
	if !helper.IsEmpty(params.FilesID) {
		filesID = strings.Split(params.FilesID, "~")
		for _, value := range filesID {
			if !helper.IsValidFileID(value) {
				return args, fmt.Errorf("Wrong Files ID")
			}
		}
	}
	return createArgs{
		Title:       params.Title,
		Description: params.Description,
		ScheduleID:  scheduleID,
		FilesID:     filesID,
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
		FilesID:     params.FilesID,
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
	var filesID []string
	// FilesID validation
	if !helper.IsEmpty(params.FilesID) {
		filesID = strings.Split(params.FilesID, "~")
		for _, value := range filesID {
			if !helper.IsValidFileID(value) {
				return args, fmt.Errorf("Wrong Files ID")
			}
		}
	}
	return upadateArgs{
		ID:          informationID,
		Title:       params.Title,
		Description: params.Description,
		ScheduleID:  scheduleID,
		FilesID:     filesID,
	}, nil

}

func (params deleteParams) validate() (deleteArgs, error) {

	var args deleteArgs
	// Information ID validation
	if helper.IsEmpty(params.ID) {
		return args, fmt.Errorf("Information ID can not be empty")
	}
	informationID, err := strconv.ParseInt(params.ID, 10, 64)
	if err != nil {
		return args, fmt.Errorf("Error convert information id to int64")
	}
	return deleteArgs{
		ID: informationID,
	}, nil
}

func (params readListParams) validate() (readListArgs, error) {
	var args readListArgs
	if helper.IsEmpty(params.Page) || helper.IsEmpty(params.Total) {
		return args, fmt.Errorf("page or total is empty")
	}

	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return args, fmt.Errorf("page must be numeric")
	}

	total, err := strconv.ParseInt(params.Total, 10, 64)
	if err != nil {
		return args, fmt.Errorf("total must be numeric")
	}

	// should be positive number
	if page < 0 || total < 0 {
		return args, fmt.Errorf("page or total must be positive number")
	}

	args = readListArgs{
		Page:  uint16(page),
		Total: uint16(total),
	}
	return args, nil
}