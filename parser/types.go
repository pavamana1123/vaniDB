package main

// Text ...
type Text struct {
	Info        Info        `json:"info"`
	Verses      []Verse     `json:"verses"`
	Synonyms    string      `json:"synonyms"`
	Translation string      `json:"translation"`
	Purport     []Paragraph `json:"purport,omitempty"`
}

// Info ...
type Info struct {
	ID     string `json:"id"`
	NextID string `json:"nextId"`
	PrevID string `json:"prevId"`
}

// Verse ...
type Verse struct {
	Roman      string `json:"roman"`
	Devanagari string `json:"devanagari,omitempty"`
	IsProse    bool   `json:"isProse"`
}

// Paragraph ...
type Paragraph struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
