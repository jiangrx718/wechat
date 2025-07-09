package common

const DEFAULT_PAGE_SIZE = 12

const (
	Math_Add_Symbol_Show = "➕"
	Math_Sub_Symbol_Show = "➖"
	Math_Mul_Symbol_Show = "✖️"
	Math_Div_Symbol_Show = "➗"
)

const (
	Math_Add_Non_Carry = 1 // 加法不进位
	Math_Add_Carry     = 2 // 加法进位
	Math_Sub_Non_Back  = 3 // 不退位减法
	Math_Sub_Back      = 4 // 退位减法
)

const (
	Math_Add_WithinFive             = 1  // 五以内加法
	Math_Sub_WithinTen              = 2  // 十以内减法
	Math_Add_Sub_WithinTen          = 3  // 十以内加减法
	Math_Add_WithinTwenty_Non_Carry = 4  //20以内加法(不进位)
	Math_Sub_WithinTwenty_Non_Back  = 5  //20以内减法(不退位)
	Math_Add_WithinTwenty_Carry     = 6  //20以内加法(进位)
	Math_Sub_WithinTwenty_Back      = 7  //20以内减法(退位)
	Math_Add_Sub_WithinTwenty       = 8  // 二十以内加减法
	Math_Add_WithinHundred          = 9  //100以内加法
	Math_Sub_WithinHundred          = 10 //100以内减法
	Math_Add_Sub_WithinHundred      = 11 //100以内加减法
)
