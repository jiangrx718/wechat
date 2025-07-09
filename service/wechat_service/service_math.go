package wechat_service

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jiangrx816/wechat/common"
	"github.com/jiangrx816/wechat/common/response"
	"github.com/jiangrx816/wechat/core/server/api"
)

/**
 * ApiServiceMathCalculationList
 * @Description 获取计算题列表
 */
func (ps *WechatService) ApiServiceMathCalculationList(ctx *gin.Context, forward, value, limit int) (resp []response.MathResponse, apiErr api.Error) {
	var initData [][2]int
	switch forward {
	case common.Math_Add_WithinFive: // 五以内加法
		initData = ps.generateAdditionWithinFive(limit, value)
	case common.Math_Sub_WithinTen: // 十以内减法
		initData = ps.generateSubtractionWithinTen(limit, value)
	case common.Math_Add_Sub_WithinTen: // 十以内加减法
		initData = ps.generateSubtractionWithinTen(limit, value)
	case common.Math_Add_WithinTwenty_Non_Carry: //20以内加法(不进位)
		initData = ps.generateStrictNonCarryAdditions(limit, value)
	case common.Math_Sub_WithinTwenty_Non_Back: //20以内减法(不退位)
		initData = ps.generateNonBorrowSubtraction(limit, value)
	case common.Math_Add_WithinTwenty_Carry: //20以内加法(进位)
		initData = ps.generateCarryAdditionProblems(limit, value)
	case common.Math_Sub_WithinTwenty_Back: //20以内减法(退位)
		initData = ps.generateBorrowSubtraction(limit, value)
	case common.Math_Add_Sub_WithinTwenty: // 二十以内加减法
		initData = ps.generateBorrowSubtraction(limit, value)
	case common.Math_Add_WithinHundred: //100以内加法
		initData = ps.generateAdditionWithinHundred(limit, value)
	case common.Math_Sub_WithinHundred: //100以内减法
		initData = ps.generateSubtractionWithinHundred(limit, value)
	case common.Math_Add_Sub_WithinHundred: //100以内加减法
		initData = ps.generateBorrowSubtraction(limit, value)
	default:
		initData = ps.generateAdditionWithinFive(limit, value)
	}

	resp = ps.FormatMathCalculationData(initData, forward, limit, value)

	return
}

// @Description 格式化计算题数据
func (ps *WechatService) FormatMathCalculationData(initData [][2]int, forward int, count int, value int) (resp []response.MathResponse) {
	var forwardSymbol string
	switch forward {
	case common.Math_Add_WithinFive, common.Math_Add_WithinTwenty_Non_Carry, common.Math_Add_WithinTwenty_Carry, common.Math_Add_WithinHundred:
		forwardSymbol = common.Math_Add_Symbol_Show
	case common.Math_Sub_WithinTen, common.Math_Sub_WithinTwenty_Non_Back, common.Math_Sub_WithinTwenty_Back, common.Math_Sub_WithinHundred:
		forwardSymbol = common.Math_Sub_Symbol_Show
	case common.Math_Add_Sub_WithinTen, common.Math_Add_Sub_WithinTwenty, common.Math_Add_Sub_WithinHundred:
		return ps.formatMathAddSubCalculationData(count, value)
	}

	for _, v := range initData {
		var forwardResult int
		if forwardSymbol == common.Math_Add_Symbol_Show {
			forwardResult = v[0] + v[1]
		}
		if forwardSymbol == common.Math_Sub_Symbol_Show {
			forwardResult = v[0] - v[1]
		}
		resp = append(resp, response.MathResponse{
			NumberOne: v[0],
			NumberTwo: v[1],
			Symbol:    forwardSymbol,
			Result:    forwardResult,
		})
	}
	return
}

// @Description 生成加减法的算式
func (ps *WechatService) formatMathAddSubCalculationData(count int, value int) (resp []response.MathResponse) {
	addSubDataList := ps.generateArithmeticWithinHundred(count, value)

	for _, v := range addSubDataList {
		var forwardResult int
		var forwardSymbol string
		if v[1] == "+" {
			num1, _ := strconv.Atoi(v[0])
			num2, _ := strconv.Atoi(v[2])
			forwardResult = num1 + num2
			forwardSymbol = common.Math_Add_Symbol_Show
			resp = append(resp, response.MathResponse{
				NumberOne: num1,
				NumberTwo: num2,
				Symbol:    forwardSymbol,
				Result:    forwardResult,
			})
		}
		if v[1] == "-" {
			num1, _ := strconv.Atoi(v[0])
			num2, _ := strconv.Atoi(v[2])
			forwardResult = num1 - num2
			forwardSymbol = common.Math_Sub_Symbol_Show
			resp = append(resp, response.MathResponse{
				NumberOne: num1,
				NumberTwo: num2,
				Symbol:    forwardSymbol,
				Result:    forwardResult,
			})
		}
	}
	return
}

