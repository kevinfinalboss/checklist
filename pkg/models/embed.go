package models

type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"`
	Footer      *Footer      `json:"footer,omitempty"`
	Image       *Image       `json:"image,omitempty"`
	Thumbnail   *Thumbnail   `json:"thumbnail,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

type Footer struct {
	Text string `json:"text"`
}

type Image struct {
	URL string `json:"url"`
}

type Thumbnail struct {
	URL string `json:"url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
