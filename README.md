# 功能介绍
内网批量探测存活网段，端口扫描，端口指纹识别，web指纹识别，端口服务爆破，目录爆破，web漏洞扫描，基础认证爆破等功能  

# 使用说明
基础用法
```
.\yscan.exe -host 127.0.0.1
.\yscan.exe -host 192.168.1.1/24
.\yscan.exe -host www.example.com
.\yscan.exe -host http://www.example.com
.\yscan.exe -hf 1.txt
```
完整参数，是不是很简洁（后面更新在增加功能的基础上还会缩减参数数量）
```
  -dict
        是否对端口服务(ssh,ftp..)进行爆破，默认不爆破
  -hf string
        指定需扫描的IP/域名/url 文件路径
  -host string
        输入待扫描的IP/域名/url地址，格式支持：192.168.21.1/24, 192.168.21.1-255, 192.168.21.1-192.168.21.255, www.example.com,http://www.example.com
  -icmp string
        icmp扫描,探测内网某子网段存活，例：-icmp 192.168.*.1,-icmp 172.*.*.1
  -noping
        不进行ip存活探测
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
编译命令：
```
go build -ldflags="-s -w "  
```
参数不多，功能也不复杂~ 希望可以用最少的参数实现最多的功能可能性  
# 使用截图  
```
yscan.exe -host 127.0.0.1
```  
![](https://raw.githubusercontent.com/pikachu-cxy/yscan/main/images/009.png)   
```
快速探测内网存在哪些网段
yscan.exe -icmp 192.168.*.1
``` 
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/003.png)  
```
并不是通过端口判断指纹，下图爆破1521端口的mysql
yscan.exe -host 127.0.0.1 -port 1521 -dict
``` 
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/004.png)  
```
目前只针对url，提供了jsfinder，basic认证，目录爆破三种功能
yscan.exe -plugins list
```
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/005.png)  
```
yscan.exe -host http://192.168.2.150:8081 -plugins jsfinder
```
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/006.png)  
```
tomcat 401 basic认证爆破
yscan.exe -host http://192.168.2.150:8081/host-manager/html -plugins 401
```
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/008.png)
```
指定漏洞tag扫描探测
yscan.exe -host http://192.168.2.150:8081/ -s tomcat 
```
![](https://github.com/pikachu-cxy/yscan/blob/71a05f31413a109f0c9577122844cf40a7b79665/images/007.png)  
# 参考&引用的项目
[https://github.com/veo/vscan](https://github.com/veo/vscan)  
[https://github.com/praetorian-inc/fingerprintx](https://github.com/praetorian-inc/fingerprintx)  
[https://github.com/niudaii/crack](https://github.com/niudaii/crack)  
[https://github.com/projectdiscovery/httpx](https://github.com/projectdiscovery/httpx)  
[https://github.com/zan8in/afrog](https://github.com/zan8in/afrog)  
[https://github.com/lcvvvv/kscan](https://github.com/lcvvvv/kscan)
