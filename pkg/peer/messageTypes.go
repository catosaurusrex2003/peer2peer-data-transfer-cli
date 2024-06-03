package peer

type FileMetadata_MessageType struct {
	Action   string `json:"action"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
}

type MessageMetadata_MessageType struct {
	Action string `json:"action"`
}
