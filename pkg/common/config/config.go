package config

var Config struct {
	App struct {
		Timezone  string `yaml:"timezone"`
		Discovery string `yaml:"discovery"`
		Env       string `yaml:"env"`
		ProxyURL  string `yaml:"proxyURL"`
	} `yaml:"app"`
	MysqlMaster struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"mysql_master"`
	MysqlSlave struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"mysql_slave"`
	Redis struct {
		ClusterMode bool     `yaml:"clusterMode"`
		Address     []string `yaml:"address"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
	} `yaml:"redis"`

	XXLJobDB struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"xxlJobDB"`

	XXLJob struct {
		AdminServer  string `yaml:"adminServer"`
		AccessToken  string `yaml:"accessToken"`
		ExecutorName string `yaml:"executorName"`
		ExecutorIp   string `yaml:"executorIp"`
		ExecutorPort string `yaml:"executorPort"`
	} `yaml:"xxlJob"`

	RpcPort struct {
		SiteRPCPort         int `yaml:"siteRPCPort"`
		AgentRPCPort        int `yaml:"agentRPCPort"`
		ActRPCPort          int `yaml:"actRPCPort"`
		FinanceRPCPort      int `yaml:"financeRPCPort"`
		GameRPCPort         int `yaml:"gameRPCPort"`
		LiveRPCPort         int `yaml:"liveRPCPort"`
		ServerApiGatePort   int `yaml:"serverApiGatePort"`
		ActivityApiGatePort int `yaml:"queueJobGatePort"`
		LiveApiGatePort     int `yaml:"liveApiGatePort"`
	} `yaml:"rpcPort"`
	RpcName struct {
		SiteRpcName         string `yaml:"siteRpcName"`
		AgentRPCName        string `yaml:"agentRPCName"`
		ActRPCName          string `yaml:"actRPCName"`
		FinanceRPCName      string `yaml:"financeRPCName"`
		GameRPCName         string `yaml:"gameRPCName"`
		LiveRPCName         string `yaml:"liveRPCName"`
		ServerApiGateName   string `yaml:"serverApiGateName"`
		ActivityApiGateName string `yaml:"queueJobGateName"`
		LiveApiGateName     string `yaml:"liveApiGateName"`
	} `yaml:"rpcName"`

	Log struct {
		LogLevel string `yaml:"logLevel"`
	} `yaml:"log"`
	Etcd struct {
		Schema   string   `yaml:"schema"`
		Addr     []string `yaml:"addr"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"etcd"`
	Sms struct {
		IsOpen                bool   `yaml:"isOpen"`
		PlatformName          string `yaml:"platformName"`
		ValidTime             int    `yaml:"validTime"`
		SendCooldown          int    `yaml:"sendCooldown"`
		MaxSendTimesOf24Hours int    `yaml:"maxSendTimesOf24Hours"`
		CodeLen               int    `yaml:"codeLen"`
		DefaultCode           string `yaml:"defaultCode"`
		Platform              struct {
			Paasoo struct {
				Endpoint string `yaml:"endpoint"`
				Key      string `yaml:"key"`
				Secret   string `yaml:"secret"`
				From     string `yaml:"from"`
			} `yaml:"paasoo"`
		} `yaml:"platform"`
	} `yaml:"sms"`
	Jwt struct {
		Key string `yaml:"key"`
	} `yaml:"jwt"`
	S3 struct {
		AccessKeyID     string `yaml:"accessKeyID"`
		SecretAccessKey string `yaml:"secretAccessKey"`
		Region          string `yaml:"region"`
		Bucket          string `yaml:"bucket"`
		Endpoint        string `yaml:"endpoint"`
	} `yaml:"s3"`
	Agora struct {
		ResetApi        string `yaml:"resetApi"`
		AppID           string `yaml:"appID"`
		AppKey          string `yaml:"appKey"`
		OrgName         string `yaml:"orgName"`
		AppName         string `yaml:"appName"`
		Certificate     string `yaml:"certificate"`
		RtcNotifySecret string `yaml:"rtcNotifySecret"`
		Rtc             struct {
			Uri            string `yaml:"uri"`
			CustomerKey    string `yaml:"customerKey"`
			CustomerSecret string `yaml:"customerSecret"`
		} `yaml:"rtc"`
	} `yaml:"agora"`
	RocketMQ struct {
		RocketMQAddr []string `yaml:"rocketMQAddr"`
	} `yaml:"rocketMQ"`
	Kafka struct {
		KafkaAddr []string `yaml:"kafkaAddr"`
	} `yaml:"kafka"`
	OpenInstall struct {
		ApiKey        string `yaml:"apiKey"`
		ApiSecret     string `yaml:"apiSecret"`
		GrowthBaseURL string `yaml:"growthBaseURL"`
		ActiveBaseURL string `yaml:"activeBaseURL"`
	} `yaml:"openInstall"`
}

func GetServiceNames() []string {
	return []string{
		Config.RpcName.SiteRpcName,
	}
}
