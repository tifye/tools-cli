package winmower

type Manifest struct {
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Releasenotes string   `json:"releasenotes"`
	Metadata     struct {
		Origin                string   `json:"origin"`
		Version               string   `json:"version"`
		UniqueDescriptiveName string   `json:"uniqueDescriptiveName"`
		Platforms             []string `json:"platforms"`
		CreationDate          string   `json:"creationDate"`
		Type                  string   `json:"type"`
		OriginalFilename      string   `json:"originalFilename"`
		GitHash               string   `json:"gitHash"`
	} `json:"metadata"`
}
