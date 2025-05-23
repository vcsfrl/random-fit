package cmd

import (
	"errors"
)

var ErrNoEnvEditor = errors.New("EDITOR environment variable is not set")

var MsgNameMissing = "Name is required."
var MsgCombinationDefinitionNameMissing = "Combination definition name is required."
var MsgPlanDefinitionNameMissing = "Combination definition name is required."
var MsgCombinationDefinition = "Combination Definition"
var MsgPlanDefinition = "Plan Definition"
var MsgList = "List"
var MsgCreate = "Create"
var MsgEdit = "Edit"
var MsgDelete = "Delete"
var MsgDone = "DONE:"
var MsgEditScript = "Editing script"
var MsgNoItemsFound = "No items found!"
