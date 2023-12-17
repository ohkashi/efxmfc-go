package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	efx "github.com/ohkashi/efxmfc-go"
)

var _sig_recv_data bool = false

func OnRecvSysMsg(_handle syscall.Handle, ctrl_id uint, sys_msg int) {
	fmt.Printf("<== OnRecvSysMsg(%#x, %d, %d)\n", _handle, ctrl_id, sys_msg)
}

func OnRecvData(h syscall.Handle, ctrl_id uint, qry_name *string, param *string) {
	fmt.Printf("<-- OnRecvData(%#x, %d, \"%s\", \"%s\"): rq_id = %d\n", h, ctrl_id, *qry_name, *param, efx.GetRecvRqID(h))
	fmt.Printf("ReqMsgCode: %s, ReqMsg: \"%s\", RtCode: %s\n", efx.GetReqMsgCode(h), efx.GetReqMessage(h), efx.GetRtCode(h))

	switch *qry_name {
	case "SCP":
		str_mk := efx.GetSingleData(h, 2)
		cur_price, _ := efx.FormatNumber(efx.GetSingleDataInt(h, 11))
		mst_ptr := efx.FindStockMaster(*param)
		if mst_ptr != nil {
			mst := (efx.StockMasterInfo)(mst_ptr)
			fmt.Printf("%s %s %s 현재가: %s원 %+.2f%%\n", mst.StockCode(), mst.StockName(), str_mk, cur_price, efx.GetSingleDataFloat(h, 14))
			if mst.IsSuspension() {
				fmt.Println("******** 거래정지 ********")
			}
			fmt.Printf("seq_no: %d, 신용가능: %v 액면가: %d, 시장경고: %v, ROE: %.2f 상장날짜: %s\n", mst.SeqNo(), mst.IsCreditable(),
				mst.ParValue(), mst.MarketWarnCode(), mst.ROE(), mst.ListingDate())
			fmt.Printf("상장주식수: %d, 신규상장: %v, 자본금: %d, 매출액: %d, 기준날짜: %s\n", mst.ListingCount(), mst.IsNewListing(), mst.Capital(),
				mst.Sales(), mst.BaseDate())
			fmt.Printf("증거금비율: %d, 신용한도초과: %v, 담보대출가능: %v, 대주가능: %v\n", mst.MarginRate(), mst.IsCreditLimtOver(), mst.IsLoanable(), mst.IsStockLoanable())
			//fmt.Printf("%v\n", mst_ptr)
		}
	case "SCAP":
		cash, _ := efx.FormatNumber(efx.GetSingleDataInt(h, 0))
		max_buy_amount, _ := efx.FormatNumber(efx.GetSingleDataInt(h, 7))
		fmt.Printf("주문가능현금: %s원\t최대매수금액: %s원\n", cash, max_buy_amount)
	default:
		fmt.Println("엥?")
	}
	fmt.Println()
	_sig_recv_data = true
}

func OnRecvRealData(h syscall.Handle, ctrl_id uint, qry_name *string, param *string) {
	fmt.Printf("<-- OnRecvRealData(%#x, %d, \"%s\")\n", h, ctrl_id, *qry_name)

	stock_code := efx.GetSingleData(h, 0)
	cur_price, _ := efx.FormatNumber(efx.GetSingleDataInt(h, 2))
	mst_ptr := efx.FindStockMaster(stock_code)
	mst := (efx.StockMasterInfo)(mst_ptr)
	fmt.Printf("[%s] %s %s 현재가: %s원  %+.2f%%\n", efx.GetSingleData(h, 1), stock_code, mst.StockName(), cur_price, efx.GetSingleDataFloat(h, 5))
}

func OnRecvError(_handle syscall.Handle, ctrl_id uint, qry_name *string, param *string) {
	fmt.Printf("<-- OnRecvError(%#x, %d, \"%s\", \"%s\")\n", _handle, ctrl_id, *qry_name, *param)
}

func wait_recv_data() {
	for !_sig_recv_data {
		efx.ProcessMessage(1)
	}
	_sig_recv_data = false
}

func init() {
	str := efx.ExecCmd("ping -n 1 192.168.0.1", ^uint32(0))
	fmt.Printf("ExecCmd: %s\n", str)
	//efx.LaunchApp("C:\\eFriend Expert\\efriendexpert\\efexpertviewer.exe", "", false, false, 5)

	start := time.Now()
	e1 := efx.Init("", OnRecvSysMsg)
	if e1 != nil {
		fmt.Printf("efx.Init: %s\n", e1)
		efx.Exit()
		os.Exit(1)
	}
	fmt.Printf("efxInit() elapsed time: %v\n", time.Since(start))

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setConsoleCtrlHandler := kernel32.NewProc("SetConsoleCtrlHandler")
	r1, r2, lastErr := setConsoleCtrlHandler.Call(
		syscall.NewCallback(func(controlType uint) uint {
			fmt.Printf("consoleControlHandler called with %v\n", controlType)
			efx.Quit(0)
			return 1
		}), 1)
	fmt.Printf("call result %v %v %v\n\n", r1, r2, lastErr)
}

func main() {
	acc_cnt := efx.GetAccountCount()
	fmt.Printf("AccountCount: %d\n", acc_cnt)
	for i := 0; i < acc_cnt; i++ {
		acnt_no := efx.GetAccount(i)
		fmt.Printf("Account[%d]: %s, %s\n", i, acnt_no, efx.GetAccountBrcode(acnt_no))
	}
	efx1 := efx.NewControl(1, OnRecvData, OnRecvRealData, OnRecvError)
	if efx1 == 0 {
		panic("efx.NewControl error!")
	}

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
			fmt.Printf("Synchonize() elapsed time: %v, sync time: %dms\n", time.Since(wait_time), sync_time)
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
