package Plugins

type pluginScan func(PluginService)

type PluginService struct {
	Host string
	//protocol string
	//user []string
	//pass []string
}

var PluginModel map[string]string

var (
	Plugin map[string]pluginScan
)

func init() {
	Plugin = map[string]pluginScan{
		"401":       Force401Scan,
		"PathCrack": PathCrack,
		"JsFinder":  JsFinder,
		"PortCrack": PortCrack,
	}
	PluginModel = map[string]string{
		"401": "针对网站401认证进行暴力破解",
		//"PathCrack": "针对网站目录进行暴力破解",
		"JsFinder": "针对网站JS爬取链接",
		//"PortCrack": "针对诸如22，3306等服务端口，进行爆破密码",
	}

}
