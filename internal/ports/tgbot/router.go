package tgbot

func (p *Port) setUpRouter() {
	p.bot.Handle("/start", p.handler.StartBot)
	p.bot.Handle("/register_user", p.handler.RegisterUser)
}
