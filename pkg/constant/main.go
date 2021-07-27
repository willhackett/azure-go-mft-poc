package constant

var (
	PublicKeyContainerName = "publickeys"
)

func AgentKeyName(agentName string, keyID string) string {
	return agentName + "/" + keyID
}