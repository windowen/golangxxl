package job

// PunCompanyJob 企业招聘职位表
type PunCompanyJob struct {
	Id            int    `gorm:"column:id" json:"id"`                                   // 职位ID（主键，自增）
	Uid           int    `gorm:"column:uid" json:"uid"`                                 // 企业用户ID
	Name          string `gorm:"column:name" json:"name"`                               // 职位名称
	ComName       string `gorm:"column:com_name" json:"com_name"`                       // 公司名称
	Hy            int    `gorm:"column:hy" json:"hy"`                                   // 行业ID
	Job1          int    `gorm:"column:job1;default:45" json:"job1"`                    // 一级职位分类ID
	Job1Son       int    `gorm:"column:job1_son;default:101" json:"job1_son"`           // 二级职位分类ID
	JobPost       int    `gorm:"column:job_post;default:806" json:"job_post"`           // 三级职位分类ID（最细分类）
	ProvinceId    int    `gorm:"column:provinceid;default:3409" json:"provinceid"`      // 工作省份ID
	CityId        int    `gorm:"column:cityid;default:3410" json:"cityid"`              // 工作城市ID
	ThreeCityId   int    `gorm:"column:three_cityid;default:3410" json:"three_cityid"`  // 工作区县ID
	Cert          string `gorm:"column:cert" json:"cert"`                               // 企业是否已认证（1是，0否）
	Type          int    `gorm:"column:type" json:"type"`                               // 职位性质（如全职、兼职）
	Number        int    `gorm:"column:number" json:"number"`                           // 招聘人数
	Exp           int    `gorm:"column:exp;default:127" json:"exp"`                     // 工作经验要求ID，默认127
	Report        int    `gorm:"column:report;default:54" json:"report"`                // 到岗时间要求，默认54
	Sex           int    `gorm:"column:sex;default:3" json:"sex"`                       // 性别要求（0不限），默认3
	Edu           int    `gorm:"column:edu;default:65" json:"edu"`                      // 学历要求，默认65
	Marriage      int    `gorm:"column:marriage;default:72" json:"marriage"`            // 婚姻状况要求，默认72
	Description   string `gorm:"column:description" json:"description"`                 // 职位描述（支持HTML）
	XuAnShang     int    `gorm:"column:xuanshang" json:"xuanshang"`                     // 是否悬赏职位（1是，0否）
	XsDate        int    `gorm:"column:xsdate" json:"xsdate"`                           // 悬赏到期时间（时间戳）
	SDate         int    `gorm:"column:sdate" json:"sdate"`                             // 发布时间（时间戳）
	EDate         int    `gorm:"column:edate" json:"edate"`                             // 职位截止时间（时间戳）
	JobHits       int    `gorm:"column:jobhits" json:"jobhits"`                         // 职位浏览次数
	LastUpdate    string `gorm:"column:lastupdate" json:"lastupdate"`                   // 职位最后更新时间
	Rec           int    `gorm:"column:rec" json:"rec"`                                 // 是否推荐职位（1推荐，0不推荐）
	CloudType     int    `gorm:"column:cloudtype" json:"cloudtype"`                     // 云同步类型（预留）
	State         int    `gorm:"column:state" json:"state"`                             // 审核状态（1通过，0未审核，2未通过）
	StatusBody    string `gorm:"column:statusbody" json:"statusbody"`                   // 审核不通过原因
	Age           int    `gorm:"column:age" json:"age"`                                 // 年龄要求范围描述
	Lang          string `gorm:"column:lang" json:"lang"`                               // 语言要求（可多项）
	Welfare       string `gorm:"column:welfare" json:"welfare"`                         // 职位福利（逗号分隔）
	Pr            int    `gorm:"column:pr;default:20" json:"pr"`                        // 企业性质ID
	Mun           int    `gorm:"column:mun;default:30" json:"mun"`                      // 企业规模ID
	ComProvinceId int    `gorm:"column:com_provinceid;default:2" json:"com_provinceid"` // 公司所在省份ID
	Rating        int    `gorm:"column:rating" json:"rating"`                           // 企业会员等级ID
	Status        int    `gorm:"column:status" json:"status"`                           // 职位状态（1显示，0隐藏）
	Urgent        int    `gorm:"column:urgent" json:"urgent"`                           // 是否紧急招聘
	RStatus       int    `gorm:"column:r_status" json:"r_status"`                       // 是否首页推荐
	EndEmail      int    `gorm:"column:end_email" json:"end_email"`                     // 接收简历邮箱
	UrgentTime    int    `gorm:"column:urgent_time" json:"urgent_time"`                 // 紧急招聘截止时间
	ComLogo       string `gorm:"column:com_logo" json:"com_logo"`                       // 公司Logo地址
	AutoType      int    `gorm:"column:autotype" json:"autotype"`                       // 是否自动刷新
	AutoTime      int    `gorm:"column:autotime" json:"autotime"`                       // 自动刷新时间戳
	IsLink        int    `gorm:"column:is_link" json:"is_link"`                         // 是否使用公司联系方式（1是，0否）
	LinkType      int    `gorm:"column:link_type;default:1" json:"link_type"`           // 联系方式类型（1默认，2自定义）
	Source        int    `gorm:"column:source;default:4" json:"source"`                 // 职位来源（10平台，1API，2采集，3apk，4telegram_bot）
	RecTime       int    `gorm:"column:rec_time" json:"rec_time"`                       // 推荐截止时间
	SNum          int    `gorm:"column:snum" json:"snum"`                               // 已投递人数
	OperaTime     int    `gorm:"column:operatime" json:"operatime"`                     // 后台操作时间
	Did           int    `gorm:"column:did" json:"did"`                                 // 分站ID（多城市支持）
	IsEmail       int    `gorm:"column:is_email" json:"is_email"`                       // 是否开启简历投递提醒邮件
	MinSalary     int    `gorm:"column:minsalary" json:"minsalary"`                     // 最低薪资（元）
	MaxSalary     int    `gorm:"column:maxsalary" json:"maxsalary"`                     // 最高薪资（元）
	SharePack     int    `gorm:"column:sharepack" json:"sharepack"`                     // 是否开启分享红包
	RewardPack    int    `gorm:"column:rewardpack" json:"rewardpack"`                   // 是否开启推荐赏金
	IsGraduate    int    `gorm:"column:is_graduate" json:"is_graduate"`                 // 是否接受应届生
	X             string `gorm:"column:x" json:"x"`                                     // 地图经度
	Y             string `gorm:"column:y" json:"y"`                                     // 地图纬度
	ZUid          int    `gorm:"column:zuid" json:"zuid"`                               // 后台创建者ID
	ExpReq        string `gorm:"column:exp_req" json:"exp_req"`                         // 经验要求文本
	EduReq        string `gorm:"column:edu_req" json:"edu_req"`                         // 学历要求文本
	SexReq        string `gorm:"column:sex_req" json:"sex_req"`                         // 性别要求文本
	MinAgeReq     int    `gorm:"column:minage_req" json:"minage_req"`                   // 最小年龄要求
	MaxAgeReq     int    `gorm:"column:maxage_req" json:"maxage_req"`                   // 最大年龄要求
}

// TableName 指定数据库表名
func (PunCompanyJob) TableName() string {
	return "pun_company_job"
}
