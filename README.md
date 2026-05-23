# linux-tutor

An interactive Linux learning app written in Go.  
It combines a Bubble Tea terminal UI, a small built-in training shell, SQLite progress tracking, and an LPIC-aligned AI tutor that generates and grades practice tasks.

## What it does

linux-tutor helps users practice Linux command-line skills in a structured way.  
It is designed around LPIC-style topics such as GNU/Unix commands, files and permissions, shell scripting, administration, networking, and security [web:6][web:1].

### Core features

- Terminal UI built with Bubble Tea.
- Built-in mini shell for guided practice.
- Multiple task types: single command, multi-command, fill-in-the-blank, multiple choice, ordering, and scenario-based tasks.
- LPIC topic catalog in JSON/YAML format.
- AI agent that generates tasks and evaluates answers.
- SQLite persistence for progress tracking.
- Adaptive learning that focuses on weak topics.

## Why this project exists

Many Linux learning tools focus only on quizzes or only on command references.  
linux-tutor combines both: it lets learners practice, make mistakes safely, and track progress over time in one interactive workflow.

## Getting started

### Requirements

- Go 1.23 or newer.
- SQLite support through the included pure-Go driver.
- A terminal that supports keyboard interaction.

### Run locally

```bash
git clone https://github.com/your-username/linux-tutor.git
cd linux-tutor
go run ./cmd/app
```

### Build

```bash
go build -o linux-tutor ./cmd/app
```

### Basic usage

When the app starts, use the arrow keys to navigate the dashboard.

Inside the mini shell, try:

```text
help
task
type
topic
next
answer pwd
```

## Project structure

```text
linux-tutor/
├─ cmd/app/                  # Application entry point
├─ internal/app/             # Bootstrapping and wiring
├─ internal/agent/           # AI tutor logic and LPIC guidelines
├─ internal/catalog/         # LPIC topic catalog in JSON/YAML
├─ internal/domain/          # Core domain models
├─ internal/infra/           # SQLite storage and repositories
├─ internal/terminal/        # Bubble Tea UI and mini shell
└─ docs/                     # Design notes and specification
```

## Learning model

The app uses an LPIC-oriented learning model:

- presents one topic at a time,
- generates practice tasks by topic code,
- grades answers with a scoring rubric,
- tracks weak areas,
- prioritizes weaker topics in the next round.

The topic coverage follows the LPIC-1 exam areas, including system architecture, commands, filesystems, shell scripting, administrative tasks, networking, and security [web:1][web:6].

## Scoring

Answer evaluation uses a simple rubric:

- exact answer,
- partial answer,
- wrong answer.

This makes it easier to distinguish between complete understanding and partial familiarity, especially for scenario-based Linux tasks.

## Roadmap

Planned improvements include:

- more task variants for each LPIC topic,
- richer shell labs,
- per-topic analytics,
- spaced repetition,
- user profiles,
- AI-generated explanations,
- exam mode with timers and weighted scoring.

## Contributing

Contributions are welcome.  
Please open an issue or pull request if you want to add new tasks, improve LPIC topic coverage, or refine the UI.

## License

Add your chosen license here.