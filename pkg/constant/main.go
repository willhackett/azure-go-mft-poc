package constant

const (
	PublicKeyContainerName = "publickeys"

	MaxConcurrentTransfers = 4

	MaxRetriesThreshold = 1
)

func AgentKeyName(agentName string, keyID string) string {
	return agentName + "/" + keyID
}

func VerifierString(m Message) []byte {
	return []byte(m.ID + m.KeyID + m.Agent + m.Type + string(m.Payload))
}

func StringInList(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
