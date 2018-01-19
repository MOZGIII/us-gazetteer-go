package usgazetteer

//go:generate go run ../internal/generate/main.go states ../sources/states.txt ../internal/generated/states.go
//go:generate gofmt -w ../internal/generated/states.go
//go:generate go run ../internal/generate/main.go counties ../sources/counties.txt ../internal/generated/counties.go
//go:generate gofmt -w ../internal/generated/counties.go
