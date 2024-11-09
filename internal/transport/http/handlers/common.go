package handlers

type H map[string]interface{}

func throw(msg string) *H {
	return &H{"error": msg}
}