// @Description 生成不超过5的加法题目
func (ps *WechatService) generateAdditionWithinFive(count int, max int) [][2]int {
	all := make([][2]int, 0)

	for a := 0; a <= max; a++ {
		for b := a; b <= max; b++ { // b >= a，避免重复
			if a+b <= max {
				all = append(all, [2]int{a, b})
			}
		}
	}

	if count > len(all) {
		count = len(all) // 避免死循环
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:count]
}

// @Description 生成10以内减法算式
func (ps *WechatService) generateSubtractionWithinTen(count int, max int) [][2]int {
	all := make([][2]int, 0)
	for a := 0; a <= max; a++ {
		for b := 0; b <= a; b++ {
			all = append(all, [2]int{a, b})
		}
	}

	if count > len(all) {
		count = len(all)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:count]
}

// @Description 不进位加法——生成严格不进位加法题目
func (ps *WechatService) generateStrictNonCarryAdditions(count int, max int) [][2]int {
	all := make([][2]int, 0)
	seen := make(map[string]bool)

	for a := 1; a <= max; a++ {
		for b := a; b <= max; b++ { // b >= a，避免交换律重复
			// 规则判断
			valid := false
			if max <= 20 {
				valid = a+b <= 20 && (a%10)+(b%10) < 10
			} else {
				valid = a >= 10 && b >= 10 &&
					a+b <= max &&
					(a%10)+(b%10) < 10
			}
			if valid {
				key := fmt.Sprintf("%d+%d", a, b)
				if !seen[key] {
					seen[key] = true
					all = append(all, [2]int{a, b})
				}
			}
		}
	}

	// 限制返回数量，避免超出范围
	if count > len(all) {
		count = len(all)
	}

	// 打乱并选取
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:count]
}

// @Description 进位加法-生成进位加法题目（严格限制和的范围）
func (ps *WechatService) generateCarryAdditionProblems(count int, max int) [][2]int {
	all := make([][2]int, 0)
	seen := make(map[string]bool)

	lowLimit := 1
	if max > 20 {
		lowLimit = 10
	}

	for a := lowLimit; a <= max; a++ {
		for b := a; b <= max; b++ { // b >= a 避免重复
			if (a%10)+(b%10) >= 10 { // 进位条件
				if max <= 20 && a+b <= 20 {
					key := fmt.Sprintf("%d+%d", a, b)
					if !seen[key] {
						seen[key] = true
						all = append(all, [2]int{a, b})
					}
				} else if max > 20 && a+b <= max {
					key := fmt.Sprintf("%d+%d", a, b)
					if !seen[key] {
						seen[key] = true
						all = append(all, [2]int{a, b})
					}
				}
			}
		}
	}

	if count > len(all) {
		count = len(all) // 防止请求超过可能题目数量
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:count]
}

// @Description 不退位减法——生成不退位减法题目
func (ps *WechatService) generateNonBorrowSubtraction(count int, max int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	problems := make([][2]int, 0, count)
	seen := make(map[string]bool) // 用于检测重复的映射

	for len(problems) < count {
		var a, b int

		if max <= 20 {
			// 20以内：允许个位数
			a = rand.Intn(max) + 1 // 1-20
			b = rand.Intn(max) + 1 // 1-20
		} else {
			// 超过20：不允许个位数（必须≥10）
			a = rand.Intn(max-10) + 10 // 10-max
			b = rand.Intn(max-10) + 10 // 10-max
		}

		// 创建唯一标识符（减法顺序很重要，不考虑交换律）
		key := fmt.Sprintf("%d-%d", a, b)

		// 检查是否已存在相同算式
		if seen[key] {
			continue
		}

		// 共同条件：被减数≥减数且个位≥减数个位（不退位）
		if a >= b && (a%10) >= (b%10) {
			// 20以内额外条件和≤20
			if max <= 20 && a <= 20 && b <= 20 {
				// 标记该算式为已见
				seen[key] = true
				// 添加到结果集
				problems = append(problems, [2]int{a, b})
			}
			// 超过20时，和不超过max
			if max > 20 && a <= max && b <= max {
				// 标记该算式为已见
				seen[key] = true
				// 添加到结果集
				problems = append(problems, [2]int{a, b})
			}
		}
	}
	return problems
}

