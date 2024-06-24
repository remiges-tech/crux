package wfinstance

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server"
	"github.com/remiges-tech/crux/types"
)

// WFInstanceList rquest format
type WFInstanceListRequest struct {
	Slice    *int32  `json:"slice" validate:"omitempty,gt=0,lt=50"`
	EntityID *string `json:"entityid" validate:"omitempty,gt=0,lt=40"`
	App      *string `json:"app" validate:"omitempty,alpha,gt=0,lt=15"`
	Workflow *string `json:"workflow" validate:"omitempty,gt=0,lt=20"`
	Parent   *int32  `json:"parent" validate:"omitempty,gt=0"`
}

// WFInstanceList response format
type WFInstanceListResponse struct {
	ID         int32            `json:"id"`
	EntityID   string           `json:"entityid"`
	Slice      int32            `json:"slice"`
	App        string           `json:"app"`
	Class      string           `json:"class"`
	Workflow   string           `json:"workflow"`
	Step       string           `json:"step"`
	LoggedDate pgtype.Timestamp `json:"loggedat"`
	DoneAt     string           `json:"doneat"`
	Nextstep   string           `json:"nextstep"`
	Parent     int32            `json:"parent,omitempty"`
}

type WFInstanceListParams struct {
	Slice    int32
	EntityID string
	App      string
	Workflow string
	Parent   int32
}

func GetWFInstanceList(c *gin.Context, s *service.Service) {
	lh := s.LogHarbour.WithClass("wfinstance")
	lh.Log("GetWFInstanceList request received")

	var (
		request WFInstanceListRequest
		params  WFInstanceListParams
	)

	// userID, err := server.ExtractUserNameFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract userID from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }

	// realm, err := server.ExtractRealmFromJwt(c)
	// if err != nil {
	// 	lh.Info().Log("unable to extract realm from token")
	// 	wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Missing, server.ErrCode_Token_Data_Missing))
	// 	return
	// }
	isCapable, _ := server.Authz_check(types.OpReq{
		User: userID,
	}, false)

	if !isCapable {
		lh.Info().LogActivity("GetWFInstanceList||unauthorized user:", userID)
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_Unauthorized, server.ErrCode_Unauthorized))
		return
	}

	// Bind request
	err := wscutils.BindJSON(c, &request)
	if err != nil {
		lh.Error(err).Log("GetWFInstanceList||error while binding json request error")
		return
	}

	// Check for validation error
	valError := wscutils.WscValidate(request, func(err validator.FieldError) []string { return []string{} })
	if len(valError) > 0 {
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, valError))
		lh.LogActivity("GetWFInstanceList||validation error:", valError)
		return
	}

	// To get request parameters
	params = GetParams(request, params)
	lh.Debug0().LogActivity("GetWFInstanceList|parameters from request :", params)

	query, ok := s.Dependencies["queries"].(*sqlc.Queries)
	if !ok {
		lh.Log("GetWFInstanceList||error while getting query instance from service dependencies")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_InternalErr, server.ErrCode_DatabaseError))
		return
	}

	// To get requested wfinstance list
	wfinstanceList, error := query.GetWFInstanceList(c, sqlc.GetWFInstanceListParams{
		Slice:    pgtype.Int4{Int32: params.Slice, Valid: params.Slice != 0},
		Entityid: pgtype.Text{String: params.EntityID, Valid: !server.IsStringEmpty(&params.EntityID)},
		App:      pgtype.Text{String: params.App, Valid: !server.IsStringEmpty(&params.App)},
		Workflow: pgtype.Text{String: params.Workflow, Valid: !server.IsStringEmpty(&params.Workflow)},
		Parent:   pgtype.Int4{Int32: params.Parent, Valid: params.Parent != 0},
	})
	if error != nil {
		lh.Error(error).Log("GetWFInstanceList||error while getting wfinstance List")
		errmsg := db.HandleDatabaseError(err)
		wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
		return
	}

	if len(wfinstanceList) == 0 {
		lh.Debug0().Log("GetWFInstanceList||requestd wfinstanceList not found")
		wscutils.SendErrorResponse(c, wscutils.NewErrorResponse(server.MsgId_NotFound, server.ErrCode_NotFound))
		return
	}

	// To Get parentList from WFInstanceList data
	parentList := getParentList(wfinstanceList)
	lh.Debug0().LogActivity("GetWFInstanceList||parentlist form WFInstanceList data:", parentList)

	for parentList != nil {
		lh.Debug0().Log("GetWFInstanceList||inside for loop : if parentList is Not Nil ")
		// To get GetWFInstanceList by parentList
		wfinstanceListByParents, err := query.GetWFInstanceListByParents(c, parentList)
		if err != nil {
			lh.Error(err).Log("GetWFInstanceList||error while getting wfinstance List by parentList")
			errmsg := db.HandleDatabaseError(err)
			wscutils.SendErrorResponse(c, wscutils.NewResponse(wscutils.ErrorStatus, nil, []wscutils.ErrorMessage{errmsg}))
			return
		}

		// Append wfinstanceListByParents data
		wfinstanceList = append(wfinstanceList, wfinstanceListByParents...)

		// Update parentList using getParentList function
		parentList = getParentList(wfinstanceListByParents)
		lh.Debug0().LogActivity("GetWFInstanceList||updated ParentList :", parentList)
	}
	var responseList []WFInstanceListResponse
	for _, val := range wfinstanceList {
		var response WFInstanceListResponse
		response.ID = val.ID
		response.EntityID = val.Entityid
		response.Slice = val.Slice
		response.App = val.App
		response.Class = val.Class
		response.Workflow = val.Workflow
		response.Step = val.Step
		response.LoggedDate = val.Loggedat
		response.Nextstep = val.Nextstep

		// Handling the omitted Parent field
		if val.Parent.Valid {
			response.Parent = val.Parent.Int32
		}

		// Handling the DoneAt field
		if !val.Doneat.Valid {
			response.DoneAt = ""
		} else {
			response.DoneAt = val.Doneat.Time.String()
		}

		responseList = append(responseList, response)
	}

	lh.Debug0().Log("Record found finished execution of GetWFInstanceList request")
	wscutils.SendSuccessResponse(c, wscutils.NewSuccessResponse(map[string][]WFInstanceListResponse{"wfinstance": responseList}))
}

// To Get Parameters from WFInstanceList Request
func GetParams(req WFInstanceListRequest, params WFInstanceListParams) WFInstanceListParams {

	if req.Slice != nil {
		params.Slice = *req.Slice
	}
	if req.EntityID != nil {
		params.EntityID = *req.EntityID
	}
	if req.App != nil {
		params.App = *req.App
	}
	if req.Workflow != nil {
		params.Workflow = *req.Workflow
	}
	if req.Parent != nil {
		params.Parent = *req.Parent
	}
	return params
}
