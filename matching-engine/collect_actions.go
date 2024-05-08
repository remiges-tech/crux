/* This file contains the collectActions() function */

package crux

func collectActions(actionSet ActionSet, ruleActions RuleActionBlock_t) ActionSet {

	newActionSet := ActionSet{
		Tasks:      []string{},
		Properties: make(map[string]string),
	}

	// Union-set of Tasks
	newActionSet.Tasks = append(newActionSet.Tasks, actionSet.Tasks...)
	for _, newTask := range ruleActions.Task {
		found := false
		for _, task := range newActionSet.Tasks {
			if newTask == task {
				found = true
				break
			}
		}
		if !found {
			newActionSet.Tasks = append(newActionSet.Tasks, newTask)
		}
	}

	// Perform "union-set" of Properties, overwriting previous property values if needed

	for name, val := range actionSet.Properties {
		newActionSet.Properties[name] = val
	}

	// Update Properties from ruleActions
	for propName, propertyVal := range ruleActions.Properties {
		found := false
		for existingPropName := range newActionSet.Properties {
			if existingPropName == propName {
				newActionSet.Properties[existingPropName] = propertyVal
				found = true
				break
			}
		}
		if !found {
			newActionSet.Properties[propName] = propertyVal
		}
	}

	return newActionSet
}