// @Description 退位减法——生成退位减法题目
func (ps *WechatService) generateBorrowSubtraction(count int, max int) [][2]int {
	all := make([][2]int, 0)
	seen := make(map[string]bool)

	lowLimit := 1
	if max > 20 {
		lowLimit = 10
	}

	for a := lowLimit; a <= max; a++ {
		for b := lowLimit; b <= max; b++ {
			if a > b && (a%10) < (b%10) { // 需要借位的条件
				// 额外条件
				if max <= 20 && a <= 20 && b <= 20 {
					key := fmt.Sprintf("%d-%d", a, b)
					if !seen[key] {
						seen[key] = true
						all = append(all, [2]int{a, b})
					}
				} else if max > 20 && a <= max && b <= max {
					key := fmt.Sprintf("%d-%d", a, b)
					if !seen[key] {
						seen[key] = true
						all = append(all, [2]int{a, b})
					}
				}
			}
		}
	}

	if count > len(all) {
		count = len(all)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:count]
}

// @Description 生成100以内加减法算式
func (ps *WechatService) generateArithmeticWithinHundred(count int, max int) [][3]string {
	rand.Seed(time.Now().UnixNano())
	problems := make([][3]string, 0, count)
	seen := make(map[string]bool) // 用于检测重复的映射

	maxNextValue := max + 1
	for len(problems) < count {
		// 随机选择加法或减法（50%概率）
		isAddition := rand.Intn(2) == 0
		var a, b int
		var problem [3]string

		if isAddition {
			// 生成加法算式
			a = rand.Intn(maxNextValue) // 0-100
			// 确保 a+b ≤ 100
			maxB := max - a
			b = rand.Intn(maxB + 1) // 0 到 maxB

			// 创建唯一标识符（考虑加法交换律）
			key1 := fmt.Sprintf("%d+%d", a, b)
			key2 := fmt.Sprintf("%d+%d", b, a)

			// 检查是否已存在相同算式
			if seen[key1] || seen[key2] {
				continue
			}

			// 标记该算式和其交换形式为已见
			seen[key1] = true
			seen[key2] = true

			problem = [3]string{fmt.Sprintf("%d", a), "+", fmt.Sprintf("%d", b)}
		} else {
			// 生成减法算式
			a = rand.Intn(maxNextValue) // 0-100
			// 确保 a ≥ b（结果非负）
			b = rand.Intn(a + 1) // 0 到 a

			// 创建唯一标识符（减法顺序重要）
			key := fmt.Sprintf("%d-%d", a, b)

			// 检查是否已存在相同算式
			if seen[key] {
				continue
			}

			// 标记该算式为已见
			seen[key] = true

			problem = [3]string{fmt.Sprintf("%d", a), "-", fmt.Sprintf("%d", b)}
		}

		// 添加到结果集
		problems = append(problems, problem)
	}
	return problems
}

// @Description 生成100以内加法算式
func (ps *WechatService) generateAdditionWithinHundred(count int, max int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	problems := make([][2]int, 0, count)
	seen := make(map[string]bool) // 用于检测重复的映射

	maxNextValue := max + 1
	for len(problems) < count {
		// 生成第一个数 (0-100)
		a := rand.Intn(maxNextValue)
		// 计算第二个数的最大值 (确保 a+b ≤ 100)
		maxB := max - a
		// 生成第二个数 (0 到 maxB)
		b := rand.Intn(maxB + 1)

		// 创建唯一标识符（考虑加法交换律）
		minVal, maxVal := a, b
		if a > b {
			minVal, maxVal = b, a
		}
		key := fmt.Sprintf("%d+%d", minVal, maxVal)

		// 检查是否已存在相同算式
		if seen[key] {
			continue
		}

		// 标记该算式为已见
		seen[key] = true
		// 添加到结果集
		problems = append(problems, [2]int{a, b})
	}
	return problems
}

// @Description 生成100以内减法算式
func (ps *WechatService) generateSubtractionWithinHundred(count int, max int) [][2]int {
	rand.Seed(time.Now().UnixNano())
	problems := make([][2]int, 0, count)
	seen := make(map[string]bool) // 用于检测重复的映射

	maxNextValue := max + 1
	for len(problems) < count {
		// 生成被减数 (0-100)
		a := rand.Intn(maxNextValue)
		// 生成减数 (0 到 a)
		b := rand.Intn(a + 1)

		// 创建唯一标识符
		key := fmt.Sprintf("%d-%d", a, b)

		// 检查是否已存在相同算式
		if seen[key] {
			continue
		}

		// 标记该算式为已见
		seen[key] = true
		// 添加到结果集
		problems = append(problems, [2]int{a, b})
	}
	return problems
}
