package vectordb

var db = map[string]string{
	"hello":   "Hi there! This is context for 'hello'.",
	"weather": "Today's weather is sunny.",
	"agent":   "Agents are autonomous entities.",
}

func Retrieve(query string) string {
	for k, v := range db {
		if k == query {
			return v
		}
	}
	return "No relevant context found."
}
