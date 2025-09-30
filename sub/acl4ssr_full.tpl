mixed-port: 7890
allow-lan: true
mode: Rule
log-level: info
external-controller: :9090
proxies:
proxy-groups:
  - name: 🚀 节点选择
    type: select
    proxies:
      - ♻️ 自动选择
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 🚀 手动切换
    type: select
    proxies:
      - ett1tdyb
  - name: ♻️ 自动选择
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - ett1tdyb
  - name: 📲 电报消息
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
      - 🇸🇬 狮城节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 💬 Ai平台
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
      - 🇸🇬 狮城节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 📹 油管视频
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
      - 🇸🇬 狮城节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 🎥 奈飞视频
    type: select
    proxies:
      - 🎥 奈飞节点
      - 🚀 节点选择
      - ♻️ 自动选择
      - 🇸🇬 狮城节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 📺 巴哈姆特
    type: select
    proxies:
      - 🇨🇳 台湾节点
      - 🚀 节点选择
      - 🚀 手动切换
      - DIRECT
  - name: 📺 哔哩哔哩
    type: select
    proxies:
      - 🎯 全球直连
      - 🇨🇳 台湾节点
      - 🇭🇰 香港节点
  - name: 🌍 国外媒体
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
      - DIRECT
  - name: 🌏 国内媒体
    type: select
    proxies:
      - DIRECT
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🚀 手动切换
  - name: 📢 谷歌FCM
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: Ⓜ️ 微软Bing
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: Ⓜ️ 微软云盘
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: Ⓜ️ 微软服务
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: 🍎 苹果服务
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: 🎮 游戏平台
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - 🇺🇲 美国节点
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: 🎶 网易音乐
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - ♻️ 自动选择
  - name: 🎯 全球直连
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
      - ♻️ 自动选择
  - name: 🛑 广告拦截
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🍃 应用净化
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🐟 漏网之鱼
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
      - DIRECT
      - 🇭🇰 香港节点
      - 🇨🇳 台湾节点
      - 🇸🇬 狮城节点
      - 🇯🇵 日本节点
      - 🇺🇲 美国节点
      - 🇰🇷 韩国节点
      - 🚀 手动切换
  - name: 🇭🇰 香港节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
  - name: 🇯🇵 日本节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
  - name: 🇺🇲 美国节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 150
    proxies:
      - DIRECT
  - name: 🇨🇳 台湾节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
  - name: 🇸🇬 狮城节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
  - name: 🇰🇷 韩国节点
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
  - name: 🎥 奈飞节点
    type: select
    proxies:
      - DIRECT
