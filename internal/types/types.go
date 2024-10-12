package types

import (
	"time"
)

// cli related.

type CliFlags struct {
	BaseUrl         string
	CookieDirectory string
	CookieFile      string
	DisplayResults  bool
	GameName        string
	ModId           int64
	OutputDirectory string
	SaveResults     bool
}

func NewScraper() *CliFlags {
	return &CliFlags{}
}

// end cli related.

// nexus mods related.
type Results struct {
	Mods ModInfo `json:"Mods"`
}

type ModInfo struct {
	ChangeLogs       []ChangeLog   `json:"ChangeLogs,omitempty"`
	Creator          string        `json:"Creator,omitempty"`
	Dependencies     []Requirement `json:"Dependencies,omitempty"`
	Description      string        `json:"Description,omitempty"`
	Files            []File        `json:"Files,omitempty"`
	LastChecked      time.Time     `json:"LastChecked,omitempty"`
	LastUpdated      string        `json:"LastUpdated,omitempty"`
	LatestVersion    string        `json:"LatestVersion,omitempty"`
	ModId            int64         `json:"ModID,omitempty"`
	ModsUsing        []Requirement `json:"ModsUsing,omitempty"`
	Name             string        `json:"Name,omitempty"`
	OriginalUpload   string        `json:"OriginalUpload,omitempty"`
	ShortDescription string        `json:"ShortDescription,omitempty"`
	Tags             []string      `json:"Tags,omitempty"`
	Uploader         string        `json:"Uploader,omitempty"`
	Url              string        `json:"Url,omitempty"`
	VirusStatus      string        `json:"VirusStatus,omitempty"`
}

type ChangeLog struct {
	Notes   []string `json:"Notes,omitempty"`
	Version string   `json:"Version,omitempty"`
}

type Requirement struct {
	Name  string `json:"Name,omitempty"`
	Notes string `json:"Notes,omitempty"`
}

type Tag struct {
	Tag string `json:"Tag,omitempty"`
}

type File struct {
	Description string `json:"description"`
	FileSize    string `json:"fileSize"`
	Name        string `json:"name"`
	TotalDLs    string `json:"totalDownloads"`
	UniqueDLs   string `json:"uniqueDownloads"`
	UploadDate  string `json:"uploadDate"`
	Version     string `json:"version"`
}

// end nexus mods related.
