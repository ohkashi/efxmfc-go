<img src="true_friend.png">

# efxMFC-go
a go wrapper for [`EfxMfc`](https://github.com/ohkashi/EfxMfc)  
한국투자증권 `eFriend Expert` Host DLL을 위한 Go 패키지

[![Licence](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

### Example

This is a simple request tr and realdata. ([main.go](main.go))

```go
package main

import (
    efx "efxmfc"
    "fmt"
    "os"
    "syscall"
    "time"
)

...

func main() {
    acc_cnt := efx.GetAccountCount()
    fmt.Printf("AccountCount: %d\n", acc_cnt)
    for i := 0; i < acc_cnt; i++ {
        acnt_no := efx.GetAccount(i)
        fmt.Printf("Account[%d]: %s, %s\n", i, acnt_no, efx.GetAccountBrcode(acnt_no))
    }
    efx1 := efx.NewControl(100, OnRecvData, OnRecvRealData, OnRecvError)

    var pwd string
    fmt.Print("Enter Password: ")
    fmt.Scanln(&pwd)
    wait_time := time.Now()
    if len(pwd) == 4 {
        fmt.Printf("Encrypt: %s\n\n", efx.GetEncryptPassword(pwd))
        sync_time := efx.Synchonize(efx.REQ_LIMIT_MS)
        fmt.Printf("Synchonize() elapsed time: %v, sync time: %dms\n", time.Since(wait_time), sync_time)
        efx.RequestAccountBalance(efx1, efx.GetAccount(acc_cnt-1), pwd)
        wait_time = time.Now()
        wait_recv_data()
    }
    fmt.Println()

    go func() {
        stock_code := []string{"254120", "081000", "005930", "294090", "149950", "052670", "078940"}
        for _, code := range stock_code {
            sync_time := efx.Synchonize(efx.REQ_LIMIT_MS)
            fmt.Printf("Synchonize() elapsed time: %v, sync time: %dms\n",
                time.Since(wait_time), sync_time)
            efx.SetSingleData(efx1, 0, "J")
            efx.SetSingleData(efx1, 1, code)
            efx.RequestData(efx1, "SCP", code)
            wait_time = time.Now()
            wait_recv_data()
        }

        var real_code string
        for _, code := range stock_code {
            str := code + "   "
            real_code += str
        }
        efx.RequestRealData(efx1, "SC_R", real_code)
        fmt.Printf("--> RequestRealData(\"SC_R\", \"%s\")\n", real_code)
    }()

    //efx.Quit(0)
    efx.MessageLoop()
    efx.Exit()
}
```
