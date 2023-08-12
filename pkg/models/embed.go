package models

type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Color       int     `json:"color,omitempty"`
	Footer      *Footer `json:"footer,omitempty"`
}

type Footer struct {
	Text string `json:"text"`
}
