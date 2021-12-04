package function

type RequestData struct {
	Version string `json:"version"`
	Author  string `json:"author"`
	Text    string `json:"text"`
	Color   string `json:"color"`
}

type Comment struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
	Color  string `json:"color"`
	Date   string `json:"date"`
	Reply  string `json:"reply"`
}

type CommentConfig struct {
	Version    string `json:"version"`
	FirstBlock int    `json:"firstBlock"`
	LastBlock  int    `json:"lastBlock"`
}
