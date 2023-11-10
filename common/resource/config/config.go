package config

import (
	"frozen-go-cms/common/mylogrus"
	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
	"os"
	"runtime"
	"strconv"
)

// 数据库的配置
type MysqlConfig struct {
	MYSQL_HOST     string
	MYSQL_USERNAME string
	MYSQL_PASSWORD string
	MYSQL_DB       string
}

type MysqlCodeConfig struct {
	MYSQL_HOST     string
	MYSQL_USERNAME string
	MYSQL_PASSWORD string
	MYSQL_DB       string
}

// redis配置
type RedisConfig struct {
	REDIS_HOST     string
	REDIS_PASSWORD string
}

// jwt
type JwtConfig struct {
	SECRET     string
	ISSUER_API string
	ISSUER_MGR string
	EXPIRE     string
}

// jwt
type GameJwtConfig struct {
	SECRET        string
	ISSUER_CLIENT string
	ISSUER_SERVER string
	EXPIRE        string
}

// oss
type OssConfig struct {
	OSS_ACCESS_KEY_ID     string
	OSS_ACCESS_KEY_SECRET string
	OSS_ROLE_ARN          string
	OSS_END_POINT         string
	OSS_BUCKET            string
	OSS_CDN               string
	OSS_EXPIRED_TIME      uint
	OSS_STS_POINT         string
	OSS_STS               string
	OSS_STS_AES           string
}

// aws
type AwsConfig struct {
	AWS_BUCKET string
	AWS_CDN    string
	AWS_DIR    string
	CONFIDENCE float32
}

// APP
type AppConfig struct {
	BIZ_SECRET              string
	WEB_SECRET              string
	OPERATION_SECRET        string
	SUPERUSER               string
	OFFICIAL_GROUP          string
	MINIMAL_VERSION_ANDROID int
	MINIMAL_VERSION_IOS     int
	MODERATE                string
}

// googlePay 配置信息
type GooglePayConfig struct {
	JsonKey []byte
}

// 融云
type RongyunConfig struct {
	RONG_CLOUD_APP_KEY    string
	RONG_CLOUD_APP_SECRET string
	RONG_CLOUD_URL        string
}

// 腾讯云
type TencentyunConfig struct {
	TENCENTYUN_APP_ID int
	TENCENTYUN_KEY    string
	TX_OVERSEA_APP_ID int
	TX_OVERSEA_KEY    string
}

// emas
type EmasConfig struct {
	ANDROID_APP_KEY    string
	ANDROID_APP_SECRET string
	REGION_ID          string
	ACCESS_KEY_ID      string
	ACCESS_KEY_SECRET  string
	IOS_APP_KEY        string
	IOS_APP_SECRET     string
	APNS               string
}

// 声网
type AgoraConfig struct {
	APP_ID          string
	APP_CERTIFICATE string
	CUSTOMER_KEY    string
	CUSTOMER_SECRET string
}

// 腾讯TRTC
type TRTCConfig struct {
	APP_ID          int
	APP_CERTIFICATE string
}

// 匹配的配置
type MatchConfig struct {
	//一开始匹配的默认时长（单位：秒）
	MATCH_FREE_TIME int
	//一开始匹配的默认时长（单位：秒）VIP
	MATCH_FREE_TIME_VIP int
	//免费加时的时长 （单位：秒）
	MATCH_ADD_TIME_FREE int
	//匹配的声网的延迟加时(单位：秒)
	MATCH_AGORA_TIME int
	//匹配周期（单位：秒）
	MATCH_CYCLE int
	//过期时间（单位：秒），用户redisCache时间
	MATCH_USER_EXPIRES int
	//pb match_success中， wait_duration 开始/下一个时间（单位：秒）
	MATCH_SUCCESS_WAIT_DURATION uint32
	//pb match_success中， single_wait_time_in_sec 单方等待连接最长时间（单位：秒）
	MATCH_SUCCESS_SINGLE_WAIT_TIME_IN_SEC uint32
	//pb match_success中， dual_wait_time_in_sec 双方连接中最长时间（单位：秒）
	MATCH_SUCCESS_DUAL_WAIT_TIME_IN_SEC uint32
}

