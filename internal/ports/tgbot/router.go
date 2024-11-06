package tgbot

func (p *Port) setUpRouter() {
	p.bot.Handle("/start", p.handler.StartBot)
}
