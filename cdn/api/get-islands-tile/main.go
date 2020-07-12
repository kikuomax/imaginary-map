package main

import (
	"context"
	runtime "github.com/aws/aws-lambda-go/lambda"
)

// Common form of a GetTileEvent.
type GetTileEvent struct {
	Zoom int `json:"zoom"`
	X int `json:"x"`
	Y int `json:"y"`
}

func HandleRequest (ctx context.Context, event GetTileEvent) ([]byte, error) {
	return []byte{ 1, 2, 3 }, nil
}

func main () {
	runtime.Start(HandleRequest)
}