rules:
  - DOMAIN-SUFFIX,acl4.ssr,🎯 全球直连
  - DOMAIN-SUFFIX,ip6-localhost,🎯 全球直连
  - DOMAIN-SUFFIX,ip6-loopback,🎯 全球直连
  - DOMAIN-SUFFIX,lan,🎯 全球直连
  - DOMAIN-SUFFIX,local,🎯 全球直连
  - DOMAIN-SUFFIX,localhost,🎯 全球直连
  - IP-CIDR,0.0.0.0/8,🎯 全球直连,no-resolve
  - IP-CIDR,10.0.0.0/8,🎯 全球直连,no-resolve
  - IP-CIDR,100.64.0.0/10,🎯 全球直连,no-resolve
  - IP-CIDR,127.0.0.0/8,🎯 全球直连,no-resolve
  - IP-CIDR,172.16.0.0/12,🎯 全球直连,no-resolve
  - IP-CIDR,192.168.0.0/16,🎯 全球直连,no-resolve
  - IP-CIDR,198.18.0.0/16,🎯 全球直连,no-resolve
  - IP-CIDR,224.0.0.0/4,🎯 全球直连,no-resolve
  - IP-CIDR6,::1/128,🎯 全球直连,no-resolve
  - IP-CIDR6,fc00::/7,🎯 全球直连,no-resolve
  - IP-CIDR6,fe80::/10,🎯 全球直连,no-resolve
  - IP-CIDR6,fd00::/8,🎯 全球直连,no-resolve
  - DOMAIN,instant.arubanetworks.com,🎯 全球直连
  - DOMAIN,setmeup.arubanetworks.com,🎯 全球直连
  - DOMAIN,router.asus.com,🎯 全球直连
  - DOMAIN,www.asusrouter.com,🎯 全球直连
  - DOMAIN-SUFFIX,hiwifi.com,🎯 全球直连
  - DOMAIN-SUFFIX,leike.cc,🎯 全球直连
  - DOMAIN-SUFFIX,miwifi.com,🎯 全球直连
  - DOMAIN-SUFFIX,my.router,🎯 全球直连
  - DOMAIN-SUFFIX,p.to,🎯 全球直连
  - DOMAIN-SUFFIX,peiluyou.com,🎯 全球直连
  - DOMAIN-SUFFIX,phicomm.me,🎯 全球直连
  - DOMAIN-SUFFIX,router.ctc,🎯 全球直连
  - DOMAIN-SUFFIX,routerlogin.com,🎯 全球直连
  - DOMAIN-SUFFIX,tendawifi.com,🎯 全球直连
  - DOMAIN-SUFFIX,zte.home,🎯 全球直连
  - DOMAIN-SUFFIX,tplogin.cn,🎯 全球直连
  - DOMAIN-SUFFIX,wifi.cmcc,🎯 全球直连
  - DOMAIN-SUFFIX,ol.epicgames.com,🎯 全球直连
  - DOMAIN-SUFFIX,dizhensubao.getui.com,🎯 全球直连
  - DOMAIN,dl.google.com,🎯 全球直连
  - DOMAIN-SUFFIX,googletraveladservices.com,🎯 全球直连
  - DOMAIN-SUFFIX,tracking-protection.cdn.mozilla.net,🎯 全球直连
  - DOMAIN,origin-a.akamaihd.net,🎯 全球直连
  - DOMAIN,fairplay.l.qq.com,🎯 全球直连
  - DOMAIN,livew.l.qq.com,🎯 全球直连
  - DOMAIN,vd.l.qq.com,🎯 全球直连
  - DOMAIN,errlog.umeng.com,🎯 全球直连
  - DOMAIN,msg.umeng.com,🎯 全球直连
  - DOMAIN,msg.umengcloud.com,🎯 全球直连
  - DOMAIN,tracking.miui.com,🎯 全球直连
  - DOMAIN,app.adjust.com,🎯 全球直连
  - DOMAIN,bdtj.tagtic.cn,🎯 全球直连
  - DOMAIN,rewards.hypixel.net,🎯 全球直连
  - DOMAIN-SUFFIX,koodomobile.com,🎯 全球直连
  - DOMAIN-SUFFIX,koodomobile.ca,🎯 全球直连
  - DOMAIN-SUFFIX,synology.me,🎯 全球直连
  - DOMAIN-SUFFIX,DiskStation.me,🎯 全球直连
  - DOMAIN-SUFFIX,i234.me,🎯 全球直连
  - DOMAIN-SUFFIX,myDS.me,🎯 全球直连
  - DOMAIN-SUFFIX,DSCloud.biz,🎯 全球直连
  - DOMAIN-SUFFIX,DSCloud.me,🎯 全球直连
  - DOMAIN-SUFFIX,DSCloud.mobi,🎯 全球直连
  - DOMAIN-SUFFIX,DSmyNAS.com,🎯 全球直连
  - DOMAIN-SUFFIX,DSmyNAS.net,🎯 全球直连
  - DOMAIN-SUFFIX,DSmyNAS.org,🎯 全球直连
  - DOMAIN-SUFFIX,FamilyDS.com,🎯 全球直连
  - DOMAIN-SUFFIX,FamilyDS.net,🎯 全球直连
  - DOMAIN-SUFFIX,FamilyDS.org,🎯 全球直连
  - DOMAIN-KEYWORD,admarvel,🛑 广告拦截
  - DOMAIN-KEYWORD,admaster,🛑 广告拦截
  - DOMAIN-KEYWORD,adsage,🛑 广告拦截
  - DOMAIN-SUFFIX,wan.2345.cn,🍃 应用净化
  - DOMAIN-SUFFIX,zhushou.2345.cn,🍃 应用净化
  - DOMAIN-SUFFIX,3600.com,🍃 应用净化
  - DOMAIN-SUFFIX,gamebox.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,jiagu.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,kuaikan.netmon.360safe.com,🍃 应用净化
  - DOMAIN-SUFFIX,leak.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,lianmeng.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,pub.se.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,s.so.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,shouji.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,soft.data.weather.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,stat.360safe.com,🍃 应用净化
  - DOMAIN-SUFFIX,stat.m.360.cn,🍃 应用净化
  - DOMAIN-SUFFIX,redirect.simba.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,rj.m.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,sdkinit.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,show.re.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,simaba.m.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,simaba.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,srd.simba.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,strip.taobaocdn.com,🍃 应用净化
  - DOMAIN-SUFFIX,tns.simba.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,tyh.taobao.com,🍃 应用净化
  - DOMAIN-SUFFIX,userimg.qunar.com,🍃 应用净化
  - DOMAIN-SUFFIX,yiliao.hupan.com,🍃 应用净化
  - DOMAIN-SUFFIX,3dns-2.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,3dns-3.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,activate-sea.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,activate-sjc0.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,activate.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,adobe-dns-2.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,adobe-dns-3.adobe.com,🍃 应用净化
  - DOMAIN-SUFFIX,adobe-dns.adobe.com,🍃 应用净化