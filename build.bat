go build -o efx-console.exe -ldflags="-s -w" -trimpath ./examples/efx-console/main.go
go build -o ./examples/walk-chart/walk-chart.exe -ldflags "-s -w -H windowsgui" -trimpath ./examples/walk-chart/walk-chart.go
