package efxmfc

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

var (
	efxdll, _                   = syscall.LoadLibrary("EfxMfc.dll")
	efxInit, _                  = syscall.GetProcAddress(efxdll, "efxInit")
	efxExit, _                  = syscall.GetProcAddress(efxdll, "efxExit")
	efxSynchronize, _           = syscall.GetProcAddress(efxdll, "efxSynchronize")
	efxLaunchApp, _             = syscall.GetProcAddress(efxdll, "efxLaunchApp")
	efxExecCmd, _               = syscall.GetProcAddress(efxdll, "efxExecCmd")
	efxProcessMessage, _        = syscall.GetProcAddress(efxdll, "efxProcessMessage")
	efxMessageLoop, _           = syscall.GetProcAddress(efxdll, "efxMessageLoop")
	efxQuit, _                  = syscall.GetProcAddress(efxdll, "efxQuit")
	efxFindStockMaster, _       = syscall.GetProcAddress(efxdll, "efxFindStockMaster")
	efxFindStockMasterByName, _ = syscall.GetProcAddress(efxdll, "efxFindStockMasterByName")
	efxGetStockMaster, _        = syscall.GetProcAddress(efxdll, "efxGetStockMaster")

	efxNewControl, _         = syscall.GetProcAddress(efxdll, "efxNewControl")
	efxGetAccountCount, _    = syscall.GetProcAddress(efxdll, "efxGetAccountCount")
	efxGetAccount, _         = syscall.GetProcAddress(efxdll, "efxGetAccount")
	efxGetAccountBrcode, _   = syscall.GetProcAddress(efxdll, "efxGetAccountBrcode")
	efxGetEncryptPassword, _ = syscall.GetProcAddress(efxdll, "efxGetEncryptPassword")
	efxAboutBox, _           = syscall.GetProcAddress(efxdll, "efxAboutBox")

	efxSetSingleData, _     = syscall.GetProcAddress(efxdll, "efxSetSingleData")
	efxSetSingleDataEx, _   = syscall.GetProcAddress(efxdll, "efxSetSingleDataEx")
	efxSetMultiData, _      = syscall.GetProcAddress(efxdll, "efxSetMultiData")
	efxSetMultiBlockData, _ = syscall.GetProcAddress(efxdll, "efxSetMultiBlockData")

	efxRequestData, _          = syscall.GetProcAddress(efxdll, "efxRequestData")
	efxIsMoreNextData, _       = syscall.GetProcAddress(efxdll, "efxIsMoreNextData")
	efxRequestNextData, _      = syscall.GetProcAddress(efxdll, "efxRequestNextData")
	efxRequestRealData, _      = syscall.GetProcAddress(efxdll, "efxRequestRealData")
	efxUnrequestRealData, _    = syscall.GetProcAddress(efxdll, "efxUnrequestRealData")
	efxUnrequestAllRealData, _ = syscall.GetProcAddress(efxdll, "efxUnrequestAllRealData")

	efxGetSingleFieldCount, _ = syscall.GetProcAddress(efxdll, "efxGetSingleFieldCount")
	efxGetSingleData, _       = syscall.GetProcAddress(efxdll, "efxGetSingleData")
	efxGetSingleDataInt, _    = syscall.GetProcAddress(efxdll, "efxGetSingleDataInt")
	efxGetSingleDataFloat, _  = syscall.GetProcAddress(efxdll, "efxGetSingleDataFloat2")

	efxGetMultiBlockCount, _  = syscall.GetProcAddress(efxdll, "efxGetMultiBlockCount")
	efxGetMultiRecordCount, _ = syscall.GetProcAddress(efxdll, "efxGetMultiRecordCount")
	efxGetMultiFieldCount, _  = syscall.GetProcAddress(efxdll, "efxGetMultiFieldCount")
	efxGetMultiData, _        = syscall.GetProcAddress(efxdll, "efxGetMultiData")
	efxGetMultiDataInt, _     = syscall.GetProcAddress(efxdll, "efxGetMultiDataInt")
	efxGetMultiDataFloat, _   = syscall.GetProcAddress(efxdll, "efxGetMultiDataFloat2")

	efxGetReqMsgCode, _ = syscall.GetProcAddress(efxdll, "efxGetReqMsgCode")
	efxGetReqMessage, _ = syscall.GetProcAddress(efxdll, "efxGetReqMessage")
	efxGetRtCode, _     = syscall.GetProcAddress(efxdll, "efxGetRtCode")
	efxGetSendRqID, _   = syscall.GetProcAddress(efxdll, "efxGetSendRqID")
	efxGetRecvRqID, _   = syscall.GetProcAddress(efxdll, "efxGetRecvRqID")
)

