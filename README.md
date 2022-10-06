# pixiv-like2tgchan-bot

将Pixiv收藏夹的内容同步到Telegram频道,以便管理自己收藏的图片

更新时间跟随RSSHub缓存刷新时间,为2小时

每次最多获取30个收藏

[查看截图](https://s2.loli.net/2022/10/06/lc97ogXRYxJFkbj.jpg)

## 如何使用

1. clone项目到本地或服务器 

   ```bash
   git clone https://github.com/SuInk/pixiv-like2tgchan-bot.git
   ```

2. 配置文件 

   将config.go.example文件更名为config.go

   ```bash
   cd config
   cp config.go.example config.go
   ```

   填入

   > UseProxy = true //是否使用代理  
   > ProxyURL    = "http://127.0.0.1:7890"//代理网址  
   > RssURL      = "https://rsshub.app/pixiv/user/bookmarks/15288095"//需要订阅的PixivRSS地址  
   > TgBotToken  = "1234567892:AAG" //Telegram 机器人token  
   > ChatID      = "@SuInks" // 频道名  
   > RefreshTime = 120 // 定时刷新分钟, RSSHub默认缓存时间为2小时
3. 获取参数
  * 如果你用的是国外服务器`UseProxy`填`false`即可
  * `ProxyURL`为[Clash](https://github.com/Dreamacro/clash/releases)默认代理地址, 可以参考[教程](https://www.idcbuy.net/it/linux/2433.html)安装
  * `RssURL`在[RSSHub](https://docs.rsshub.app/social-media.html#pixiv)查看
  * `TgBotToken`在Telegram申请
  * `ChatID `为你的频道chatid或英文唯一ID
4. 运行

   ```bash
   go build main.go
   ./main
   ```

   

## Thanks

[RSSHub](https://github.com/DIYgod/RSSHub)

