package main


import "askvart.ru/snippetbox/pkg/models"


type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}