const REQ_LIMIT_MS uint = 10 // 초당 조회TR 및 주문TR 제한: 20ms 간 5개 TR 이하 호출

type MarketType int32

const (
	KOSPI MarketType = iota
	KOSDAQ
)

type MarketWarning int32

const (
	CODE_NONE MarketWarning = iota
	CODE_CAUTION
	CODE_WARNING
	CODE_DANGER
)

type InvestInfo struct {
	is_suspension         bool          // 거래정지
	is_clearance_sale     bool          // 정리매매
	is_management         bool          // 관리종목
	market_warn_code      MarketWarning // 시장 경고 코드
	is_market_warning     bool          // 경고 예고 여부
	is_insincerity_notice bool          // 불성실 공시 여부
	is_backdoor_listing   bool          // 우회상장 여부
	is_creditable         bool          // 신용주문 가능
	margin_rate           int16         // 증거금 비율
	par_value             int32         // 액면가
	listing_date          [9]uint8      // 상장 날짜
	listing_count         int64         // 상장 주수
	is_new_listing        bool          // 신규상장 여부
	_dummy                [4]uint8      // C구조체와 정렬 차이 때려맞춤!
	capital               int64         // 자본금
	is_short_selling      bool          // 공매도주문가능 여부
	is_abnormal_rise      bool          // 이상급등종목 여부
	sales                 int32         // 매출액
	operating_profit      int32         // 영업이익
	ordinary_profit       int32         // 경상이익
	net_income            int32         // 당기순이익
	ROE                   float32       // 자기자본이익률
	base_year_month       [9]uint8      // 기준년월
	preday_market_cap     uint32        // 전일기준 시가총액(단위: 억원)
	is_credit_limt_over   bool          // 회사신용한도초과 여부
	is_loanable           bool          // 담보대출가능 여부
	is_stock_loanable     bool          // 대주가능 여부
}

type StockMasterItem struct {
	seq_no          int32
	stock_code      [10]uint8
	standard_code   [13]uint8
	stock_name      [61]uint8
	market_type     MarketType
	is_venture      bool
	is_kospi50      bool
	is_kospi100     bool
	kospi200_sector int8 // KOSPI200 섹터업종
	is_krx100       bool
	is_krx300       bool
	is_etf          bool
	is_acquisition  bool     // 기업인수목적회사 여부
	_dummy          [4]uint8 // C구조체와 정렬 차이 때려맞춤!
	invest_info     InvestInfo
}

type StockMasterInfo interface {
	SeqNo() int
	StockCode() string
	StockName() string
	MarketType() MarketType
	IsETF() bool
	IsSuspension() bool            // 거래정지
	IsClearanceSale() bool         // 정리매매
	IsManagement() bool            // 관리종목
	MarketWarnCode() MarketWarning // 시장 경고 코드
	IsMarketWarning() bool         // 경고 예고 여부
	IsInsincerityNotice() bool     // 불성실 공시 여부
	IsBackdoorListing() bool       // 우회상장 여부
	IsCreditable() bool            // 신용주문 가능
	MarginRate() int               // 증거금 비율
	ParValue() int                 // 액면가
	ListingDate() string           // 상장 날짜
	ListingCount() int64
	IsNewListing() bool
	Capital() int64
	Sales() int // 매출액
	ROE() float32
	BaseDate() string // 기준년월
	IsCreditLimtOver() bool
	IsLoanable() bool      // 담보대출가능 여부
	IsStockLoanable() bool // 대주가능 여부
}