// 在线
type OnlineConfig struct {
	//在线周期
	ONLINE_CYCLE int
	//在线过期时间
	ONLINE_USER_EXPIRES int
}

// 1对1视频
type VideoConfig struct {
	VIDEO_DAILY_FREE_NUM int
	//一开始匹配的默认时长（单位：秒）
	VIDEO_FREE_TIME int
	//一开始匹配的默认时长（单位：秒）,vip
	VIDEO_FREE_TIME_VIP int
	//免费加时的时长 （单位：秒）
	VIDEO_ADD_TIME_FREE int
	//声网的延迟加时(单位：秒)
	VIDEO_AGORA_TIME int
	//1分钟视频，普通用户价格
	VIDEO_MINUTE_NORMAL int
	//1分钟视频，公会用户价格
	VIDEO_MINUTE_UNION int
}

// 会话
type SessionConfig struct {
	SESSION_DAILY_FREE_NUM int
	GUILD_USER_HELLO_DAY   int
}

type BeanConfig struct {
	DIAMOND_BEAN_RATE int
}

type GemConfig struct {
	DIAMOND_GEM_RATE int
}

type H5Config struct {
	USER_LEVEL               string
	GROUP_SUPPORT            string
	LUCKY_WHEEL              string
	WEEKLY_STAR              string
	WEEKLY_CP                string
	COUNTRY_STAR             string
	NOBLE_BUY_IOS            string
	GUILD_DATA_URL           string
	MGR_GUILD_DATA_URL       string
	RANKING_PINK_DIAMOND_URL string
	AGENT_SHARE_URL          string
	AGENT_SHARE_ICON         string
}

type GroupImConfig struct {
	MSG_SORT_EXPIRE   int
	MSG_SORT_SNAP     int
	MSG_PARALLEL_SIZE int
}

type GradeConfig struct {
	//魅力速度
	CHARM_SPEED_VIP int
	//活跃
	ACTITY_SPEED_VIP int
	//财富
	WEALTH_SPEED_VIP int
}

type LikeConfig struct {
	//喜欢人数
	I_LIKE_NUM int
	//喜欢人数VIP
	I_LIKE_NUM_VIP int
	//喜欢人数贵族
	I_LIKE_NUM_NOBLE int
}

type ApplePayConfig struct {
	PASSWORD string
}

type RegisterConfig struct {
	IMEI_TOTAL          int
	IMEI_OAUTH          int
	ACCOUNT_IP          int
	ACCOUNT_IP_DURATION int
}

type BannerConfig struct {
	GIFT_BANNER_LEVEL1 int
	GIFT_BANNER_LEVEL2 int
	GIFT_BANNER_LEVEL3 int
}

type DiamondConfig struct {
	DAILY_LOGIN_IMEI_LIMIT int
	DAILY_LOGIN_IP_LIMIT   int
	PRIVATE_GIFT_RETURN    int
	NEW_USER_INVITE_AWARD  uint32
}

type LuckWheelConfig struct {
	MINIMAL_PARTICIPANT   int // 轮盘开始最少需要的参与人数
	WAIT_TIMELONG         int // 等待轮盘开始的时长（分钟）
	WINNER_DIAMOND_BANNER int //全服广播钻石门槛
}

// 自定义主题
type GroupCustomThemeConfig struct {
	PIC_LIMIT int //图片数量
	DAY       int //有效天数
}

type GiftConfig struct {
	WALL_DIAMOND int //上礼物墙，礼物钻石金额
}

type DailyConfig struct {
	LOGIN_COMMON int
	LOGIN_VIP    int
}

type FruitTycoonConfig struct {
	BIG_WINNER_THRESDHOLD uint
	BIG_WINNER_LOW        uint
	BIG_WINNER_HIGH       uint
	POOL_RATIO            uint32
	WATERMELON_RATIO      uint32
}

type ActivityConfig struct {
	COUNTRY_STAR_POOL_RATIO     uint32
	COUNTRY_STAR_ORDINARY_RATIO uint32
}

