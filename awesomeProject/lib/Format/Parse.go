package Format

import (
	"awesomeProject/lib/Plugins"
	"fmt"
	"strings"
)

func ParsePlugins(plugins string, ps Plugins.PluginService) {

	if plugins == "list" {
		println("plugins list: ")
		for name, des := range Plugins.PluginModel {
			println(fmt.Sprintf("   %s : %s", name, des))
		}
		return
	}

	if IsUrl(ps.Host) {
		lists := strings.Split(plugins, ",")

		for _, list := range lists {
			Plugins.Plugin[list](ps)
		}
	} else if ps.Host != "" {
		println("Plugins 功能目前只支持url! ,暂不支持IP/域名~")
		return
	}
}

func PluginScan(host string) {

}
