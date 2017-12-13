package information

import (
	"net/http"
	"strconv"
	"time"

	"github.com/melodiez14/meiko/src/util/conn"

	"github.com/melodiez14/meiko/src/util/helper"
	"github.com/melodiez14/meiko/src/webserver/template"

	"github.com/julienschmidt/httprouter"
	"github.com/melodiez14/meiko/src/module/course"
	cs "github.com/melodiez14/meiko/src/module/course"
	fs "github.com/melodiez14/meiko/src/module/file"
	inf "github.com/melodiez14/meiko/src/module/information"
	rg "github.com/melodiez14/meiko/src/module/rolegroup"
	"github.com/melodiez14/meiko/src/util/alias"
	"github.com/melodiez14/meiko/src/util/auth"
)

// GetSummaryHandler func ...
func GetSummaryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess := r.Context().Value("User").(*auth.User)

	// get enrolled course
	schedulesID, err := course.SelectIDByUserID(sess.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError).
			AddError("Internal server error"))
		return
	}

	// get information list
	informations, err := inf.SelectByScheduleID(schedulesID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError).
			AddError("Internal server error"))
		return
	}

	// convert informations to response
	var informationResponses []informationResponse
	t2 := time.Now()
	for _, val := range informations {
		informationResponses = append(informationResponses, informationResponse{
			Title:       val.Title,
			Date:        helper.DateToString(val.CreatedAt, t2),
			Description: val.Description.String,
		})
	}

	// if informations has only 5, so last and recent will be the same
	// else it has 5 last information and other is recent
	var res getSummaryResponse
	if len(informationResponses) <= alias.InformationMinimumLast {
		res = getSummaryResponse{
			Last:   informationResponses,
			Recent: informationResponses,
		}
	} else {
		res = getSummaryResponse{
			Last:   informationResponses[:alias.InformationMinimumLast],
			Recent: informationResponses[alias.InformationMinimumLast:],
		}
	}

	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetData(res))
	return
}

// CreateHandler func ...
func CreateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess := r.Context().Value("User").(*auth.User)
	if !sess.IsHasRoles(rg.ModuleInformation, rg.RoleCreate, rg.RoleXCreate) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusForbidden).
			AddError("You don't have privilege"))
		return
	}
	params := createParams{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		ScheduleID:  r.FormValue("schedule_did"),
		FilesID:     r.FormValue("file_id"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	tx := conn.DB.MustBegin()
	// Insert
	tableID, err := inf.Insert(args.Title, args.Description, args.ScheduleID, tx)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	if args.FilesID != nil {
		for _, fileID := range args.FilesID {
			// need to change
			err := fs.UpdateRelation(fileID, TableNameInformation, tableID, tx)
			if err != nil {
				tx.Rollback()
				template.RenderJSONResponse(w, new(template.Response).
					SetCode(http.StatusBadRequest).
					AddError("Wrong File ID"))
				return
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError))
		return
	}

	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetMessage("Information created successfully"))
	return

}

// UpdateHandler func ...
func UpdateHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess := r.Context().Value("User").(*auth.User)
	if !sess.IsHasRoles(rg.ModuleInformation, rg.RoleUpdate, rg.RoleXUpdate) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusForbidden).
			AddError("You don't have privilege"))
		return
	}
	params := updateParams{
		ID:          ps.ByName("id"),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		ScheduleID:  r.FormValue("schedule_id"),
		FilesID:     r.FormValue("file_id"),
	}

	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	// check is information id exist?
	if !inf.IsInformationIDExist(args.ID) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Information ID does not exist"))
		return
	}
	// check is shedule ID exist
	if args.ScheduleID != 0 {
		if !cs.IsExistScheduleID(args.ScheduleID) {
			template.RenderJSONResponse(w, new(template.Response).
				SetCode(http.StatusBadRequest).
				AddError("Schedule ID does not exist"))
			return
		}
	}
	tx := conn.DB.MustBegin()
	err = inf.Update(args.Title, args.Description, args.ScheduleID, args.ID, tx)
	if err != nil {
		tx.Rollback()
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError))
		return
	}
	// Get All relations with
	filesIDDB, err := fs.GetByStatus(fs.StatusExist, args.ID)
	if err != nil {
		tx.Rollback()
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusInternalServerError))
		return
	}
	var tableID = strconv.FormatInt(args.ID, 10)
	// Add new file
	for _, fileID := range args.FilesID {
		if !fs.IsExistID(fileID) {
			filesIDDB = append(filesIDDB, fileID)
			// Update relation
			// need to change
			err := fs.UpdateRelation(fileID, TableNameInformation, tableID, tx)
			if err != nil {
				tx.Rollback()
				template.RenderJSONResponse(w, new(template.Response).
					SetCode(http.StatusInternalServerError))
				return
			}
		}
	}
	for _, fileIDDB := range filesIDDB {
		isSame := 0
		for _, fileIDUser := range args.FilesID {
			if fileIDUser == fileIDDB {
				isSame = 1
			}
		}
		if isSame == 0 {
			err := fs.UpdateStatusFiles(fileIDDB, fs.StatusDeleted, tx)
			if err != nil {
				tx.Rollback()
				template.RenderJSONResponse(w, new(template.Response).
					SetCode(http.StatusInternalServerError))
				return
			}
		}
	}
	err = tx.Commit()
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetMessage("Update information succesfully"))
	return
}