func (smi StockMasterItem) SeqNo() int {
	return int(smi.seq_no)
}

func (smi StockMasterItem) StockCode() string {
	return PtrToString(unsafe.Pointer(&smi.stock_code[0]), 10)
}

func (smi StockMasterItem) StockName() string {
	return PtrToString(unsafe.Pointer(&smi.stock_name[0]), 61)
}

func (smi StockMasterItem) MarketType() MarketType {
	return smi.market_type
}

func (smi StockMasterItem) IsETF() bool {
	return smi.is_etf
}

func (smi StockMasterItem) IsSuspension() bool {
	return smi.invest_info.is_suspension
}

func (smi StockMasterItem) IsClearanceSale() bool {
	return smi.invest_info.is_clearance_sale
}

func (smi StockMasterItem) IsManagement() bool {
	return smi.invest_info.is_management
}

func (smi StockMasterItem) MarketWarnCode() MarketWarning {
	return smi.invest_info.market_warn_code
}

func (smi StockMasterItem) IsMarketWarning() bool {
	return smi.invest_info.is_market_warning
}

func (smi StockMasterItem) IsInsincerityNotice() bool {
	return smi.invest_info.is_insincerity_notice
}

func (smi StockMasterItem) IsBackdoorListing() bool {
	return smi.invest_info.is_backdoor_listing
}

func (smi StockMasterItem) IsCreditable() bool {
	return smi.invest_info.is_creditable
}

func (smi StockMasterItem) MarginRate() int {
	return int(smi.invest_info.margin_rate)
}

func (smi StockMasterItem) ParValue() int {
	return int(smi.invest_info.par_value)
}

func (smi StockMasterItem) ListingDate() string {
	return PtrToString(unsafe.Pointer(&smi.invest_info.listing_date[0]), 9)
}

func (smi StockMasterItem) ListingCount() int64 {
	return smi.invest_info.listing_count
}

func (smi StockMasterItem) IsNewListing() bool {
	return smi.invest_info.is_new_listing
}

func (smi StockMasterItem) Capital() int64 {
	return smi.invest_info.capital
}

func (smi StockMasterItem) Sales() int {
	return int(smi.invest_info.sales)
}

func (smi StockMasterItem) ROE() float32 {
	return smi.invest_info.ROE
}

func (smi StockMasterItem) BaseDate() string {
	return PtrToString(unsafe.Pointer(&smi.invest_info.base_year_month[0]), 9)
}

func (smi StockMasterItem) IsCreditLimtOver() bool {
	return smi.invest_info.is_credit_limt_over
}

func (smi StockMasterItem) IsLoanable() bool {
	return smi.invest_info.is_loanable
}

func (smi StockMasterItem) IsStockLoanable() bool {
	return smi.invest_info.is_stock_loanable
}

type EventCallback func(syscall.Handle, uint, *string, *string)
type SysMsgCallback func(syscall.Handle, uint, int)

func Init(efx_dir string, sysmsg_cb SysMsgCallback) (err error) {
	cb := syscall.NewCallback(func(h syscall.Handle, id uint, msg int) uintptr {
		sysmsg_cb(h, id, msg)
		return 0
	})
	var dir_ptr unsafe.Pointer = nil
	if len(efx_dir) > 0 {
		str := append([]byte(efx_dir), 0)
		dir_ptr = unsafe.Pointer(&str[0])
	}
	r1, _, e1 := syscall.SyscallN(efxInit, uintptr(dir_ptr), cb)
	if r1 != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
			switch int(r1) {
			case -1:
				err = errors.New("administrative privileges required")
			case 1:
				err = errors.New("incorrect function")
			case 2:
				err = errors.New("file not found")
			case 3:
				err = errors.New("path not found")
			default:
				err = errors.New("undefined")
			}
			return err
		}
	}
	return
}

func Exit() {
	_, _, e1 := syscall.SyscallN(efxExit, 0, 0)
	if e1 != 0 {
		abort("Call efxExit", error(e1))
	}
}

