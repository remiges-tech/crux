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
	"github.com/remiges-tech/logharbour/logharbour"

	"github.com/remiges-tech/alya/service"
)

const (
	doneProp = "done"
	WFE      = "W"
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
func DoMarkDone(c *gin.Context, s *service.Service, qtx *sqlc.Queries, instanceID int32, entity map[string]string) (wfinstance.WFInstanceNewResponse, error) {
	l := s.LogHarbour.WithClass("DoMarkDone")
	l.Debug1().Log("DoMarkDone function execution started")

	wfinst, err := qtx.GetWFInstanceFromId(c, instanceID)
	if err != nil {
		l.Debug1().Error(err)
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while GetWFInstanceFromId() in WFInstanceMarkDone")
	}
	cruxCache, ok := s.Dependencies["cruxCache"].(*crux.Cache)
	if !ok {
		l.Debug0().Debug1().Log("Error while getting cruxCache instance from service Dependencies")
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while getting cruxCache instance from service Dependencies")
	}

	stepfailed, err := strconv.ParseBool(entity["stepfailed"])
	if err != nil {
		l.Debug1().Error(err)
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while converting stepfailed val from string to bool")
	}
	step := entity["step"]

	schema, ruleset, err := cruxCache.RetriveRuleSchemasAndRuleSetsFromCache(WFE, wfinst.App, realmName, wfinst.Class, wfinst.Workflow, wfinst.Slice)
	if err != nil {
		l.Debug0().Error(err).Log("error while Retrieve RuleSchemas and RuleSets FromCache")
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while Retrieve RuleSchemas and RuleSets FromCache: %v", err)
	} else if schema == nil || ruleset == nil {
		l.Debug0().Error(err).Log("didn't find any data in RuleSchemas or RuleSets cache")
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("didn't find any data in RuleSchemas or RuleSets cache: ")
	}

	entity_t := crux.Entity{
		Realm: realmName,
		App:   wfinst.App,
		Slice: wfinst.Slice,
		Class: wfinst.Class,
		Attrs: entity,
	}
	err = crux.VerifyEntity(entity_t, schema)
	if err != nil {
		l.Debug0().Error(err).Log("error while verifying entityFromCache")
		return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("error while verifying entityFromCache: %v", err)
	}

	var response wfinstance.WFInstanceNewResponse

	if stepfailed {

		actionSet := crux.ActionSet{}
		seenRuleSets := make(map[string]struct{})

		// Call the doMatch function passing the entity.entity, ruleset, and the empty actionSet and seenRuleSets
		actionset, _, err, _ := crux.DoMatch(entity_t, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})

		if err != nil {
			l.Debug0().Error(err).Log("error while performing DoMatch")
			return wfinstance.WFInstanceNewResponse{}, err
		}

		if actionset.Properties[doneProp] == "true" {
			res, err := doneTrue(l, qtx, instanceID, entity_t, wfinst)
			if err != nil {
				l.Info().Error(err).Log("error while deleting wfinstance")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return res, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task

			// call addTasks

			entity["class"] = wfinst.Class

			task := wfinstance.AddTaskRequest{
				Steps:    actionset.Tasks,
				Nextstep: actionset.Properties["nextstep"],
				Request: wfinstance.WFInstanceNewRequest{
					Slice:    wfinst.Slice,
					App:      wfinst.App,
					EntityID: wfinst.Entityid,
					Entity:   entity,
					Workflow: wfinst.Workflow,
				},
			}

			response, err = wfinstance.AddTasks(task, s, c)
			if err != nil {
				l.Info().Error(err).Log("Error while AddTasks")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return response, nil
		} else {
			entity["step"] = actionset.Tasks[0]
			err = UpdateWFInstanceStep(qtx, instanceID, entity_t, actionset.Tasks[0], ruleset.SetName)
			if err != nil {
				l.Info().Error(err).Log("Error while Update WFInstance Step")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			response = wfinstance.WFInstanceNewResponse{
				Tasks:    []map[string]int32{{step: instanceID}},
				Loggedat: pgtype.Timestamp{Time: wfinst.Loggedat.Time, Valid: true},
			}
			return response, nil
		}
	}
	// We come here k	wing that the previous step didn't fail. We can now proceed to the next step; the previous step was successful
	recordcount, _ := GetWFInstanceCountForEntity(qtx, instanceID, entity_t, ruleset.SetName)
	if recordcount == 1 {
		// markDoneReq.Step = step

		// Invoke doMatch() with
		// entity = the object received
		// ruleset = the ruleset name retrieved from wfinstance
		// actionset and seenrulesets: empty
		actionSet := crux.ActionSet{}
		seenRuleSets := make(map[string]struct{})


		actionset, match, err, _ := crux.DoMatch(entity_t, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})

		if err != nil {
			l.Info().Error(err).Log("error while performing DoMatch")
			return wfinstance.WFInstanceNewResponse{}, err
		}

		if !match {
			return wfinstance.WFInstanceNewResponse{}, fmt.Errorf("rule are not matched")
		}
		if actionset.Properties[doneProp] == "true" {
			res, err := doneTrue(l, qtx, instanceID, entity_t, wfinst)
			if err != nil {
				l.Info().Error(err).Log("error while deleting wfinstance")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return res, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task
			// Return the full set of tasks and their record IDs

			err := deleteWFInstance(qtx, instanceID, entity_t)
			if err != nil {
				l.Info().Error(err).Log("Error while deleteWFInstance() in DoMarkDone")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			dclog := l.WithClass("WFInstance").WithInstanceId(string(instanceID))
			dclog.LogDataChange("insert ruleset", logharbour.ChangeInfo{
				Entity: "WFInstance",
				Op:     "delete",
				Changes: []logharbour.ChangeDetail{
					{
						Field:  "entityid",
						OldVal: nil,
						NewVal: wfinst.Entityid,
					},
					{
						Field:  "slice",
						OldVal: nil,
						NewVal: wfinst.Slice,
					},
					{
						Field:  "app",
						OldVal: nil,
						NewVal: wfinst.App,
					},
					{
						Field:  "class",
						OldVal: nil,
						NewVal: wfinst.Class,
					},
					{
						Field:  "workflow",
						OldVal: nil,
						NewVal: wfinst.Workflow,
					},
					{
						Field:  "step",
						OldVal: nil,
						NewVal: step,
					},
				},
			})

			doneAtTimeStamp := time.Now()
			err = UpdateWFInstanceDoneAt(qtx, instanceID, entity_t, doneAtTimeStamp, ruleset.SetName)
			if err != nil {
				l.Info().Error(err).Log("Error while update wfinstance Done At() in DoMarkDone")
				return wfinstance.WFInstanceNewResponse{}, err
			}

			// call addTasks

			entity["class"] = wfinst.Class

			addTaskParam := wfinstance.AddTaskRequest{
				Steps:    actionset.Tasks,
				Nextstep: actionset.Properties["nextstep"],
				Request: wfinstance.WFInstanceNewRequest{
					Slice:    wfinst.Slice,
					App:      wfinst.App,
					EntityID: wfinst.Entityid,
					Entity:   entity,
					Workflow: wfinst.Workflow,
				},
			}

			response, err = wfinstance.AddTasks(addTaskParam, s, c)
			if err != nil {
				l.Error(err).Log("Error while AddTasks")
				return wfinstance.WFInstanceNewResponse{}, err
			}
			return response, nil
		} else {
			step = actionset.Tasks[0]
			UpdateWFInstanceStep(qtx, instanceID, entity_t, actionset.Tasks[0], ruleset.SetName)
			response = wfinstance.WFInstanceNewResponse{
				Tasks:    []map[string]int32{{step: instanceID}},
				Loggedat: pgtype.Timestamp{Time: wfinst.Doneat.Time, Valid: true},
			}
			return response, nil
		}
	} else if recordcount > 1 {
		// At this point, we have found multiple records with the same entity ID and workflow, which means they differ only in the value of "step"
		// Set the doneat field in the current wfinstance record to the current timestamp
		doneAtTimeStamp := time.Now()
		err := UpdateWFInstanceDoneAt(qtx, instanceID, entity_t, doneAtTimeStamp, ruleset.SetName)
		if err != nil {
			l.Info().Error(err).Log("Error while update wfinstance Done At() in DoMarkDone")
			return wfinstance.WFInstanceNewResponse{}, err
		}
		// Look through all the other wfinstance records which have matching tuple (slice,app,workflow,entityid)
		wfInstances, err := getWFInstanceList(qtx, instanceID, entity_t, ruleset.SetName)
		if err != nil {
			l.Info().Error(err).Log("Error while getWFInstanceList() in DoMarkDone")
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

			entity_t.Attrs["step"] = wfinst.Nextstep


			actionset, _, err, _ := crux.DoMatch(entity_t, ruleset, schema, actionSet, seenRuleSets, crux.Trace_t{})

			if err != nil {
				l.Info().Error(err).Log("error while performing DoMatch")
				return wfinstance.WFInstanceNewResponse{}, err
			}

			if actionset.Properties[doneProp] == "true" {

				res, err := doneTrue(l, qtx, instanceID, entity_t, wfinst)
				if err != nil {
					l.Info().Error(err).Log("error while deleting wfinstance")
					return wfinstance.WFInstanceNewResponse{}, err
				}
				return res, nil
			}

			if len(actionset.Tasks) >= 1 {
				// Delete the old record from wfinstance and create fresh records, one per task
				// Return the full set of tasks and their record IDs
				err := deleteWFInstance(qtx, instanceID, entity_t)
				if err != nil {
					l.Info().Error(err).Log("Error while deleteWFInstance() in DoMarkDone")
					return wfinstance.WFInstanceNewResponse{}, err
				}
				dclog := l.WithClass("WFInstance").WithInstanceId(string(instanceID))
				dclog.LogDataChange("insert ruleset", logharbour.ChangeInfo{
					Entity: "WFInstance",
					Op:     "delete",
					Changes: []logharbour.ChangeDetail{
						{
							Field:  "entityid",
							OldVal: nil,
							NewVal: wfinst.Entityid,
						},
						{
							Field:  "slice",
							OldVal: nil,
							NewVal: wfinst.Slice,
						},
						{
							Field:  "app",
							OldVal: nil,
							NewVal: wfinst.App,
						},
						{
							Field:  "class",
							OldVal: nil,
							NewVal: wfinst.Class,
						},
						{
							Field:  "workflow",
							OldVal: nil,
							NewVal: wfinst.Workflow,
						},
						{
							Field:  "step",
							OldVal: nil,
							NewVal: step,
						},
					},
				})

				// call addTasks

				entity["class"] = wfinst.Class

				task := wfinstance.AddTaskRequest{
					Steps:    actionset.Tasks,
					Nextstep: actionset.Properties["nextstep"],
					Request: wfinstance.WFInstanceNewRequest{
						Slice:    wfinst.Slice,
						App:      wfinst.App,
						EntityID: wfinst.Entityid,
						Entity:   entity,
						Workflow: wfinst.Workflow,
					},
				}

				response, err = wfinstance.AddTasks(task, s, c)
				if err != nil {
					l.Error(err).Log("Error while AddTasks")
					return wfinstance.WFInstanceNewResponse{}, err
				}
				return response, nil
			} else {
				step = actionset.Tasks[0]
				UpdateWFInstanceStep(qtx, instanceID, entity_t, actionset.Tasks[0], ruleset.SetName)
				response = wfinstance.WFInstanceNewResponse{
					Tasks:    []map[string]int32{{step: instanceID}},
					Loggedat: pgtype.Timestamp{Time: wfinst.Doneat.Time, Valid: true},
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
			id := strconv.Itoa(int(instanceID))
			response := wfinstance.WFInstanceNewResponse{
				ID:       id,
				Loggedat: pgtype.Timestamp{Time: wfinst.Doneat.Time, Valid: true},
			}
			return response, nil
		}

	}
	return wfinstance.WFInstanceNewResponse{}, errors.New("schema Realmkey not match")
}
