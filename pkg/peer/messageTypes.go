package peer

type FileMetadata_MessageType struct {
	Action       string `json:"action"` //will be fileInfoSend
	FileName     string `json:"fileName"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	IsDir        bool   `json:"isDir"`
}

type FileAccept_MessageType struct {
	Action string `json:"action"` // will be "fileAcceptMessage"
	Value  string `json:"value"`  // can be "accept" or "deny"
}

type MessageMetadata_MessageType struct {
	Action  string `json:"action"`
	Content string `json:"content"`
}
