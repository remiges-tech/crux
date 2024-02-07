/* This file contains the collectActions() function */

package main

func collectActions(actionSet ActionSet, ruleActions RuleActions) ActionSet {

	newActionSet := ActionSet{
		tasks:      []string{},
		properties: make(map[string]string),
	}

	// Union-set of tasks
	newActionSet.tasks = append(newActionSet.tasks, actionSet.tasks...)
	for _, newTask := range ruleActions.tasks {
		found := false
		for _, task := range newActionSet.tasks {
			if newTask == task {
				found = true
				break
			}
		}
		if !found {
			newActionSet.tasks = append(newActionSet.tasks, newTask)
		}
	}


	// Perform "union-set" of properties, overwriting previous property values if needed

	for name, val := range actionSet.properties {
		newActionSet.properties[name] = val
	}

	// Update properties from ruleActions
	for propName, propertyVal := range ruleActions.properties {
		found := false
		for existingPropName := range newActionSet.properties {
			if existingPropName == propName {
				newActionSet.properties[existingPropName] = propertyVal
				found = true
				break
			}
		}
		if !found {
			newActionSet.properties[propName] = propertyVal
		}
	}

	return newActionSet
}
