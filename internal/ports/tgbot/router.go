package tgbot

import "github.com/ARUMANDESU/go-revise/internal/ports/tgbot/button"

func (p *Port) setUpRouter() {
	p.bot.Handle("/start", p.handler.StartBot)

	p.bot.Handle("/register", p.handler.RegisterUser)
	p.bot.Handle(&button.RegistrationConfirmI, p.handler.RegisterUserConfirmed)

	p.bot.Handle("/revise_create", p.handler.CreateItem)
}
