# Oyster-IoT
English | [简体中文](./README-zh.md)
<p>
	<p align="center">
		<img src="https://img.gejiba.com/images/375278c7eb92bab1d7ba09f27e2fe8b4.png">
	</p>
	<p align="center">
		<font size=6 face="宋体">打造海洋养殖专业物联平台</font>
	<p>
</p>
<p align="center">
<img alt="Go" src="https://img.shields.io/badge/Go-1.7%2B-blue">
<img alt="Mysql" src="https://img.shields.io/badge/Mysql-5.7%2B-brightgreen">
<img alt="Redis" src="https://img.shields.io/badge/Redis-6.2%2B-yellowgreen">
<img alt="influxDB" src="https://img.shields.io/badge/influxDB-2.3%2B-orange">
<img alt="EMQX" src="https://img.shields.io/badge/EMQX-4.4+%2B-yellow">
<img alt="license" src="https://img.shields.io/badge/license-GPL-lightgrey">
</p>

> Oyster-IoT: 采用Beego框架，结合Mysql+Redis+influxDB+EMQX开发的物联网管理平台。当前支持单机部署，硬件探测设备使用MQTT交互。

**演示平台:** https://oyster-iot.cloud/  
**微信公众号:** 扫描添加，联系获取体验账号

<img src="https://img.gejiba.com/images/7ba3b6926304f02ca20edb6ef4dec764.jpg" alt="image-20220610100236465" />

## 环境依赖
Mysql v5.7+

Redis v6.2.6

influxDB v2.3.0

EMQX v4.4.3

Note: 请选择合适的版本进行安装,否则服务启动失败。

## 配置文件
```bash
# Configure(必须) mysql redis connection infomation 
# Option(可选): qiniu influxdb
# Note: 
# 七牛云配置不配，则不能使用视频监控功能
# If the qiniu configuration is not compatible, the video surveillance function cannot be used
# influxdb配置不配，则默认使用mysql存储设备上传的数据
# If the configuration is not supported, use mysql to store data uploaded from the device by default
# See Beego Configuration for more app.conf
# 关于更多 app.conf 可以查看beego配置

vim ./conf/app.conf

# 配置设备接入EMQX信息和influxDB
# Configure(必须) MQTT EMQX connection infomation 
# Option(可选): influxdb
# Note:
# influxdb配置不配，则默认使用mysql存储设备上传的数据
# If the configuration is not supported, use mysql to store data uploaded from the device by default

vim ./devaccess/config.ini

```
## 编译运行
```bash
# check go verison
go verison: go1.17.8

# clone the project
git clone https://github.com/ptonlix/oyster-iot

# compile and run
make build
./oyster

# Web request
1.Apifox
2.oyster-iot-admin-vue
```

## 一.前言介绍

​随着生蚝的快速普及，各地对生蚝的的需求量日益增长，近海养殖中，生蚝养殖占据了一席之地。以广西为例，钦州港和铁山港均有上千亩的生蚝养殖厂。

​生蚝主要是以筏式养殖为主，养殖地点应当选择在风浪较小，水质稳定，没有工业污染，底质为泥质或泥沙质的海区，水深保证在8m以上，水温变化幅度较小，夏季不超过30℃，表层流速在0.3-0.5m/s之间，比重在1.008-1.02之间。

​影响生蚝生长的两个主要环境因素是，海水温度和盐度。如：

1. 近江牡蛎，其生长适宜温度范围为10-33℃，适盐范围为5-25‰。
2. 长牡蛎，其生存温度范围为零下3℃至32℃，生长适宜温度范围为5-28℃，适盐范围为10-37‰，尤以20-30‰更为合适。
3. 密鳞牡蛎，其生存温度范围与长牡蛎相同，均为零下3℃至32℃，但是适盐范围较窄，为27-34‰。
4. 褶牡蛎，其生长适宜温度范围为6-25℃，对盐度的适应范围较广（与近江牡蛎接近），内湾以及近海均有分布。
5. 大连湾牡蛎，其生长适宜水温范围为6-25℃，适盐范围较窄，为25-34‰。

​海水的盐度，决定了生蚝的鲜美程度。

​和其他贝类动物一样，生蚝为了对抗海水的盐分，必须在体内积累足够的氨基酸才能得以生存。而这些氨基酸中最主要的便是鲜味物质的代表——谷氨酸，所以不同海域的生蚝鲜味浓淡才会有所不同。

​海水的温度也会影响生蚝的滋味。

​生蚝有一个独特的技能——能够根据周围环境自由转换性别。在温暖的水域，食物丰盛，生蚝通常会变成肥腴鲜美的雌性，此时的生蚝肉质柔滑饱满，体内的蚝卵也同样提供了鲜美的风味；在冷凉的水域，生蚝生长速度放缓，更容易积累风味物质，通常会变为雄性生蚝，肉质清瘦爽脆。

