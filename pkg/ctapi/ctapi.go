package ctapi

import (
	"golang.org/x/sys/windows"
)

const CT_OPEN_RECONNECT = 0x02
const CT_OPEN_READ_ONLY = 0x04
const CT_OPEN_BATCH = 0x08
const CT_OPEN_EXTENDED = 0x10
const CT_OPEN_WINDOWSUSER = 0x20

const DBTYPE_EMPTY = 0
const DBTYPE_NULL = 1
const DBTYPE_I2 = 2
const DBTYPE_I4 = 3
const DBTYPE_R4 = 4
const DBTYPE_R8 = 5
const DBTYPE_CY = 6
const DBTYPE_DATE = 7
const DBTYPE_BSTR = 8
const DBTYPE_IDISPATCH = 9
const DBTYPE_ERROR = 10
const DBTYPE_BOOL = 11
const DBTYPE_VARIANT = 12
const DBTYPE_IUNKNOWN = 13
const DBTYPE_DECIMAL = 14
const DBTYPE_UI1 = 17
const DBTYPE_ARRAY = 0x2000
const DBTYPE_BYREF = 0x4000
const DBTYPE_I1 = 16
const DBTYPE_UI2 = 18
const DBTYPE_UI4 = 19
const DBTYPE_I8 = 20
const DBTYPE_UI8 = 21
const DBTYPE_GUID = 72
const DBTYPE_VECTOR = 0x1000
const DBTYPE_RESERVED = 0x8000
const DBTYPE_BYTES = 128
const DBTYPE_STR = 129
const DBTYPE_WSTR = 130
const DBTYPE_NUMERIC = 131
const DBTYPE_UDT = 132
const DBTYPE_DBDATE = 133
const DBTYPE_DBTIME = 134
const DBTYPE_DBTIMESTAMP = 135

var nullptr = uintptr(0)

type CtApi struct {
	dll   *windows.DLL
	procs *ctApiProcs
}
type ctApiProcs struct {
	//AlmQuery *windows.Proc
	//CtAPIAlarm *windows.Proc
	//CtAPITrend *windows.Proc
	//TrnQuery                   *windows.Proc
	ctCancelIO *windows.Proc
	//ctCiCode                   *windows.Proc
	ctClientCreate        *windows.Proc
	ctClientDestroy       *windows.Proc
	ctClose               *windows.Proc
	ctCloseEx             *windows.Proc
	ctEngToRaw            *windows.Proc
	ctFindClose           *windows.Proc
	ctFindFirst           *windows.Proc
	ctFindFirstEx         *windows.Proc
	ctFindNext            *windows.Proc
	ctFindNumRecords      *windows.Proc
	ctFindPrev            *windows.Proc
	ctFindScroll          *windows.Proc
	ctGetOverlappedResult *windows.Proc
	ctGetProperty         *windows.Proc
	//ctHasOverlappedIoCompleted *windows.Proc
	ctListNew                *windows.Proc
	ctListAdd                *windows.Proc
	ctListAddEx              *windows.Proc
	ctListData               *windows.Proc
	ctListDelete             *windows.Proc
	ctListEvent              *windows.Proc
	ctListFree               *windows.Proc
	ctListItem               *windows.Proc
	ctListRead               *windows.Proc
	ctListWrite              *windows.Proc
	ctOpen                   *windows.Proc
	ctOpenEx                 *windows.Proc
	ctRawToEng               *windows.Proc
	ctSetManagedBinDirectory *windows.Proc
	ctTagGetProperty         *windows.Proc
	ctTagRead                *windows.Proc
	ctTagReadEx              *windows.Proc
	ctTagWrite               *windows.Proc
	ctTagWriteEx             *windows.Proc
}
