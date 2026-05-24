package gui

import (
	"fmt"
	"linux-tutor/internal/domain"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func buildQuestionView(task domain.Question) fyne.CanvasObject {
	questionHeader := NewQuestionTitle("Question")
	questionBody := NewMutedText(renderQuestionBody(task))
	topic := NewMutedText(fmt.Sprintf("Topic: %s", task.Topic.Title))
	area := NewMutedText(fmt.Sprintf("Area: %s", task.Topic.Area))
	kind := NewMutedText(fmt.Sprintf("Kind: %s", task.Kind))
	hint := NewMutedText(fmt.Sprintf("Hint: %s", task.Hint))
	return container.NewVBox(questionHeader, questionBody, topic, area, kind, hint)
}

func renderQuestionBody(task domain.Question) string {
	if task.ID == "" {
		return "No question loaded yet."
	}
	parts := []string{task.Prompt}
	if len(task.Choices) > 0 {
		parts = append(parts, "Choices:")
		for i, c := range task.Choices {
			parts = append(parts, fmt.Sprintf("%d) %s", i+1, c))
		}
	}
	return strings.Join(parts, "")
}
