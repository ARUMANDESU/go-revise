package tgbot

func (t *TgBot) setUpRouter() {
	t.bot.Handle("/start", t.handler.StartBot)
}
