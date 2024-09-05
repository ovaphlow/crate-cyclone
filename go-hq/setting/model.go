package setting

type Setting struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	RootID   string `json:"rootId"`
	ParentID string `json:"parentId"`
	Tags     string `json:"tags"`
	Name     string `json:"name"`
	Detail   string `json:"detail"`
}
