package server

type sharedData struct {
	*Config
}

// data that should be accessible across handlers
var shared *sharedData

func init() {
	shared = &sharedData{}
}