type CheckoutConfig struct {
	URL             string
	AUTHORIZATION   string
	H5              string
	HILO_SECRET_KEY string
}

type RiskControlConfig struct {
	USER_QPS_LIMIT     int64
	USER_URL_QPS_LIMIT int64
}

type PayerMaxConfig struct {
	URL                string
	KEY                string
	MERCHANT_ID        string
	BIZ_TYPE           string
	VERSION            string
	FRONT_CALLBACK_URL string
	SHOW_RESULT        string
	EXPIRE_TIME        string
	LANGUAGE           string
}

type SudConfig struct {
	API_LIST string
}

type URLConfig struct {
	BIZ_HTTP string
}

const (
	LOCAL   string = "local"
	DEBUG   string = "debug"
	RELEASE string = "release"
)

var mysqlConfigData MysqlConfig
var mysqlCodeConfigData MysqlCodeConfig
var redisConfigData RedisConfig
var jwtConfigData JwtConfig
var gameJwtConfigData GameJwtConfig
var appConfigData AppConfig
var ossConfigData OssConfig
var awsConfigData AwsConfig
var googlePayData GooglePayConfig
var rongyunData RongyunConfig
var tencentyunData TencentyunConfig
var emasData EmasConfig
var agora AgoraConfig
var trtc TRTCConfig
var matchData MatchConfig
var onlineData OnlineConfig
var sessionData SessionConfig
var videoData VideoConfig
var beanData BeanConfig
var gemData GemConfig
var h5Data H5Config
var groupImData GroupImConfig
var gradeData GradeConfig
var likeData LikeConfig
var applePayData ApplePayConfig
var registerData RegisterConfig
var bannerConfig BannerConfig
var diamondConfig DiamondConfig
var luckyWheelConfig LuckWheelConfig
var groupCustomThemeConfig GroupCustomThemeConfig
var giftConfig GiftConfig
var dailyConfig DailyConfig
var fruitTycoonConfig FruitTycoonConfig
var activityConfig ActivityConfig
var checkoutConfig CheckoutConfig
var riskControl RiskControlConfig
var payerMaxConfig PayerMaxConfig
var mode string
var master bool
var sudConfig SudConfig
var urlConfig URLConfig

func GetConfigMysql() MysqlConfig {
	return mysqlConfigData
}

func GetConfigMysqlCode() MysqlCodeConfig {
	return mysqlCodeConfigData
}

func GetConfigRedis() RedisConfig {
	return redisConfigData
}

func GetConfigJWT() JwtConfig {
	return jwtConfigData
}

func GetConfigGameJWT() GameJwtConfig {
	return gameJwtConfigData
}

func GetConfigApp() AppConfig {
	return appConfigData
}

func GetConfigOss() OssConfig {
	return ossConfigData
}

func GetConfigAws() AwsConfig {
	return awsConfigData
}

func GetConfigGooglePay() GooglePayConfig {
	return googlePayData
}

func GetMode() string {
	return mode
}

func AppIsRelease() bool {
	return GetMode() == RELEASE
}

func AppIsLocal() bool {
	return GetMode() == LOCAL
}

func IsMaster() bool {
	return master
}

func GetOssCDN() string {
	return ossConfigData.OSS_CDN
}

func GetRongyunAppKey() string {
	return rongyunData.RONG_CLOUD_APP_KEY
}

func GetRongyunAppSecret() string {
	return rongyunData.RONG_CLOUD_APP_SECRET
}

func GetRongyunUrl() string {
	return rongyunData.RONG_CLOUD_URL
}

func GetTencentyunAppId() int {
	return tencentyunData.TENCENTYUN_APP_ID
}

func GetTencentyunKey() string {
	return tencentyunData.TENCENTYUN_KEY
}

func GetTxOverSeaAppId() int {
	return tencentyunData.TX_OVERSEA_APP_ID
}

func GetTxOverSeaAppKey() string {
	return tencentyunData.TX_OVERSEA_KEY
}

func GetEmasRegionId() string {
	return emasData.REGION_ID
}

