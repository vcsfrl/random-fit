package cmd

import "fmt"

var errNoEnvEditor = fmt.Errorf("EDITOR environment variable is not set")

var msgNameMissing = "Name is required."
var msgCombinationDefinitionNameMissing = "Combination definition name is required."
var msgPlanDefinitionNameMissing = "Combination definition name is required."
var msgCombinationDefinition = "Combination Definition"
var msgPlanDefinition = "Plan Definition"
var msgList = "List"
var msgCreate = "Create"
var msgEdit = "Edit"
var msgDelete = "Delete"
var msgDone = "DONE:"
var msgEditScript = "Editing script"
var msgRemoveScript = "Removing script"
var msgNoItemsFound = "No items found!"
