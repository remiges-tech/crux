/* This file contains the collectActions() function */

package main

func collectActions(actionSet ActionSet, ruleActions RuleActions) ActionSet {
	newActionSet := ActionSet{}

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
	newActionSet.properties = append(newActionSet.properties, actionSet.properties...)
	for _, newProperty := range ruleActions.properties {
		found := false
		for i, property := range newActionSet.properties {
			if property.name == newProperty.name {
				newActionSet.properties[i].val = newProperty.val
				found = true
				break
			}
		}
		if !found {
			newActionSet.properties = append(newActionSet.properties, newProperty)
		}
	}
	return newActionSet
}
