package server

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

type request struct {
	Entity       crux.Markdone_t `json:"entity"`
	Step         string          `json:"step"`
	WorkflowName string          `json:"workflowName"`
}

func MarkDone(c *gin.Context, s *service.Service) {
	l := s.LogHarbour
	l.Debug0().Log("starting execution of CapGet()")
	var req request

	err := wscutils.BindJSON(c, &req)
	if err != nil {
		l.Error(err).Log("Error Unmarshalling Query parameters to struct:")
		return
	}

	query, _ := s.Dependencies["queries"].(*sqlc.Queries)
	ResponseData, err := crux.DoMarkDone(query, req.Entity, req.Step, req.WorkflowName)
	if err != nil {
		l.Debug1().LogDebug("Error while marshaling patternSchema", err)
		wscutils.SendErrorResponse(c, &wscutils.Response{Data: err})
		return
	}
	wscutils.SendSuccessResponse(c, &wscutils.Response{Status: wscutils.SuccessStatus, Data: ResponseData, Messages: nil})

}
