package crux

import (
	"errors"
	"time"
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

func doMarkDone(entity markdone_t, step string, SetName string) (markdone_t, []response_t, error) {
	ruleset := retriveRuleSetsFromCache(entity.Entity.Realm, entity.Entity.App, entity.Entity.Class, entity.Entity.Slice)
	if entity.Stepfailed == true {
		// Step supplied in the request
		entity.Step = step

		actionSet := ActionSet{}
		seenRuleSets := make(map[string]struct{})

		// Call the doMatch function passing the entity.entity, ruleset, and the empty actionSet and seenRuleSets
		actionset, _, err := DoMatch(entity.Entity, ruleset, actionSet, seenRuleSets)
		if err != nil {
			return markdone_t{}, []response_t{}, err
		}

		if actionset.Properties[done] == "done" {
			// Delete the wfinstance record
			err := deleteWFInstance(entity)
			if err != nil {
				entity.Stepfailed = true
				return markdone_t{}, []response_t{}, err
			} else {
				entity.Step = done
			}
			return entity, []response_t{}, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task
			err := deleteWFInstance(entity)
			if err != nil {
				entity.Stepfailed = true
				return markdone_t{}, []response_t{}, err
			}
			var response []response_t

			for _, freshtask := range actionset.Tasks {
				newrecord, err := createFreshRecord(entity, SetName, freshtask, actionset.Properties)
				if err != nil {
					if len(newrecord) > 0 {
						record := newrecord[0] // Assuming only one record is returned
						taskMap := make(map[string]int32)
						taskMap[freshtask] = record.ID

						param := response_t{
							Task:     taskMap,
							Loggedat: record.Loggedat.Time, // Assuming Loggedat is a pgtype.Timestamp
						}

						response = append(response, param)
					}
				}
			}
			return entity, response, nil
		} else {
			// Update the old record in wfinstance to set the value of "step" = the task returned
			entity.Step = actionset.Tasks[0]
			UpdateWFInstanceStep(entity, actionset.Tasks[0])
			return entity, []response_t{}, nil // Return the task and other data in response
		}
	}

	// We come here knowing that the previous step didn't fail. We can now proceed to the next step; the previous step was successful
	recordcount, _ := GetWorkFlowInstance(entity, ruleset.SetName)
	if recordcount == 1 {
		entity.Step = step

		// Invoke doMatch() with
		// entity = the object received
		// ruleset = the ruleset name retrieved from wfinstance
		// actionset and seenrulesets: empty
		actionSet := ActionSet{}
		seenRuleSets := make(map[string]struct{})

		actionset, _, err := DoMatch(entity.Entity, ruleset, actionSet, seenRuleSets)
		if err != nil {
			return markdone_t{}, []response_t{}, err
		}

		if actionset.Properties[done] == "true" {
			// Delete the wfinstance record
			// Return specifying that the workflow is completed
			err := deleteWFInstance(entity)
			if err != nil {
				entity.Stepfailed = true
			} else {
				entity.Step = done
			}
			return entity, []response_t{}, nil
		}

		if len(actionset.Tasks) > 1 {
			// Has more than one task then delete the old record from wfinstance and create fresh records, one per task
			// Return the full set of tasks and their record IDs
			err := deleteWFInstance(entity)
			if err != nil {
				entity.Stepfailed = true
				return markdone_t{}, []response_t{}, err
			}
			var response []response_t

			for _, freshtask := range actionset.Tasks {
				newrecord, err := createFreshRecord(entity, ruleset.SetName, freshtask, actionset.Properties)
				if err != nil {
					if len(newrecord) > 0 {
						record := newrecord[0] // at time only one record made
						taskMap := make(map[string]int32)
						taskMap[freshtask] = record.ID

						param := response_t{
							Task:     taskMap,
							Loggedat: record.Loggedat.Time, // Assuming Loggedat is a pgtype.Timestamp
						}

						response = append(response, param)
					}
				}
			}
			return entity, response, nil
		} else {
			// Update the old record in wfinstance to set the value of "step" to the task returned
			var response []response_t
			taskMap := make(map[string]int32)
			taskMap[actionset.Tasks[0]] = entity.Id

			param := response_t{
				Task:     taskMap,
				Loggedat: time.Now(), // Assuming Loggedat is a pgtype.Timestamp
			}

			response = append(response, param)
			UpdateWFInstanceStep(entity, actionset.Tasks[0])
			return entity, response, nil // Return the task and other data in response
		}
	} else if recordcount > 1 {
		// At this point, we have found multiple records with the same entity ID and workflow, which means they differ only in the value of "step"
		// Set the doneat field in the current wfinstance record to the current timestamp
		err := UpdateWFInstanceDoneAt(entity, time.Now(), ruleset.SetName)
		if err != nil {
			return markdone_t{}, []response_t{}, err
		}
		// Look through all the other wfinstance records which have matching tuple (slice,app,workflow,entityid)
		wfInstances, err := getWFInstanceList(entity, ruleset.SetName)
		if err != nil {
			return markdone_t{}, []response_t{}, err
		}

		// Check if all other wfinstance records have doneat set
		allDone := true
		for _, wfInstance := range wfInstances {
			v, err := wfInstance.Doneat.Value()

			if err != nil || v == 0 {
				allDone = false
				break
			}
		}

		if allDone {
			wfinstance, err := getCurrentWFINstance(entity, ruleset.SetName)
			if err == nil {
				return markdone_t{}, []response_t{}, err
			}
			entity.Step = wfinstance.Nextstep // The value of "nextstep" from the current wfinstance record
			// Invoke doMatch() with
			//  entity = the object received
			// ruleset = the ruleset name retrieved from wfinstance
			// actionset and seenrulesets: empty
			actionSet := ActionSet{}
			seenRuleSets := make(map[string]struct{})

			actionset, _, err := DoMatch(entity.Entity, ruleset, actionSet, seenRuleSets)
			if err != nil {
				return markdone_t{}, []response_t{}, err
			}

			if actionset.Properties[done] == "true" {
				// Delete all wfinstance records with tuple matching (slice, app, workflow, entityid)
				// Return specifying that the workflow is completed
				err := deleteWFInstance(entity)
				if err != nil {
					entity.Stepfailed = true
				} else {
					entity.Step = done
				}
				return entity, []response_t{}, nil
			}

			if len(actionset.Tasks) > 1 {
				// Delete the old record from wfinstance and create fresh records, one per task
				// Return the full set of tasks and their record IDs
				err := deleteWFInstance(entity)
				if err != nil {
					entity.Stepfailed = true
					return markdone_t{}, []response_t{}, err
				}
				var response []response_t

				for _, freshtask := range actionset.Tasks {
					newrecord, err := createFreshRecord(entity, ruleset.SetName, freshtask, actionset.Properties)
					if err != nil {
						if len(newrecord) > 0 {
							record := newrecord[0] // Assuming only one record is returned
							taskMap := make(map[string]int32)
							taskMap[freshtask] = record.ID

							param := response_t{
								Task:     taskMap,
								Loggedat: record.Loggedat.Time, // Assuming Loggedat is a pgtype.Timestamp
							}

							response = append(response, param)
						}
					}
				}
				return entity, response, nil
			} else {
				// Update the old record in wfinstance to set the value of "step" to the task returned
				// Return the task and other data in response
				var response []response_t
				taskMap := make(map[string]int32)
				taskMap[actionset.Tasks[0]] = entity.Id

				param := response_t{
					Task:     taskMap,
					Loggedat: time.Now(), // Assuming Loggedat is a pgtype.Timestamp
				}

				response = append(response, param)
				UpdateWFInstanceStep(entity, actionset.Tasks[0])
				return entity, response, nil // Return the task and other data in response

			}
		} else {
			// We come here when our current step is one of a set of concurrent steps
			// and one or more of the other concurrent steps is yet to complete.
			// In that, we have nothing else to do other than mark the current step complete
			// and return to the caller saying "We have marked it done, there is nothing more
			// to do till one more of the other concurrent steps completes. Keep walking."
			// Return with details of success of mark-done.
			return entity, []response_t{}, nil
		}
	}
	return markdone_t{}, []response_t{}, errors.New("schema Realmkey not match")
}
