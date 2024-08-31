# tennis-tg-bot 🏓

[![[branch]](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml/badge.svg)](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tubopo/tennis-tg-bot)](https://goreportcard.com/report/github.com/tubopo/tennis-tg-bot)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/tubopo/tennis-tg-bot/blob/main/LICENSE)

## schedule user flow

```mermaid
---
 title: Training schedule - create
---
flowchart LR
 appointment([Записаться на тернировку])
 appointment-->date[Выбрать дату]
 appointment-->place[Выбрать зал]

 date--"`Available Time Slots`"-->calendar[place, dd-mm-yy, hh:mm]

 place-->date
```

```mermaid
---
 title: Training schedule - view
---

flowchart LR
 view([Просмотреть тренировки])
 view-->date[Выбрать дату]
 date-->calendar[place, dd-mm-yy, hh:mm]
 calendar-->participants[Участники]

```
