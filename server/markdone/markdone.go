package markdone

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	crux "github.com/remiges-tech/crux/matching-engine"
	"github.com/remiges-tech/crux/server/wfinstance"

	"github.com/remiges-tech/alya/service"
)

const (
	doneProp = "done"
)

/*
//
// we first check to see if the previous step failed. If yes, we ask the ruleset
// what to do next -- typically expecting an alternate step or remedial step
//
if stepfailed == true then
set entity.step = the step supplied in the request

	    invoke doMatch() with
			entity = the object received
			ruleset = the ruleset name retrieved from wfinstance
			actionset and seenruleset = empty array
	    if doMatch() returns a critical error then
	        return with the critical error details
	    endif
	    if actionset.properties.done == true then
	        delete the wfinstance record
	        return specifying that the workflow is completed
	    endif
	    if actionset.tasks has more than one task then
	        delete the old record from wfinstance and create fresh records, one per task
	        return the full set of tasks and their record IDs
	    else
	        update the old record in wfinstance to set the value of "step" = the task returned
	        return the task and other data in response
	    endif

endif

//
// we come here knowing that the previous step didn't fail. We can now proceed
// to the next step; the previous step was successful
//
recordcount = count_of(SELECT from wfinstance with matching tuple (slice,app,workflow,entityid))
if recordcount == 1 then

	set entity.step = the step supplied in the request
	invoke doMatch() with
	        entity = the object received
	        ruleset = the ruleset name retrieved from wfinstance
	        actionset and seenrulesets: empty

	if doMatch() returns a critical error then
	    return with the critical error details
	endif
	if actionset.properties.done == true then
	    delete the wfinstance record
	    return specifying that the workflow is completed
	endif
	if actionset.tasks has more than one task then
	    delete the old record from wfinstance and create fresh records, one per task
	    return the full set of tasks and their record IDs
	else
	    update the old record in wfinstance to set the value of "step" to the task returned
	    return the task and other data in response
	endif

else (this means count > 1)

	//
	// at this point, we have found multiple records with the same entity ID and
	// workflow, which means they differ only in the value of "step", i.e. the last
	// step done. These records represent multiple asynchronous steps which were being
	// executed in parallel by the application on one entity as part of one workflow.
	//
	// Which of these multiple wfinstances represents my current call to MarkDone()? This can
	// be identified by matching the "step" value in the input with that in the record.
	// Other records refer to other asynchronous steps.
	//
	// When we complete one step out of a set of asynchronous steps, we log this completion
	// and see if there is a next step which we can start right away, or we need to wait for
	// other asynchronous steps to complete.
	//
	set the doneat field in the current wfinstance record to the current timestamp
	look through all the other wfinstance records which have matching tuple (slice,app,workflow,entityid)
	if all of them have doneat set      // this means they are all complete
	    set entity.step = the value of "nextstep" from the current wfinstance record
	    invoke doMatch() with
	            entity = the object received
	            ruleset = the ruleset name retrieved from wfinstance
	            actionset and seenrulesets: empty
	    if doMatch() returns with a critical error then
	        return with error details
	    endif
	    if actionset.properties.done == true then
	        delete all wfinstance records with tuple matching (slice,app,workflow,entityid)
	        return specifying that the workflow is completed
	    endif
	    if actionset.tasks has more than one task then
	        delete the old record from wfinstance and create fresh records, one per task
	        return the full set of tasks and their record IDs
	    else
	        update the old record in wfinstance to set the value of "step" to the task returned
	        return the task and other data in response
	    endif
	else
	    // we come here when our current step is one of a set of concurrent steps
	    // and one or more of the other concurrent steps is yet to complete. In that
	    // we have nothing else to do other than mark the current step complete and
	    // return to the caller saying "We have marked it done, there is nothing more
	    // to do till one more of the other concurrent steps completes. Keep walking."
	    return with details of success of mark-done.
	endif

endif
*/
func DoMarkDone(c *gin.Context, s *service.Service, queries *sqlc.Queries, DoMarkDoneParam Markdone_t) (wfinstance.WFInstanceNewResponse, error) {
	lh := s.LogHarbour.WithClass("wfinstance")
	lh.Log("GetWFinstanceNew request received")

	stepfailed, err := strconv.ParseBool(DoMarkDoneParam.Entity.Attrs["stepfailed"])
	if err != nil {
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while converting stepfailed val from string to bool")
	}
	step := DoMarkDoneParam.Entity.Attrs["step"]

	var cacheData = crux.Cache{
		Ctx:          c,
		Query:        queries,
		Slice:        crux.Slice_t(DoMarkDoneParam.Entity.Slice),
		App:          crux.App_t(DoMarkDoneParam.Entity.App),
		Class:        crux.ClassName_t(DoMarkDoneParam.Entity.Class),
		Realm:        crux.Realm_t(DoMarkDoneParam.Entity.Realm),
		WorkflowName: DoMarkDoneParam.Workflow,
	}

	schema, ruleset, err := cacheData.RetriveRuleSchemasAndRuleSetsFromCache("W")
	if err != nil {
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while Retrieve RuleSchemas and RuleSets FromCache: %v", err)
	} else if schema == nil || ruleset == nil {
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("didn't find any data in RuleSchemas or RuleSets cache: ")
	}

	var response wfinstance.WFInstanceNewResponse

	if stepfailed {

		actionSet := crux.ActionSet{}
		seenRuleSets := make(map[string]struct{})

		// Call the doMatch function passing the entity.entity, ruleset, and the empty actionSet and seenRuleSets
		actionset, _, err, _ := crux.DoMatch(DoMarkDoneParam.Entity, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})
		if err != nil {
			return wfinstance.WFInstanceNewResponse{}, err
		}

		if actionset.Properties[doneProp] == "true" {
			// Delete the wfinstance record

			err := deleteWFInstance(queries, DoMarkDoneParam)
			if err != nil {
				return wfinstance.WFInstanceNewResponse{}, err
			} else {

				response = wfinstance.WFInstanceNewResponse{
					Done: "true",
				}

			}
			return response, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task

			// call addTasks

			DoMarkDoneParam.Entity.Attrs["class"] = DoMarkDoneParam.Entity.Class

			task := wfinstance.AddTaskRequest{
				Steps:    actionset.Tasks,
				Nextstep: actionset.Properties["nextstep"],
				Request: wfinstance.WFInstanceNewRequest{
					Slice:    DoMarkDoneParam.Entity.Slice,
					App:      DoMarkDoneParam.Entity.App,
					EntityID: DoMarkDoneParam.EntityID,
					Entity:   DoMarkDoneParam.Entity.Attrs,
					Workflow: DoMarkDoneParam.Workflow,
				},
			}

			response, err = wfinstance.AddTasks(task, s, c)
			if err != nil {
				lh.Error(err).Log("Error while AddTasks")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return response, nil
		} else {
			DoMarkDoneParam.Entity.Attrs["step"] = actionset.Tasks[0]
			err = UpdateWFInstanceStep(queries, DoMarkDoneParam, actionset.Tasks[0], ruleset.SetName)
			if err != nil {
				lh.Error(err).Log("Error while Update WFInstance Step")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			response = wfinstance.WFInstanceNewResponse{
				Tasks:    []map[string]int32{{step: DoMarkDoneParam.InstanceID}},
				Loggedat: pgtype.Timestamp{Time: DoMarkDoneParam.Loggedat, Valid: true},
			}
			return response, nil
		}
	}
	// We come here k	wing that the previous step didn't fail. We can now proceed to the next step; the previous step was successful
	recordcount, _ := GetWFInstanceCountForEntity(queries, DoMarkDoneParam, ruleset.SetName)
	if recordcount == 1 {
		// markDoneReq.Step = step

		// Invoke doMatch() with
		// entity = the object received
		// ruleset = the ruleset name retrieved from wfinstance
		// actionset and seenrulesets: empty
		actionSet := crux.ActionSet{}
		seenRuleSets := make(map[string]struct{})

		actionset, match, err, _ := crux.DoMatch(DoMarkDoneParam.Entity, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})
		if err != nil {
			return wfinstance.WFInstanceNewResponse{}, err
		}

		if !match {
			return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("rule are not matched")
		}
		if actionset.Properties[doneProp] == "true" {
			// Delete the wfinstance record
			// Return specifying that the workflow is completed
			err := deleteWFInstance(queries, DoMarkDoneParam)
			if err != nil {

				return response, err
			} else {

				response = wfinstance.WFInstanceNewResponse{
					Done: "true",
				}

			}
			return response, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task
			// Return the full set of tasks and their record IDs

			err := deleteWFInstance(queries, DoMarkDoneParam)
			if err != nil {
				return wfinstance.WFInstanceNewResponse{}, err
			}
			doneAtTimeStamp := time.Now()
			err = UpdateWFInstanceDoneAt(queries, DoMarkDoneParam, doneAtTimeStamp, ruleset.SetName)
			if err != nil {
				return wfinstance.WFInstanceNewResponse{}, err
			}

			// call addTasks

			DoMarkDoneParam.Entity.Attrs["class"] = DoMarkDoneParam.Entity.Class

			addTaskParam := wfinstance.AddTaskRequest{
				Steps:    actionset.Tasks,
				Nextstep: actionset.Properties["nextstep"],
				Request: wfinstance.WFInstanceNewRequest{
					Slice:    DoMarkDoneParam.Entity.Slice,
					App:      DoMarkDoneParam.Entity.App,
					EntityID: DoMarkDoneParam.EntityID,
					Entity:   DoMarkDoneParam.Entity.Attrs,
					Workflow: DoMarkDoneParam.Workflow,
				},
			}

			response, err = wfinstance.AddTasks(addTaskParam, s, c)
			if err != nil {
				lh.Error(err).Log("Error while AddTasks")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return response, nil
		} else {
			step = actionset.Tasks[0]
			UpdateWFInstanceStep(queries, DoMarkDoneParam, actionset.Tasks[0], ruleset.SetName)
			response = wfinstance.WFInstanceNewResponse{
				Tasks:    []map[string]int32{{step: DoMarkDoneParam.InstanceID}},
				Loggedat: pgtype.Timestamp{Time: DoMarkDoneParam.Loggedat, Valid: true},
			}
			return response, nil
		}
	} else if recordcount > 1 {
		// At this point, we have found multiple records with the same entity ID and workflow, which means they differ only in the value of "step"
		// Set the doneat field in the current wfinstance record to the current timestamp
		doneAtTimeStamp := time.Now()
		err := UpdateWFInstanceDoneAt(queries, DoMarkDoneParam, doneAtTimeStamp, ruleset.SetName)
		if err != nil {
			lh.Error(err).Log("Error while deleteWFInstance() in DoMarkDone")
			return wfinstance.WFInstanceNewResponse{}, err
		}
		// Look through all the other wfinstance records which have matching tuple (slice,app,workflow,entityid)
		wfInstances, err := getWFInstanceList(queries, DoMarkDoneParam, ruleset.SetName)
		if err != nil {
			return wfinstance.WFInstanceNewResponse{}, err
		}

		// Check if all other wfinstance records have doneat set
		allDone := true
		for _, wfInstance := range wfInstances {
			v, err := wfInstance.Doneat.Value()
			if err != nil {
				return wfinstance.WFInstanceNewResponse{}, err
			}

			if v == nil {
				allDone = false
				break
			}
		}

		if allDone {

			// The value of "nextstep" from the current wfinstance record
			// Invoke doMatch() with
			//  entity = the object received
			// ruleset = the ruleset name retrieved from wfinstance
			// actionset and seenrulesets: empty
			actionSet := crux.ActionSet{}
			seenRuleSets := make(map[string]struct{})

			actionset, _, err, _ := crux.DoMatch(DoMarkDoneParam.Entity, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})
			if err != nil {
				return wfinstance.WFInstanceNewResponse{}, err
			}

			if actionset.Properties[doneProp] == "true" {
				// Delete all wfinstance records with tuple matching (slice, app, workflow, entityid)
				// Return specifying that the workflow is completed

				err := deleteWFInstance(queries, DoMarkDoneParam)
				if err != nil {

					return response, err
				} else {

					response = wfinstance.WFInstanceNewResponse{
						Done: "true",
					}

				}
				return response, nil
			}

			if len(actionset.Tasks) > 1 {
				// Delete the old record from wfinstance and create fresh records, one per task
				// Return the full set of tasks and their record IDs
				err := deleteWFInstance(queries, DoMarkDoneParam)
				if err != nil {
					lh.Error(err).Log("Error while deleteWFInstance() in DoMarkDone")
					return wfinstance.WFInstanceNewResponse{}, err
				}

				// call addTasks

				DoMarkDoneParam.Entity.Attrs["class"] = DoMarkDoneParam.Entity.Class

				task := wfinstance.AddTaskRequest{
					Steps:    actionset.Tasks,
					Nextstep: actionset.Properties["nextstep"],
					Request: wfinstance.WFInstanceNewRequest{
						Slice:    DoMarkDoneParam.Entity.Slice,
						App:      DoMarkDoneParam.Entity.App,
						EntityID: DoMarkDoneParam.EntityID,
						Entity:   DoMarkDoneParam.Entity.Attrs,
						Workflow: DoMarkDoneParam.Workflow,
					},
				}

				response, err = wfinstance.AddTasks(task, s, c)
				if err != nil {
					lh.Error(err).Log("Error while AddTasks")
					return wfinstance.WFInstanceNewResponse{}, err
				}
				return response, nil
			} else {
				step = actionset.Tasks[0]
				UpdateWFInstanceStep(queries, DoMarkDoneParam, actionset.Tasks[0], ruleset.SetName)
				response = wfinstance.WFInstanceNewResponse{
					Tasks:    []map[string]int32{{step: DoMarkDoneParam.InstanceID}},
					Loggedat: pgtype.Timestamp{Time: DoMarkDoneParam.Loggedat, Valid: true},
				}
				return response, nil
			}

		} else {
			// We come here when our current step is one of a set of concurrent steps
			// and one or more of the other concurrent steps is yet to complete.
			// In that, we have nothing else to do other than mark the current step complete
			// and return to the caller saying "We have marked it done, there is nothing more
			// to do till one more of the other concurrent steps completes. Keep walking."
			// Return with details of success of mark-done.
			id := strconv.Itoa(int(DoMarkDoneParam.InstanceID))
			response := wfinstance.WFInstanceNewResponse{
				ID:       id,
				Loggedat: pgtype.Timestamp{Time: DoMarkDoneParam.Loggedat, Valid: true},
			}
			return response, nil
		}

	}
	return wfinstance.WFInstanceNewResponse{}, errors.New("schema Realmkey not match")
}
