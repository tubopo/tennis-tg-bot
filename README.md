# tennis-tg-bot 🏓

[![[branch]](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml/badge.svg)](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tubopo/tennis-tg-bot)](https://goreportcard.com/report/github.com/tubopo/tennis-tg-bot)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/tubopo/tennis-tg-bot/blob/main/LICENSE)

## schedule user flow

```mermaid
---
       title: Trainning schedule user flow
---
flowchart TD
    start([Записаться на тернировку])
    start-->date[Выбрать дату]
    start-->place[Выбрать зал]

    date--"`Avalible Time Slots`"-->calendar[place, dd-mm-yy, hh:mm]

    place-->date
```
