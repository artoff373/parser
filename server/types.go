package main

import "html/template"

type roots struct {
	Name template.HTML
	ID   template.HTML
}
type rootHTML []roots

type profileHTML struct {
	ID      template.HTML
	Name    template.HTML
	Keys    []template.HTML
	Last    template.HTML
	Sources []sourceHTML
}
type sourceHTML struct {
	ID       template.HTML
	Name     template.HTML
	URL      template.HTML
	Selector template.HTML
}

type postHTML struct {
	Postik []postiks
}
type postiks struct {
	Id      template.HTML
	Title   template.HTML
	PubDate template.HTML
	Text    template.HTML
	Link    template.HTML
	Relev   template.HTML
	Source  template.HTML
}

type DbHtml struct {
	Host     template.HTML
	Port     template.HTML
	User     template.HTML
	Password template.HTML
	DbName   template.HTML
}

type Archive struct {
	Read template.HTML
	Val  template.HTML
}
