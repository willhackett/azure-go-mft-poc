package constant

import "encoding/json"

const (
	FileRequestMessageType = "FileRequest"

	FileHandshakeMessageType = "FileHandshake"

	FileHandshakeResponseMessageType = "FileHandshakeResponse"

	FileAvailableMessageType = "FileAvailable"
)

// Message contains the overall structure of all messages sent to the queue
type Message struct {
	ID        string          `json:"id"`
	KeyID     string          `json:"key_id"`
	Agent     string          `json:"agent"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Signature string          `json:"signature"`
}

// FileRequestMessage contains the structure of the file request message
type FileRequestMessage struct {
	FileName            string `json:"file_name"`
	DestinationAgent    string `json:"destination_agent"`
	DestinationFileName string `json:"destination_file_name"`
}

// FileHandshakeMessage contains the structure of the file handshake message
type FileHandshakeMessage struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
}

// FileHandshakeResponseMessage contains the structure of the file handshake response message
type FileHandshakeResponseMessage struct {
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason"`
}

// FileAvailableMessage contains the structure of the file available message
type FileAvailableMessage struct {
	SignedURL string `json:"signed_url"`
	FileName  string `json:"file_name"`
}