func Synchonize(millisec uint) uint32 {
	r1, _, _ := syscall.SyscallN(efxSynchronize, uintptr(millisec), 0)
	return uint32(r1)
}

func LaunchApp(file_path string, args string, waitInputIdle bool, waitTerminate bool, showFlag int32) syscall.Handle {
	str_path := append([]byte(file_path), 0)
	path_ptr := unsafe.Pointer(&str_path[0])
	str_arg := append([]byte(args), 0)
	arg_ptr := unsafe.Pointer(&str_arg[0])
	var wait_input_idle, wait_terminate uint8 = 0, 0
	if waitInputIdle {
		wait_input_idle = 1
	}
	if waitTerminate {
		wait_terminate = 1
	}
	r1, _, _ := syscall.SyscallN(efxLaunchApp, uintptr(path_ptr), uintptr(arg_ptr), uintptr(wait_input_idle), uintptr(wait_terminate), uintptr(showFlag))
	return syscall.Handle(r1)
}

func ExecCmd(cmd_arg string, wait_time uint32) string {
	str_cmd := append([]byte(cmd_arg), 0)
	cmd_ptr := unsafe.Pointer(&str_cmd[0])
	data := make([]byte, 4096)
	r1, _, _ := syscall.SyscallN(efxExecCmd, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(cmd_ptr), uintptr(wait_time))
	return string(data[:r1])
}

func ProcessMessage(millisec uint) int {
	r1, _, _ := syscall.SyscallN(efxProcessMessage, uintptr(millisec), 0)
	return int(r1)
}

func MessageLoop() {
	_, _, e1 := syscall.SyscallN(efxMessageLoop, 0, 0)
	if e1 != 0 {
		abort("Call efxMessageLoop", error(e1))
	}
}

func Quit(exit_code int) {
	syscall.SyscallN(efxQuit, 0, 0)
}

func AboutBox() {
	syscall.SyscallN(efxAboutBox, 0, 0)
}

func PtrToString(p unsafe.Pointer, max_len int) string {
	if p == nil {
		return ""
	}
	arr := make([]uint8, max_len)
	str_len := 0
	for i := 0; i < max_len; i++ {
		if c := *(*uint8)(unsafe.Add(p, uintptr(i))); c != 0 {
			arr[i] = c
		} else {
			str_len = i
			break
		}
	}
	return string(arr[:str_len])
}

func FormatNumber(N interface{}) (string, error) {
	switch N.(type) {
	case int:
	case int16:
	case int32:
	case int64:
	case int8:
	case float32:
	case float64:
	default:
		return "", fmt.Errorf("Not a valid number!")
	}

	n := fmt.Sprintf("%v", N)
	n = strings.ReplaceAll(n, ",", "")
	dec := ""
	if strings.Index(n, ".") != -1 {
		dec = n[strings.Index(n, ".")+1 : len(n)]
		n = n[0:strings.Index(n, ".")]
	}
	for i := 0; i <= len(n); i = i + 4 {
		a := n[0 : len(n)-i]
		b := n[len(n)-i : len(n)]
		n = a + "," + b
	}
	if n[0:1] == "," {
		n = n[1:len(n)]
	}
	if n[len(n)-1:len(n)] == "," {
		n = n[0 : len(n)-1]
	}
	if dec != "" {
		n = n + "." + dec
	}
	return n, nil
}

func FindStockMaster(stock_code string) *StockMasterItem {
	str := append([]byte(stock_code), 0)
	code_ptr := unsafe.Pointer(&str[0])
	r1, _, _ := syscall.SyscallN(efxFindStockMaster, uintptr(code_ptr), 0, 0)
	var pItem *StockMasterItem = nil
	if r1 != 0 {
		pItem = (*StockMasterItem)(unsafe.Pointer(r1))
	}
	return pItem
}

