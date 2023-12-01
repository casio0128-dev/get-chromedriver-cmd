package main

type ChromeDriverData struct {
	Timestamp  string               `json:"timestamp"`
	Milestones map[string]Milestone `json:"milestones"`
}

type Milestone struct {
	Milestone string    `json:"milestone"`
	Version   string    `json:"version"`
	Revision  string    `json:"revision"`
	Downloads Downloads `json:"downloads"`
}

type Downloads struct {
	Chrome              []PlatformDownload `json:"chrome"`
	ChromeDriver        []PlatformDownload `json:"chromedriver"`
	ChromeHeadlessShell []PlatformDownload `json:"chrome-headless-shell,omitempty"` // omitempty if this field may not be present
}

type PlatformDownload struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}