不同的生蚝种类，需要合适温度和盐度，所以生蚝养殖非常清楚海水的温度和盐度，以帮助养殖户及时了解生蚝的生长环境，判断生蚝的生长情况，更能预测生蚝长成的质量。

#### 传统的测量方式

传统的探测海水温度盐度的方式，是采用人工测量的方式，使用光学盐度计和海水温度计分别测量。

​<img src="https://img.gejiba.com/images/8ead700f79e1a17f1f3657ec22994052.png" alt="image-20220610100236465" width="180" height="320" />&nbsp; &nbsp; &nbsp; <img src="https://img.gejiba.com/images/416baebb474b629caaedb7bd68a551e8.png" alt="image-20220610100518325" width="180" height="320" />

光学盐度计&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; 海水温度计

人工测量费时费力，上千亩的生蚝基地，如果纯人工测试量，每天花费的人工成本巨大，测量结束还需要手动记录，录入系统，而且无法形成系统性的报表，以进一步的分析预测。

#### Oyster智能生蚝养殖系统

针对传统的测量方式的弊端，Oyster智能生蚝养殖系统，采用海水温度盐度智能一体化探测，部署完成后，无需人工介入，智能探测。通过5G系统实时将数据传输到Oyster数据中台，进行分析整合，输出温度盐度的周期数据和实时数据。

Oyster数据中台采用时序数据库，针对海水的温度盐度进行大数据分析，同时融合了深度学习，结合历史数据和生蚝长成效果对生蚝的成长进行预测，为养殖户提供决策依据。

同时根据部署环境的条件，还提供太阳能供电能力和传统电线供电两种方式，为在不同条件下部署提供保障。

## 二.系统篇

![Oyster-IOT平台架构图](https://img.gejiba.com/images/1d08f2b79fa4ecc29b9668d42382a6c9.png)

Oyster- IOT平台主要分为：   API接口层、业务层、数据层、设备接入层  这四个部分。

1.API接口层：主要提供管理平台和小程序等相关调用接口
  	
2.业务层：目前平台主要分为 业务管理、资产管理、应用管理、自动化系统、可视化管理 这五大业务子系统
  	
3.数据层：

- InfluxDB时序数据库，用于结构化存储IOT设备上传的海量数据

- Mysql数据库，用于存储业务系统数据

- Redis缓存数据库，用户存储状态数据，加快访问

- TensorFlow深度学习框架，用户分析IOT数据，预测结果

4.设备接入层： 为各种IOT设备提供接入能力，兼容各种设备接入。

## 三.智能硬件篇

![Oyster-IOT智能网关硬件架构图](https://img.gejiba.com/images/eabc1e3b45326b4d4ee15e329c22d279.png)

​	  	

​Oyster- IOT智能网关主要分为： 主控、通信模块、GNSS定位系统、电源模块 等四个部分

​Oyster- IOT智能网关作为智能探测系统总控部分，与其它传感器系统进行星型组网，构成整个探测系统。

![Oyster-IOT硬件组网](https://img.gejiba.com/images/d370bbd2485e4ba98fcfb50e73d1cdc2.png)

​探测器与网关目前主要采用WI-FI方式进行通信，可以在生蚝养殖基地进行大范围的组网，解决生蚝养殖面积大，组网难度大的问题。

​Oyster-IOT主要采用MQTT协议与平台进行通信。

## 四.小程序篇

​Oyster智能生蚝养殖系统，目前主要使用微信小程序作为客户端，提供养殖户进行使用。

<img src="https://img.gejiba.com/images/5ec4ad3bb3e0cd9c96be99ec1e383ff7.jpg" alt="WechatIMG135" width="240" height="480" />

请通过微信公众号,可以访问微信小程序和获取体验账号

## 五.管理平台

​Oyster智能生蚝养殖系统，后台管理平台由Oyster智能物联管理平台统一管理。

目前Oyster智能物联管理平台，数据大屏、用户管理、日志管理等部分组成

详情请看: oyster-iot-admin-vue 项目

[github地址](https://github.com/ptonlix/oyster-iot-admin-vue)
[gitlab地址](https://gitee.com/ptonlix/oyster-iot-admin-vue)

​		

## 六.结尾

以上就是Oyster智能生蚝养殖系统整体介绍了，目前该系统主要应用在生蚝养殖行业，后续可以拓展到其它养殖行业，或者更大胆一点，可以兼容更广泛的物联网应用。

欢迎大家联系交流咨询，添加公众号一起讨论～

<p align="center">
  <b>SPONSORED BY</b>
</p>
<p align="center">
   <a href="https://www.gogeek.com.cn/" title="gogeek" target="_blank">
      <img height="200px" src="https://img.gejiba.com/images/96b6d150bd758b13d66aec66cb18044e.jpg" title="gogeek">
   </a>
</p>