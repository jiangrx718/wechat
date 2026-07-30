package miniwechat

// impl_home_data.go 首页玩法卡片配置数据
//
// 与小程序前端 pages/home/home.js 中的 categories 保持一致，
// 后端通过 visible 字段控制每张卡片（及每个分类）是否在小程序首页展示。
// 数据为准静态配置，改动后重启服务即可生效。

// categoryDTO 分类
type categoryDTO struct {
	Key     string    `json:"key"`     // 分类唯一标识
	Name    string    `json:"name"`    // 分类名称
	Icon    string    `json:"icon"`    // 分类图标（emoji）
	Sub     string    `json:"sub"`     // 分类副标题
	Sort    int       `json:"sort"`    // 分类排序（升序）
	Visible bool      `json:"visible"` // 是否展示该分类
	Games   []gameDTO `json:"games"`   // 该分类下的玩法卡片
}

// gameDTO 玩法卡片
type gameDTO struct {
	Key     string `json:"key"`     // 玩法唯一标识（同时是前端 wxss 类名后缀 gcard-<key>）
	Name    string `json:"name"`    // 卡片标题
	Icon    string `json:"icon"`    // 卡片图标（emoji）
	Desc    string `json:"desc"`    // 卡片描述
	Tag     string `json:"tag"`     // 右上角角标（自由文本）
	URL     string `json:"url"`     // 点击跳转路径
	Sort    int    `json:"sort"`    // 卡片排序（升序）
	Visible bool   `json:"visible"` // 是否展示该卡片
}

// dailyEntryDTO 首页「今日推荐玩法」醒目入口卡片
//
// 该卡片的内容（玩法名 / 图标 / 色调）仍由前端按「周一-周日」本地循环生成
// （utils/dailySchedule.js），后端只控制：
//  1. Visible：整张入口卡片是否在首页展示（false 时前端不渲染该卡片）；
//  2. Sub：卡片副标题文案（如「今日推荐玩法 · 连续打卡赢连胜」），便于运营改文案。
//
// 当天推荐的玩法本身是否可见，依赖 categories 中对应 game 的 Visible 控制，
// 这里不再重复下发每天的玩法（避免与本地调度冲突）。
type dailyEntryDTO struct {
	Visible bool   `json:"visible"` // 是否展示「今日推荐玩法」入口卡片
	Sub     string `json:"sub"`     // 卡片副标题文案
}

// dailyEntryData 「今日推荐玩法」入口卡片配置
//
// 调整是否展示该卡片：把 Visible 改为 false，重启服务生效。
// 调整副标题文案：修改 Sub。
var dailyEntryData = dailyEntryDTO{
	Visible: true,
	Sub:     "今日推荐玩法 · 连续打卡赢连胜",
}

