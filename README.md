# linux-tutor

linux-tutor is an interactive Linux learning app written in Go. It combines a terminal UI, a graphical UI, an LPIC-style question engine, and progress tracking in SQLite.

## What it is

The app helps users practice Linux command-line skills in a structured way. It is built around LPIC-inspired topics such as GNU and Unix commands, process management, permissions, shell scripting, user administration, networking, and host security.

## Features

- Terminal UI built with Bubble Tea.
- Graphical UI built with Fyne.
- Adaptive question flow based on weak topics.
- LPIC topic catalog stored in JSON.
- Question generation and answer evaluation in the agent layer.
- SQLite-backed progress tracking.
- Separate lesson, question, and feedback views.

## Architecture

The project is organized by responsibility:

- `cmd/` — entry points for each executable.
  - `cmd/app/` — shared or default launcher.
  - `cmd/tui/` — terminal application entry point.
  - `cmd/gui/` — graphical application entry point.
- `internal/app/` — application wiring, configuration, and runtime setup.
- `internal/agent/` — question generation, evaluation rules, prompts, guidelines, and tutoring logic.
- `internal/catalog/` — LPIC topic catalog data.
- `internal/domain/` — core domain types such as topics, tasks, sessions, answers, and attempts.
- `internal/infra/` — persistence, repositories, and other infrastructure code.
- `internal/terminal/` — TUI model, screens, widgets, mini shell, and terminal flow.
- `internal/gui/` — GUI application, theme, and desktop flow.
- `internal/usecase/` — application use cases such as starting lessons, generating questions, evaluating progress, and recommending next topics.

This structure keeps UI concerns, business logic, and infrastructure separate, which makes the project easier to extend and test.

## How it works

A typical learning flow looks like this:

1. The app loads the LPIC catalog.
2. The agent generates a question for a selected topic.
3. The user answers in the terminal or GUI.
4. The answer is evaluated and scored.
5. Progress is saved to SQLite.
6. The app recommends the next topic, often favoring weaker areas.

## Learning model

The learning model is designed around topic-level practice rather than random quizzes. Each topic can generate different task types, including:

- single command.
- multi-command.
- fill in the blank.
- multiple choice.
- ordering.
- scenario-based tasks.

The app keeps track of correct and wrong answers so it can focus on weaker topics over time.

## Getting started

### Requirements

- Go 1.22 or newer.
- SQLite support.
- A terminal for the TUI version.
- A desktop environment for the GUI version.

### Run the app

```bash
go run ./cmd/app
```

### Build the binary

```bash
go build -o bin/linux-tutor ./cmd/app
```

### Run the TUI

```bash
go run ./cmd/tui
```

### Run the GUI

```bash
go run ./cmd/gui
```

## Project layout

```text
linux-tutor/
├─ cmd/
│  ├─ app/
│  ├─ gui/
│  └─ tui/
├─ internal/
│  ├─ agent/
│  ├─ app/
│  ├─ catalog/
│  ├─ domain/
│  ├─ gui/
│  ├─ infra/
│  ├─ terminal/
│  └─ usecase/
├─ Makefile
└─ README.md
```

## Scoring

Answer evaluation uses a simple rubric:

- exact answer.
- partial answer.
- wrong answer.

This makes it possible to distinguish between full understanding and partial familiarity, especially in shell and scenario-based tasks.

## Contributing

Contributions are welcome. If you add new topics, update the catalog. If you add a new task type, extend the agent and use-case layers consistently. If you change persistence, keep the repository and infra layers aligned.

## License

MIT
