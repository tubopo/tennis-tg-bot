# tennis-tg-bot 🏓

[![[branch]](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml/badge.svg)](https://github.com/tubopo/tennis-tg-bot/actions/workflows/branch.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tubopo/tennis-tg-bot)](https://goreportcard.com/report/github.com/tubopo/tennis-tg-bot)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/tubopo/tennis-tg-bot/blob/main/LICENSE)

## schedule user flow

```mermaid
---
 title: Training schedule - create
---
graph TD
    A[Start] --> B{User Command}
    B -->|/start| C[Welcome Message]
    B -->|/new_training| D[Send Place Selection]
    D --> E{User Selects Place}
    E -->|Смолячкова, 9| F[Show Time Slots for Смолячкова, 9]
    E -->|Ленина, 27| G[Show Time Slots for Ленина, 27]
    F --> H{User Selects Time Slot}
    G --> H
    H -->|Уровень 1| I[Confirm Training Level 1]
    H -->|Уровень 2| J[Confirm Training Level 2]
    I --> K[End]
    J --> K
    C --> K
```

```mermaid
---
 title: Training schedule - view
---

flowchart LR
 view([View Training])
 view --> date[Select date]
 date --> calendar[place, dd-mm-yy, hh:mm]
 calendar --> participants[View Participants]
```
