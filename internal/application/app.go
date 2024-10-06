package application

import (
	"github.com/ARUMANDESU/go-revise/internal/application/notification"
	"github.com/ARUMANDESU/go-revise/internal/application/reviseitem"
	"github.com/ARUMANDESU/go-revise/internal/application/user"
)

type Application struct {
	User         user.Application
	ReviseItem   reviseitem.Application
	Notification notification.Application
}
