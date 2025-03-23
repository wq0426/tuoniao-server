package common

import (
	"time"
)

const (
	ONE_DAY_SECONDS       = 86400
	HOUR_SECONDS          = 3600
	STAR_EXPLORE_DURATION = 30 * 24 * time.Hour

	STANDARD_TIMESTR = "2006-01-02 15:04:05"
	DATE_FORMAT      = "2006-01-02"

	REDIS_PREFIX_REGISTER = "register."
	REDIS_PREFIX_LOGIN    = "login."
	REDIS_PREFIX_RESET    = "reset."
	REDIS_PREFIX_CHANGE   = "change."
	REDIS_PREFIX_REBIND   = "rebind."
	REDIS_PREFIX_FRIEND   = "friend."
	REDIS_PREFIX_REALTIME = "realtime."

	CODE_TYPE_REGISTER = 1
	CODE_TYPE_LOGIN    = 2
	CODE_TYPE_RESET    = 3
	CODE_TYPE_CHANGE   = 4
	// 重新绑定
	CODE_TYPE_REBIND = 5
	// 交友密码
	CODE_TYPE_FRIEND = 6
	// 实时查询
	CODE_TYPE_REALTIME = 7
	// 设备升级
	CODE_TYPE_DEVICE_UPGRADE = 8
	// login expire
	CODE_TYPE_LOGIN_EXPIRE = 9

	REDIS_PREFIX_DEVICE_UPGRADE = "country.device_upgrade."
	REDIS_PREFIX_LOGIN_EXPIRE   = "country.login_expire."

	NICKNAME_DEFAULT   = "Unknown"
	AVATAR_URL_DEFAULT = "https://txtimages.oss-cn-beijing.aliyuncs.com/profile/avatar.png"

	ROLE_NORMAL         = 1 // 普通会员
	ROLE_VIP            = 2 // 高级会员
	ROLE_FARMER_BEGIN   = 3 // 初级农场主
	ROLE_FARMER_VIP     = 4 // 高级农场主
	ROLE_FARMER_EXPERT  = 5 // 资深农场主
	ROLE_FARMER_PARTNER = 6 // 合伙人

	STATUS_NORMAL = 1 //  正常
	STATUS_CANCEL = 2 // 注销

	STATUS_PLUNDER_NO  = 0 //  未被掠夺
	STATUS_PLUNDER_YES = 1 // 已被掠夺

	AUTH_STATUS_NOT = 0 // 未认证
	AUTH_STATUS_YES = 1 // 已认证

	PUSH_STATUS_NOT = 0 // 未开启推送
	PUSH_STATUS_YES = 1 // 已开启推送

	DEMO_STATUS_NOT = 0 // 未开启演示
	DEMO_STATUS_YES = 1 // 已开启演示

	MAX_BARRACKS_LEVEL = 5

	// 商品状态（1:正常  2:下架）
	PRODUCT_STATUS_NORMAL = 1
	PRODUCT_STATUS_OFF    = 2

	REGION_CN_BEIJING = "cn-beijing"
	OSS_ENDPOINT      = "oss-cn-beijing.aliyuncs.com"
	HTTP_PREFIX       = "http://"
	HTTPS_PREFIX      = "https://"
	BUCKET_NAME       = "countrybattle"

	RESOURCE_AVATAR = "resource/avatar"

	POWER_DEFAULT          = 1000
	BARRACKS_LEVEL_DEFAULT = 1

	IS_FIRST_USED    = 0
	NOT_IS_FIRST_USE = 1

	MAX_CARDS_NUM = 80

	REWARD_TYPE_GRAIN       = 1
	REWARD_TYPE_HORSE       = 2
	REWARD_TYPE_IRON        = 3
	REWARD_TYPE_PLUNDER     = 4
	REWARD_TYPE_PUTDOWN     = 5
	REWARD_TYPE_BATTLE      = 6
	REWARD_TYPE_ESCORT      = 7
	REWARD_TYPE_ADVERTISING = 8
	REWARD_TYPE_DONATE      = 9

	MODULE_TYPE_USER_INFO = "user_info"
	MODULE_TYPE_ESCORT    = "escort"

	// 查询缓存prefix
	PREFFIX_USER_INFO = "user_info."
	// 用户资产
	PREFFIX_USER_ASSET = "user_asset."
	// 商品信息
	PREFFIX_PRODUCT_INFO = "product_info."
	// 购买商品列表
	PREFFIX_PRODUCT_LIST = "product_list."
	// 我的商品列表
	PREFFIX_MY_PRODUCT = "my_product."
	// 获取设置缓存
	PREFFIX_GLOBAL_SETTINGS = "global_settings."
	// 转盘缓存
	PREFFIX_TURNTABLE = "turntable.percent"
	// 打卡奖励配置
	PUNCHING_REWARD_SETTINGS = "punching_reward.settings"

	PUBLISH_PRODUCT_STATUS_NORMAL = 1 // 挂单中
	PUBLISH_PRODUCT_STATUS_BARGIN = 2 // 已成交
	PUBLISH_PRODUCT_STATUS_CANCEL = 3 // 下架

	PUBLISH_TYPE_SELL = 1 // 出售
	PUBLISH_TYPE_BUY  = 2 // 求购

	// 资产类型(1:积分 2:余额)
	ASSET_TYPE_POINT   = 1
	ASSET_TYPE_BALANCE = 2
	ASSET_TYPE_COUPON  = 3

	// 业务类型(1:充值 2:提现 3:兑换 4:订单)
	BUSINESS_TYPE_RECHARGE = 1
	BUSINESS_TYPE_WITHDRAW = 2
	BUSINESS_TYPE_EXCHANGE = 3
	BUSINESS_TYPE_ORDER    = 4

	// 动作类型(1:使用 2:奖励 3:购买 4:提现 5:充值 6:兑换 7:退款)
	ACTION_TYPE_USE      = 1
	ACTION_TYPE_REWARD   = 2
	ACTION_TYPE_BUY      = 3
	ACTION_TYPE_WITHDRAW = 4
	ACTION_TYPE_RECHARGE = 5
	ACTION_TYPE_EXCHANGE = 6
	ACTION_TYPE_REFUND   = 7

	// 产品类型(1:生肖 2:命数)
	PRODUCT_TYPE_ZODIAC = 1
	PRODUCT_TYPE_LIFE   = 2

	// 方法名称
	ACITON_NAME_CREATE_ROOM   = "create_room"
	ACITON_NAME_LIST_ROOM     = "list_room"
	ACTION_NAME_JOIN_ROOM     = "join_room"
	ACTION_NAME_MATCH_ROOM    = "match_room"
	ACTION_NAME_GET_ROOM      = "get_room"
	ACTION_NAME_GET_USER      = "get_user"
	ACTION_NAME_PLAY_CARD     = "play_card"
	ACTION_NAME_OPEN_CARD     = "open_card"
	ACTION_NAME_EXIT_ROOM     = "exit_room"
	ACTION_NAME_SEND_CARD     = "send_card"
	ACTION_NAME_DO_ACTION     = "do_action"
	ACTION_NAME_SHOT          = "shot"
	ACTION_NAME_PUSH_SHOT     = "push_shot"
	ACTION_NAME_PLAYER_UPDATE = "player_update"
	ACTION_NAME_GAME_OVER     = "game_over"
	ACTION_NAME_AUDIO         = "audio"
	ACTION_NAME_PUSH_AUDIO    = "push_audio"

	// 房间状态
	ROOM_STATUS_MATCHING    = 1
	ROOM_STATUS_IN_PROGRESS = 2
	ROOM_STATUS_CHANGING    = 3
	ROOM_STATUS_FINISHED    = 4

	// 房间类型
	ROOM_TYPE_TWO_PLAYERS   = 1
	ROOM_TYPE_THREE_PLAYERS = 2
	ROOM_TYPE_FOUR_PLAYERS  = 3

	// 游戏模式
	GAME_MODE_SINGLE_PLAYER = 1
	GAME_MODE_MULTIPLAYER   = 2

	// 卡牌名称和ID
	CARD_NAME_Q          = 1
	CARD_NAME_K          = 2
	CARD_NAME_A          = 3
	CARD_NAME_JOKER      = 4
	CARD_NAME_Q_NAME     = "Q"
	CARD_NAME_K_NAME     = "K"
	CARD_NAME_A_NAME     = "A"
	CARD_NAME_JOKER_NAME = "JOKER"

	// 座位编号
	SEAT_NUMBER_ONE   = 1
	SEAT_NUMBER_TWO   = 2
	SEAT_NUMBER_THREE = 3
	SEAT_NUMBER_FOUR  = 4

	// 执行状态
	ACTION_STATUS_WAITING  = 1
	ACTION_STATUS_OPENING  = 2
	ACTION_STATUS_SHOOTING = 3

	// 短信服务，模版ID和公司名
	SMS_SIGN     = "深圳市壹号熊网络科技"
	SMS_TEMPLATE = "SMS_305080252"

	STORE_ID   = 1
	STORE_NAME = "鸵小妥"
	STORE_LOGO = "https://txtimages.oss-cn-beijing.aliyuncs.com/icons/dianpu.png"

	// 订单状态
	ORDER_STATUS_PENDING  = 0 // 待付款
	ORDER_STATUS_SHIPPED  = 1 // 待发货
	ORDER_STATUS_RECEIVED = 2 // 待收货
	ORDER_STATUS_EVALUATE = 3 // 待评价
	ORDER_STATUS_COMPLETE = 4 // 已完成
	ORDER_STATUS_CLOSED   = 5 // 已关闭
	ORDER_STATUS_EXPIRED  = 6 // 已过期
	ORDER_STATUS_REFUNDED = 7 // 已退款
)