// cardsData 全量卡片配置（全量 = true，下发给前端；= false 仅内部维护，不下发）
//
// 调整某张卡片是否在首页展示：把对应 game 的 Visible 改为 false 即可。
// 调整整个分类：把 category 的 Visible 改为 false 即可（该分类下所有卡片都不会展示）。
// 新增玩法：在对应分类的 Games 末尾追加，Sort 取当前最大值 +1。
var cardsData = []categoryDTO{
	{
		Key: "classic", Name: "巧思解谜", Icon: "🧩", Sub: "经典复原 · 逐关突破", Sort: 1, Visible: true,
		Games: []gameDTO{
			{Key: "puzzle", Name: "思维拼图休闲", Icon: "🧩", Desc: "选图拼合 · 多种难度", Tag: "经典", URL: "/pages/index/index", Sort: 1, Visible: true},
			{Key: "huarong", Name: "华容道", Icon: "🔢", Desc: "数字滑块 · 复原闯关", Tag: "经典", URL: "/pages/huarong/huarong", Sort: 2, Visible: true},
			{Key: "lights", Name: "烽火台", Icon: "🔥", Desc: "点灯解谜 · 逐关闯关", Tag: "20关", URL: "/pages/lights/lights", Sort: 3, Visible: true},
			{Key: "flow", Name: "行军路线", Icon: "🛤️", Desc: "同色连线 · 逐关闯关", Tag: "20关", URL: "/pages/flow/flow", Sort: 4, Visible: true},
			{Key: "nonogram", Name: "数织阵图", Icon: "🖼️", Desc: "数字推图 · 画出古风", Tag: "12关", URL: "/pages/nonogram/nonogram", Sort: 5, Visible: true},
			{Key: "minesweeper", Name: "经典扫雷", Icon: "🚩", Desc: "长按插旗 · 推理排雷", Tag: "3档", URL: "/pages/minesweeper/minesweeper", Sort: 6, Visible: true},
			{Key: "spot", Name: "双图找茬", Icon: "🔍", Desc: "程序生图 · 圈出不同", Tag: "找茬", URL: "/pages/spot/spot", Sort: 7, Visible: true},
			{Key: "sokoban", Name: "运粮入库", Icon: "📦", Desc: "推箱入库 · 12关闯关", Tag: "经典", URL: "/pages/sokoban/sokoban", Sort: 8, Visible: true},
			{Key: "bonfire", Name: "烽火传令", Icon: "🔥", Desc: "一笔遍历 · 20关闯关", Tag: "20关", URL: "/pages/bonfire/bonfire", Sort: 9, Visible: true},
			{Key: "dujiang", Name: "都江堰", Icon: "💧", Desc: "旋转接水 · 12关闯关", Tag: "12关", URL: "/pages/dujiang/dujiang", Sort: 10, Visible: true},
		},
	},
	{
		Key: "word", Name: "文墨雅趣", Icon: "📚", Sub: "诗词成语 · 妙笔生花", Sort: 2, Visible: true,
		Games: []gameDTO{
			{Key: "poetry", Name: "飞花令", Icon: "🌸", Desc: "诗词填字 · 按卷闯关", Tag: "4大类别", URL: "/pages/poetry/poetry", Sort: 1, Visible: true},
			{Key: "wordchain", Name: "三国拼词", Icon: "📜", Desc: "成语接龙 · 百回闯关", Tag: "120回", URL: "/pages/wordchain/wordchain", Sort: 2, Visible: true},
			{Key: "idiomchain", Name: "成语接龙", Icon: "🔗", Desc: "同音相接 · 百回闯关", Tag: "120关", URL: "/pages/idiomchain/idiomchain", Sort: 3, Visible: true},
			{Key: "idiomspot", Name: "成语找茬", Icon: "🔍", Desc: "找出错字 · 30关精选", Tag: "300题", URL: "/pages/idiomspot/idiomspot", Sort: 4, Visible: true},
			{Key: "idiommatch", Name: "成语消消乐", Icon: "💥", Desc: "点字成词 · 限时连消", Tag: "限时", URL: "/pages/idiommatch/idiommatch", Sort: 5, Visible: true},
			{Key: "poetryflash", Name: "诗词闪电", Icon: "⚡", Desc: "限时选字 · 无尽挑战", Tag: "上瘾", URL: "/pages/poetryflash/poetryflash", Sort: 6, Visible: true},
			{Key: "wordle", Name: "字谜方阵", Icon: "🔤", Desc: "成语 Wordle · 绿黄灰破译", Tag: "每日", URL: "/pages/wordle/wordle", Sort: 7, Visible: true},
		},
	},
	{
		Key: "speed", Name: "手速挑战", Icon: "⚡", Sub: "反应合成 · 极限连击", Sort: 3, Visible: true,
		Games: []gameDTO{
			{Key: "fibonacci", Name: "校场点兵", Icon: "⚔️", Desc: "斐波那契 · 滑动合并", Tag: "3+5→8", URL: "/pages/fibonacci/fibonacci", Sort: 1, Visible: true},
			{Key: "stack", Name: "通天塔", Icon: "🗼", Desc: "一指叠塔 · 极限挑战", Tag: "上瘾", URL: "/pages/stack/stack", Sort: 2, Visible: true},
			{Key: "mergecube", Name: "幂次方块", Icon: "🟧", Desc: "数字合成 · 连消翻倍", Tag: "2048", URL: "/pages/mergecube/mergecube", Sort: 3, Visible: true},
			{Key: "block", Name: "方块奇兵", Icon: "🟦", Desc: "七形旋转 · 消行连击", Tag: "Tetris", URL: "/pages/block/block", Sort: 4, Visible: true},
			{Key: "tri", Name: "飞花三消", Icon: "🎯", Desc: "交换连字 · 成语暴击", Tag: "三消", URL: "/pages/tri/tri", Sort: 5, Visible: true},
			{Key: "ride", Name: "千里走单骑", Icon: "🐎", Desc: "跳跃跨障 · 无尽跑酷", Tag: "跑酷", URL: "/pages/ride/ride", Sort: 6, Visible: true},
			{Key: "lianhuan", Name: "连环计", Icon: "🎆", Desc: "点选连消 · 大块高分", Tag: "连消", URL: "/pages/lianhuan/lianhuan", Sort: 7, Visible: true},
			{Key: "liuhe", Name: "六合阵", Icon: "⬡", Desc: "六边形版 · 2048合成", Tag: "2048", URL: "/pages/liuhe/liuhe", Sort: 8, Visible: true},
		},
	},
	{
		Key: "strategy", Name: "策略防御", Icon: "🛡️", Sub: "塔防守卫 · 排兵布阵", Sort: 4, Visible: true,
		Games: []gameDTO{
			{Key: "td", Name: "萝卜保卫战", Icon: "🥕", Desc: "建塔防御 · 12关闯关", Tag: "塔防", URL: "/pages/td/td", Sort: 1, Visible: true},
			{Key: "autochess", Name: "合成战棋", Icon: "♟️", Desc: "合成升星 · 无尽防守", Tag: "自走棋", URL: "/pages/autochess/autochess", Sort: 2, Visible: true},
			{Key: "tank", Name: "坦克大战", Icon: "🛡️", Desc: "驾驶坦克 · 消灭敌军", Tag: "FC", URL: "/pages/tank/tank", Sort: 3, Visible: true},
		},
	},
}