func FindStockMasterByName(stock_name string) *StockMasterItem {
	str := append([]byte(stock_name), 0)
	name_ptr := unsafe.Pointer(&str[0])
	r1, _, _ := syscall.SyscallN(efxFindStockMasterByName, uintptr(name_ptr), 0, 0)
	var pItem *StockMasterItem = nil
	if r1 != 0 {
		pItem = (*StockMasterItem)(unsafe.Pointer(r1))
	}
	return pItem
}

func GetStockMasterItem(seq_no int32) *StockMasterItem {
	r1, _, _ := syscall.SyscallN(efxGetStockMaster, uintptr(seq_no), 0, 0)
	var pItem *StockMasterItem = nil
	if r1 != 0 {
		pItem = (*StockMasterItem)(unsafe.Pointer(r1))
	}
	return pItem
}

func NewControl(id uint, recv_cb EventCallback, real_cb EventCallback, error_cb EventCallback) syscall.Handle {
	var cb1, cb2, cb3 uintptr = 0, 0, 0
	if recv_cb != nil {
		cb1 = syscall.NewCallback(func(h syscall.Handle, id uint, qry unsafe.Pointer, param unsafe.Pointer) uintptr {
			qry_name := PtrToString(qry, 16)
			str_param := PtrToString(param, 256)
			recv_cb(h, id, &qry_name, &str_param)
			return 0
		})
	}
	if real_cb != nil {
		cb2 = syscall.NewCallback(func(h syscall.Handle, id uint, qry unsafe.Pointer, param unsafe.Pointer) uintptr {
			qry_name := PtrToString(qry, 16)
			real_cb(h, id, &qry_name, nil)
			return 0
		})
	}
	if error_cb != nil {
		cb3 = syscall.NewCallback(func(h syscall.Handle, id uint, qry unsafe.Pointer, param unsafe.Pointer) uintptr {
			qry_name := PtrToString(qry, 16)
			str_param := PtrToString(param, 256)
			error_cb(h, id, &qry_name, &str_param)
			return 0
		})
	}
	r1, _, _ := syscall.SyscallN(efxNewControl, uintptr(id), cb1, cb2, cb3)
	return syscall.Handle(r1)
}

func GetAccountCount() int {
	r1, _, _ := syscall.SyscallN(efxGetAccountCount, 0, 0, 0)
	return int(r1)
}