func GetEmasAccessKeyId() string {
	return emasData.ACCESS_KEY_ID
}

func GetEmasAccessKeySecret() string {
	return emasData.ACCESS_KEY_SECRET
}

func GetEmasAndroidAppKey() string {
	return emasData.ANDROID_APP_KEY
}

func GetEmasIosAppKey() string {
	return emasData.IOS_APP_KEY
}

func GetEmasApns() string {
	return emasData.APNS
}

func GetAgoraAppId() string {
	return agora.APP_ID
}

func GetAgoraAppCertificate() string {
	return agora.APP_CERTIFICATE
}

func GetAgoraCustomerKey() string {
	return agora.CUSTOMER_KEY
}

func GetAgoraCustomerSecret() string {
	return agora.CUSTOMER_SECRET
}

func GetMatchConfig() *MatchConfig {
	return &matchData
}

func GetOnlineConfig() *OnlineConfig {
	return &onlineData
}

func GetSessionConfig() SessionConfig {
	return sessionData
}

func GetVideoConfig() VideoConfig {
	return videoData
}

func GetBeanConfig() BeanConfig {
	return beanData
}

func GetGemConfig() GemConfig {
	return gemData
}

func GetH5Config() H5Config {
	return h5Data
}

func GetGroupImConfig() GroupImConfig {
	return groupImData
}

func GetGradeConfig() GradeConfig {
	return gradeData
}

func GetLikeConfig() LikeConfig {
	return likeData
}

func GetApplePayConfig() ApplePayConfig {
	return applePayData
}

func GetRegisterConfig() RegisterConfig {
	return registerData
}

func GetBannerConfig() BannerConfig {
	return bannerConfig
}

func GetDiamondConfig() DiamondConfig {
	return diamondConfig
}

func GetLuckyWheelConfig() LuckWheelConfig {
	return luckyWheelConfig
}

func GetGroupCustomThemeConfig() GroupCustomThemeConfig {
	return groupCustomThemeConfig
}

func GetGiftConfig() GiftConfig {
	return giftConfig
}

func GetDailyConfig() DailyConfig {
	return dailyConfig
}

func GetFruitTycoonConfig() FruitTycoonConfig {
	return fruitTycoonConfig
}

func GetActivityConfig() ActivityConfig {
	return activityConfig
}

func GetCheckoutConfig() CheckoutConfig {
	return checkoutConfig
}

func GetRiskControlConfig() RiskControlConfig {
	return riskControl
}

func GetPayerMaxConfig() PayerMaxConfig {
	return payerMaxConfig
}

func GetSudConfig() SudConfig {
	return sudConfig
}

func GetUrlConfig() URLConfig {
	return urlConfig
}

func GetTRTCConfig() TRTCConfig {
	return trtc
}

