# 功能介绍
内网批量探测存活网段，端口扫描，端口指纹识别，web指纹识别，端口服务爆破，目录爆破，web漏洞扫描，基础认证爆破等功能  

# 使用说明
基础用法
```
.\yscan.exe -host 127.0.0.1
.\yscan.exe -host 192.168.1/24
.\yscan.exe -host www.example.com
.\yscan.exe -host http://www.example.com
.\yscan.exe -hf 1.txt
```
完整参数，是不是很简洁（后面更新在增加功能的基础上还会缩减参数数量）
```
.\yscan.exe -h
  -dict
        是否对端口服务(ssh,ftp..)进行爆破，默认不爆破
  -hf string
        指定需扫描的IP/域名/url 文件路径
  -host string
        输入待扫描的IP/域名/url地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255, www.example.com,http://www.example.com
  -icmp string
        icmp扫描,探测内网某子网段存活，例：-icmp 192.168.*.1,-icmp 172.*.*.1
  -output string
        导出扫描结果到指定文件 (default "output.txt")
  -plugins string
        针对url进行测试，可指定使用某个插件(401,jsfinder...),-plugins list 查看程序目前支持的所有插件详情
  -pocs
        是否进行自动poc探测,默认不探测
  -port string
        输入需要扫描的端口,支持如下参数：full（全端口扫描）,top100,top1000 ,HttPorts(常见http端口）,1-65535（自定义端口 范围） (default "top100")
  -s string
        跳过指纹识别,对目标网址指定poc探测,例：-s shiro,seeyon,weblogic,thinkphp  
```
编辑命令：
```
go build -ldflags="-s -w "  
```
参数不多，功能也不复杂~ 希望可以用最少的参数实现最多的功能可能性  
# 使用截图  
![yscan.exe -host 127.0.0.1](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/002.png)  
![yscan.exe -icmp 192.168.*.1](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/003.png)  