// GetDetailByAdminHandler func ...
func GetDetailByAdminHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess := r.Context().Value("User").(*auth.User)
	if !sess.IsHasRoles(rg.ModuleInformation, rg.RoleRead, rg.RoleXRead) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusForbidden).
			AddError("You don't have privilege"))
		return
	}
	params := detailInfromationParams{
		ID: ps.ByName("id"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	// check is information id exist?
	if !inf.IsInformationIDExist(args.ID) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Information ID does not exist"))
		return
	}
	res, err := inf.GetByID(args.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	id := res.ScheduleID.Int64
	if id != 0 {
		if !cs.IsEnrolled(sess.ID, id) {
			template.RenderJSONResponse(w, new(template.Response).
				SetCode(http.StatusBadRequest).
				AddError("You does not have permission"))
			return
		}
	}
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusBadRequest).
		SetData(res))
	return

}

// GetListHandler func ...
func GetListHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess := r.Context().Value("User").(*auth.User)
	if !sess.IsHasRoles(rg.ModuleInformation, rg.RoleRead, rg.RoleXRead) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusForbidden).
			AddError("You don't have privilege"))
		return
	}
	params := readListParams{
		Total: r.FormValue("ttl"),
		Page:  r.FormValue("pg"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	scheduleID, err := cs.SelectScheduleIDByUserID(sess.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	offset := (args.Page - 1) * args.Total
	result, err := inf.SelectByPage(scheduleID, args.Total, offset)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetData(result))
	return
}

// DeleteHandler func ...
func DeleteHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess := r.Context().Value("User").(*auth.User)
	if !sess.IsHasRoles(rg.ModuleInformation, rg.RoleDelete, rg.RoleXDelete) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusForbidden).
			AddError("You don't have privilege"))
		return
	}
	params := deleteParams{
		ID: ps.ByName("id"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError(err.Error()))
		return
	}
	// check is information id exist?
	if !inf.IsInformationIDExist(args.ID) {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Information ID does not exist"))
		return
	}
	// delete query
	err = inf.Delete(args.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Delete failed"))
		return
	}
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetMessage("Delete information successfully"))
	return

}

// GetDetailHandler func ..
func GetDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	sess := r.Context().Value("User").(*auth.User)
	params := detailInfromationParams{
		ID: ps.ByName("id"),
	}
	args, err := params.validate()
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).AddError(err.Error()))
		return
	}
	scheduleID := inf.GetScheduleIDByID(args.ID)
	if scheduleID != 0 {
		if !course.IsEnrolled(sess.ID, scheduleID) {
			template.RenderJSONResponse(w, new(template.Response).
				SetCode(http.StatusBadRequest).
				AddError("you do not have permission to this informations"))
			return
		}
	}
	res, err := inf.GetByID(args.ID)
	if err != nil {
		template.RenderJSONResponse(w, new(template.Response).
			SetCode(http.StatusBadRequest).
			AddError("Information does not exist"))
		return
	}
	template.RenderJSONResponse(w, new(template.Response).
		SetCode(http.StatusOK).
		SetData(res))
	return
}
