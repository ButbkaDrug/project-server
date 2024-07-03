package main

type TableOfContentEntry struct {
	Title string
	Link  string
}

type Meta struct {
	Title          string
	TableOfContent []TableOfContentEntry
}

type Page struct {
	Page Meta
	Data interface{}
	Map  interface{}
}
