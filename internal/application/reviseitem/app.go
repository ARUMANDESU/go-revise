package reviseitem

import (
	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem/command"
	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem/query"
)

type Application struct {
	Query   Query
	Command Command
}

type Query struct {
	GetReviseItem       query.GetReviseItemHandler
	ListUserReviseItems query.ListUserReviseItemsHandler
}

type Command struct {
	NewReviseItem     command.NewReviseItemHandler
	DeleteReviseItem  command.DeleteReviseItemHandler
	ChangeDescription command.ChangeDescriptionHandler
	ChangeName        command.ChangeNameHandler
	AddTags           command.AddTagsHandler
	RemoveTags        command.RemoveTagsHandler
	Review            command.ReviewHandler
}
