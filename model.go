package main

type Url struct {
	Loc     string `xml:"loc"`
	checked bool   `xml:"-"`
}