func GetAccount(idx int) string {
	data := make([]byte, 16)
	r1, _, _ := syscall.SyscallN(efxGetAccount, uintptr(idx), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func GetAccountBrcode(acnt_no string) string {
	str := append([]byte(acnt_no), 0)
	acnt_ptr := unsafe.Pointer(&str[0])
	data := make([]byte, 16)
	r1, _, _ := syscall.SyscallN(efxGetAccountBrcode, uintptr(acnt_ptr), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func GetEncryptPassword(pwd string) string {
	str := append([]byte(pwd), 0)
	pwd_ptr := unsafe.Pointer(&str[0])
	data := make([]byte, 100)
	r1, _, _ := syscall.SyscallN(efxGetEncryptPassword, uintptr(pwd_ptr), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func SetSingleData(h syscall.Handle, field_idx int, data string) {
	str := append([]byte(data), 0)
	data_ptr := unsafe.Pointer(&str[0])
	syscall.SyscallN(efxSetSingleData, uintptr(h), uintptr(field_idx), uintptr(data_ptr))
}

func SetSingleDataEx(h syscall.Handle, block_idx int, field_idx int, data string) {
	str := append([]byte(data), 0)
	data_ptr := unsafe.Pointer(&str[0])
	syscall.SyscallN(efxSetSingleDataEx, uintptr(h), uintptr(block_idx), uintptr(field_idx), uintptr(data_ptr))
}

func SetMultiData(h syscall.Handle, rec_idx int, field_idx int, data string) {
	str := append([]byte(data), 0)
	data_ptr := unsafe.Pointer(&str[0])
	syscall.SyscallN(efxSetMultiData, uintptr(h), uintptr(rec_idx), uintptr(field_idx), uintptr(data_ptr))
}

func SetMultiBlockData(h syscall.Handle, block_idx int, rec_idx int, field_idx int, data string) {
	str := append([]byte(data), 0)
	data_ptr := unsafe.Pointer(&str[0])
	syscall.SyscallN(efxSetMultiBlockData, uintptr(h), uintptr(block_idx), uintptr(rec_idx), uintptr(field_idx), uintptr(data_ptr))
}

func RequestData(h syscall.Handle, qry_name string, param string) {
	str := append([]byte(qry_name), 0)
	qry_ptr := unsafe.Pointer(&str[0])
	str2 := append([]byte(param), 0)
	param_ptr := unsafe.Pointer(&str2[0])
	syscall.SyscallN(efxRequestData, uintptr(h), uintptr(qry_ptr), uintptr(param_ptr))
}

func IsMoreNextData(h syscall.Handle) bool {
	r1, _, _ := syscall.SyscallN(efxIsMoreNextData, uintptr(h), 0, 0, 0)
	return r1 > 0
}

func RequestNextData(h syscall.Handle, qry_name string) {
	str := append([]byte(qry_name), 0)
	qry_ptr := unsafe.Pointer(&str[0])
	syscall.SyscallN(efxRequestNextData, uintptr(h), uintptr(qry_ptr), 0)
}

func RequestRealData(h syscall.Handle, qry_name string, code string) {
	str := append([]byte(qry_name), 0)
	qry_ptr := unsafe.Pointer(&str[0])
	str2 := append([]byte(code), 0)
	code_ptr := unsafe.Pointer(&str2[0])
	syscall.SyscallN(efxRequestRealData, uintptr(h), uintptr(qry_ptr), uintptr(code_ptr))
}

func UnrequestRealData(h syscall.Handle, qry_name string, code string) {
	str := append([]byte(qry_name), 0)
	qry_ptr := unsafe.Pointer(&str[0])
	str2 := append([]byte(code), 0)
	code_ptr := unsafe.Pointer(&str2[0])
	syscall.SyscallN(efxUnrequestRealData, uintptr(h), uintptr(qry_ptr), uintptr(code_ptr))
}

func UnrequestAllRealData(h syscall.Handle) {
	syscall.SyscallN(efxUnrequestAllRealData, uintptr(h), 0, 0)
}

func GetSingleFieldCount(h syscall.Handle) int {
	r1, _, _ := syscall.SyscallN(efxGetSingleFieldCount, uintptr(h), 0, 0, 0)
	return int(r1)
}

func GetSingleData(h syscall.Handle, n ...int) string {
	data := make([]byte, 256)
	field_idx, attr_type := 0, 0
	args := len(n)
	if args > 0 {
		field_idx = n[0]
		if args > 1 {
			attr_type = n[1]
		}
	}
	r1, _, _ := syscall.SyscallN(efxGetSingleData, uintptr(h), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(field_idx), uintptr(attr_type))
	return string(data[:r1])
}

func GetSingleDataInt(h syscall.Handle, n ...int) int {
	field_idx, attr_type := 0, 0
	args := len(n)
	if args > 0 {
		field_idx = n[0]
		if args > 1 {
			attr_type = n[1]
		}
	}
	r1, _, _ := syscall.SyscallN(efxGetSingleDataInt, uintptr(h), uintptr(field_idx), uintptr(attr_type))
	return int(r1)
}

func GetSingleDataFloat(h syscall.Handle, n ...int) float32 {
	field_idx, attr_type := 0, 0
	args := len(n)
	if args > 0 {
		field_idx = n[0]
		if args > 1 {
			attr_type = n[1]
		}
	}
	var result float32 = 0
	syscall.SyscallN(efxGetSingleDataFloat, uintptr(h), uintptr(unsafe.Pointer(&result)), uintptr(field_idx), uintptr(attr_type))
	return result
}

func GetMultiBlockCount(h syscall.Handle) int {
	r1, _, _ := syscall.SyscallN(efxGetMultiBlockCount, uintptr(h), 0, 0, 0)
	return int(r1)
}

func GetMultiRecordCount(h syscall.Handle, block_idx int) int {
	r1, _, _ := syscall.SyscallN(efxGetMultiRecordCount, uintptr(h), uintptr(block_idx), 0, 0)
	return int(r1)
}

func GetMultiFieldCount(h syscall.Handle, block_idx int, rec_idx int) int {
	r1, _, _ := syscall.SyscallN(efxGetMultiFieldCount, uintptr(h), uintptr(block_idx), uintptr(rec_idx), 0)
	return int(r1)
}

func GetMultiData(h syscall.Handle, n ...int) string {
	data := make([]byte, 256)
	var block_idx, rec_idx, field_idx, attr_type int = 0, 0, 0, 0
	args := len(n)
	if args > 0 {
		block_idx = n[0]
		if args > 1 {
			rec_idx = n[1]
			if args > 2 {
				field_idx = n[2]
				if args > 3 {
					attr_type = n[3]
				}
			}
		}
	}
	r1, _, _ := syscall.SyscallN(efxGetMultiData, uintptr(h), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), uintptr(block_idx), uintptr(rec_idx), uintptr(field_idx), uintptr(attr_type))
	return string(data[:r1])
}

func GetMultiDataInt(h syscall.Handle, n ...int) int {
	var block_idx, rec_idx, field_idx, attr_type int = 0, 0, 0, 0
	args := len(n)
	if args > 0 {
		block_idx = n[0]
		if args > 1 {
			rec_idx = n[1]
			if args > 2 {
				field_idx = n[2]
				if args > 3 {
					attr_type = n[3]
				}
			}
		}
	}
	r1, _, _ := syscall.SyscallN(efxGetMultiDataInt, uintptr(h), uintptr(block_idx), uintptr(rec_idx), uintptr(field_idx), uintptr(attr_type))
	return int(r1)
}

func GetMultiDataFloat(h syscall.Handle, n ...int) float32 {
	var block_idx, rec_idx, field_idx, attr_type int = 0, 0, 0, 0
	args := len(n)
	if args > 0 {
		block_idx = n[0]
		if args > 1 {
			rec_idx = n[1]
			if args > 2 {
				field_idx = n[2]
				if args > 3 {
					attr_type = n[3]
				}
			}
		}
	}
	var result float32 = 0
	syscall.SyscallN(efxGetMultiDataFloat, uintptr(h), uintptr(unsafe.Pointer(&result)), uintptr(block_idx), uintptr(rec_idx), uintptr(field_idx), uintptr(attr_type))
	return result
}

func GetReqMsgCode(h syscall.Handle) string {
	data := make([]byte, 100)
	r1, _, _ := syscall.SyscallN(efxGetReqMsgCode, uintptr(h), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func GetReqMessage(h syscall.Handle) string {
	data := make([]byte, 256)
	r1, _, _ := syscall.SyscallN(efxGetReqMessage, uintptr(h), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func GetRtCode(h syscall.Handle) string {
	data := make([]byte, 100)
	r1, _, _ := syscall.SyscallN(efxGetRtCode, uintptr(h), uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)))
	return string(data[:r1])
}

func GetSendRqID(h syscall.Handle) int {
	r1, _, _ := syscall.SyscallN(efxGetSendRqID, uintptr(h), 0, 0)
	return int(r1)
}

func GetRecvRqID(h syscall.Handle) int {
	r1, _, _ := syscall.SyscallN(efxGetRecvRqID, uintptr(h), 0, 0)
	return int(r1)
}

func RequestAccountBalance(h syscall.Handle, account string, param ...string) {
	SetSingleData(h, 0, account[:8])
	SetSingleData(h, 1, account[8:10])
	enc_pwd := GetEncryptPassword(param[0])
	SetSingleData(h, 2, enc_pwd)
	SetSingleData(h, 3, "")
	SetSingleData(h, 4, "")
	SetSingleData(h, 5, "00")
	SetSingleData(h, 6, "N")
	if len(param) > 1 {
		RequestData(h, "SCAP", param[1])
	} else {
		RequestData(h, "SCAP", account)
	}
}