func init() {
	str, _ := os.Getwd()
	mylogrus.MyLog.Info(str)

	envDir := ".env"

	//加载环境变量
	if err := godotenv.Load(envDir); err != nil {
		mylogrus.MyLog.Fatalf("Error loading .env err:%v", err)
	}

	//获取环境变量
	mode = os.Getenv("MODE")
	var err error
	master, _ = strconv.ParseBool(os.Getenv("MASTER"))
	mylogrus.MyLog.Infof("My role is %t", master)

	iniDir := mode + ".ini"
	if runtime.GOOS == "darwin" { // mac本地调试
		iniDir = "/var/log/hilo/" + iniDir
	}
	//根据环境变量获取具体的配置，实现多环境配置
	//var conf *ini.File
	conf, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, iniDir)
	if err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	//加载mysql的配置
	if err := conf.Section("DATABASE").MapTo(&mysqlConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("DATABASECODE").MapTo(&mysqlCodeConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("REDIS").MapTo(&redisConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("JWT").MapTo(&jwtConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("GAMEJWT").MapTo(&gameJwtConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("APP").MapTo(&appConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		mylogrus.MyLog.Infof("APP: %+v", appConfigData)
	}

	if err := conf.Section("OSS").MapTo(&ossConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("AWS").MapTo(&awsConfigData); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		if awsConfigData.CONFIDENCE <= 50 {
			awsConfigData.CONFIDENCE = 80
		}
		mylogrus.MyLog.Infof("AWS: %+v", awsConfigData)
	}

	if err := conf.Section("RONGYUN").MapTo(&rongyunData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("TENCENTYUN").MapTo(&tencentyunData); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		mylogrus.MyLog.Info("TENCENTYUN: ", tencentyunData)
	}

	if err := conf.Section("EMAS").MapTo(&emasData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("AGORA").MapTo(&agora); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("TRTC").MapTo(&trtc); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("MATCH").MapTo(&matchData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("ONLINE").MapTo(&onlineData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("SESSION").MapTo(&sessionData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("VIDEO").MapTo(&videoData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("BEAN").MapTo(&beanData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("GEM").MapTo(&gemData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("H5").MapTo(&h5Data); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("GROUPIM").MapTo(&groupImData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("GRADE").MapTo(&gradeData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("LIKE").MapTo(&likeData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("APPLEPAY").MapTo(&applePayData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("REGISTER").MapTo(&registerData); err != nil {
		mylogrus.MyLog.Fatal(err)
	}

	if err := conf.Section("BANNER").MapTo(&bannerConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("DIAMOND").MapTo(&diamondConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		if diamondConfig.NEW_USER_INVITE_AWARD <= 0 {
			diamondConfig.NEW_USER_INVITE_AWARD = 5000
		}
	}
	if err := conf.Section("LUCKY_WHEEL").MapTo(&luckyWheelConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("GROUP_CUSTOM_THEME").MapTo(&groupCustomThemeConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("GIFT").MapTo(&giftConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("DAILY").MapTo(&dailyConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("CHECKOUT").MapTo(&checkoutConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("PAYER_MAX").MapTo(&payerMaxConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("FRUIT_TYCOON").MapTo(&fruitTycoonConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		// 防止未配置或配置错误
		if fruitTycoonConfig.BIG_WINNER_LOW <= 0 {
			fruitTycoonConfig.BIG_WINNER_LOW = 10000
		}
		if fruitTycoonConfig.BIG_WINNER_HIGH <= 0 {
			fruitTycoonConfig.BIG_WINNER_HIGH = 20000
		}
		if fruitTycoonConfig.POOL_RATIO <= 0 || fruitTycoonConfig.POOL_RATIO > 100 {
			fruitTycoonConfig.POOL_RATIO = 20
		}
		if fruitTycoonConfig.WATERMELON_RATIO <= 0 || fruitTycoonConfig.WATERMELON_RATIO > 100 {
			fruitTycoonConfig.WATERMELON_RATIO = 70
		}
		mylogrus.MyLog.Infof("FRUIT_TYCOON: %+v", fruitTycoonConfig)
	}

	if err := conf.Section("ACTIVITY").MapTo(&activityConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		// 防止未配置或配置错误
		if activityConfig.COUNTRY_STAR_POOL_RATIO <= 0 {
			activityConfig.COUNTRY_STAR_POOL_RATIO = 20
		}
		if activityConfig.COUNTRY_STAR_ORDINARY_RATIO <= 0 {
			activityConfig.COUNTRY_STAR_ORDINARY_RATIO = 20
		}
		mylogrus.MyLog.Infof("ACTIVITY: %+v", activityConfig)
	}

	if err := conf.Section("RISK_CONTROL").MapTo(&riskControl); err != nil {
		mylogrus.MyLog.Fatal(err)
	} else {
		if riskControl.USER_QPS_LIMIT <= 0 {
			riskControl.USER_QPS_LIMIT = 128
		}
		if riskControl.USER_URL_QPS_LIMIT <= 0 {
			riskControl.USER_URL_QPS_LIMIT = 64
		}
		mylogrus.MyLog.Infof("RISK_CONTROL: %+v", riskControl)
	}
	if err := conf.Section("SUD").MapTo(&sudConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
	if err := conf.Section("URL").MapTo(&urlConfig); err != nil {
		mylogrus.MyLog.Fatal(err)
	}
}
