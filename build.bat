go build -o efx-console.exe -ldflags="-s -w" -trimpath main.go
go build -ldflags "-s -w -H windowsgui" walk-chart.go
