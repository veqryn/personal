package main

import (
	"fmt"
	"log"

	"github.com/avct/uasurfer"
	"github.com/mssola/user_agent"
	"github.com/ua-parser/uap-go/uaparser"
)

var (
	userAgents = [][]string{
		// Random
		{"Chris", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36`},
		{"Linux", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.11`},
		{"Bot", `Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)`},
		{"Amazon Silk", `Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_3; en-us; Silk/1.1.0-80) AppleWebKit/533.16 (KHTML, like Gecko) Version/5.0 Safari/533.16 Silk-Accelerated=true`},
		{"Mac", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36`},
		{"Outlook App", `Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; Microsoft Outlook 16.0.8201`},
		{"Thunderbird", `Mozilla/5.0 (Windows; U; Windows NT 5.1; cs; rv:1.8.1.21) Gecko/20090302 Lightning/0.9 Thunderbird/2.0.0.21`},
		{"MS Office", `Microsoft Office/14.0 (Windows NT 5.1; Microsoft Outlook 14.0.4536; Pro; MSOffice 14)`},
		{"Lotus Notes", `Mozilla/4.0 (compatible; Lotus-Notes/5.0; Windows-NT)`},
		{"Google Proxy", `Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.0.7) Gecko/2009021910 Firefox/3.0.7 (via ggpht.com)`},
		{"Apple Mail", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/536.26.14 (KHTML, like Gecko)`},
		{"Yahoo Mobile", `YahooMobile/1.0 (mail; 3.0.5.1311380); (Linux; U; Android 4.0.3; htc_runnymede Build/ICE_CREAM_SANDWICH_MR1);`},
		{"Android", `DDG-Android-3.0.12`},
		{"Slack Mail", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_2) AppleWebKit/537.36 (KHTML, like Gecko) AtomShell/2.3.1 Chrome/52.0.2743.82 Electron/1.3.8 Safari/537.36 Slack_SSB/2.3.1`},

		// Webmail
		{"Maxthon",`Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/532.4 (KHTML, like Gecko) Maxthon/3.0.6.27 Safari/532.4`},
		{"QQ",`Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.80 Safari/537.36 QQBrowser/9.3.6874.400`},
		{"Sogou Explorer",`Mozilla/5.0 (Windows NT 6.1) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.802.30 Safari/535.1 SE 2.X MetaSr 1.0`},
		{"UC Browser",`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.155 UBrowser/5.4.5426.1034 Safari/537.36`},
		{"Yandex",`Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) YaBrowser/1.0.1084.5402 Chrome/19.0.1084.5402 Safari/536.5`},
		{"360",`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36 QIHU 360SE`},
		{"Chromium",`Mozilla/5.0 (X11; U; Linux x86_64; en-US) AppleWebKit/534.10 (KHTML, like Gecko) Ubuntu/10.10 Chromium/8.0.552.237 Chrome/8.0.552.237 Safari/534.10`},
		{"Brave",`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) brave/0.7.11 Chrome/47.0.2526.110 Brave/0.36.5 Safari/537.36`},
		{"Baidu",`Mozilla/5.0 (Windows; U; Windows NT 5.1; zh_CN) AppleWebKit/534.7 (KHTML, like Gecko) Chrome/7.0 baidubrowser/1.x Safari/534.7`},

		// Mobile
		{"Samsung",`Mozilla/5.0 (Linux; Android 7.0; SAMSUNG SM-G920P Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/5.0 Chrome/51.0.2704.106 Mobile Safari/537.36`},
		{"360",`Mozilla/5.0 (Linux; U; Android 4.3; zh-cn; GT-I9308 Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30; 360browser(securitypay,securityinstalled); 360(android,uppayplugin); 360 Aphone Browser (6.0.1)`},
		{"Android",`Mozilla/5.0 (Linux; U; Android 4.2; en-us; Nexus 10 Build/JVP15I) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Safari/534.30`},
		{"Android Webview",`Mozilla/5.0 (Linux; Android 7.1.1; Pixel XL Build/NOF26V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/56.0.2924.87 Mobile Safari/537.36 GSA/6.13.22.21.arm64`},
		{"Blackberry",`Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; nl) AppleWebKit/534.8+ (KHTML, like Gecko) Version/6.0.0.668 Mobile Safari/534.8+`},
		{"Chromium",`Mozilla/5.0 (Linux; Ubuntu 15.04 like Android 4.4) AppleWebKit/537.36 Chromium/55.0.2883.75 Mobile Safari/537.36`},
		{"Kindle",`Mozilla/5.0 (X11; U; Linux armv7l like Android; en-us) AppleWebKit/531.2+ (KHTML, like Gecko) Version/5.0 Safari/533.2+ Kindle/3.0+`},
		{"Opera Mini",`Mozilla/5.0 (iPhone; CPU iPhone OS 7_1_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) OPiOS/8.0.0.78129 Mobile/11D201 Safari/9537.53`},
		{"Opera Mini",`Opera/10.61 (J2ME/MIDP; Opera Mini/5.1.21219/19.999; en-US; rv:1.9.3a5) WebKit/534.5 Presto/2.6.30`},
		{"Silk",`Mozilla/5.0 (Linux; U; Android 4.4.3; de-de; KFTHWI Build/KTU84M) AppleWebKit/537.36 (KHTML, like Gecko) Silk/3.47 like Chrome/37.0.2026.117 Safari/537.36`},
		{"UC Mobile",`Mozilla/5.0 (Linux; Android 4.4.2; Q510s Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 Mobile UCBrowser/3.4.1.483`},
		{"Soguo mobile", `Mozilla/5.0 (Linux; Android 7.0; VIE-L29 Build/HUAWEIVIE-L29; wv) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.92 Mobile Safari/537.36 SogouMSE,SogouMobileBrowser/5.6.0`},

		// Desktop Applications
		{"Apple Mail",`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/536.26.14 (KHTML, like Gecko)`},
		{"Gmail",`Mozilla/5.0 (Windows; U; Windows NT 5.1; de; rv:1.9.0.7) Gecko/2009021910 Firefox/3.0.7 (via ggpht.com)`},
		{"Lotus Notes",`Mozilla/4.0 (compatible; Lotus-Notes/6.0; Windows-NT)`},
		{"Outlook 2007",`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727; MSOffice 12)`},
		{"Outlook 2010",`Microsoft Office/14.0 (Windows NT 5.1; Microsoft Outlook 14.0.4536; Pro; MSOffice 14)`},
		{"Outlook 2013",`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; Trident/6.0; Microsoft Outlook 15.0.4420)`},
		{"Outlook 2016",`Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 10.0; WOW64; Trident/8.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.30729; .NET CLR 3.5.30729; Microsoft Outlook 16.0.6366; ms-office; MSOffice 16)`},
		{"Thunderbird",`Mozilla/5.0 (X11; Linux i686; rv:7.0.1) Gecko/20110929 Thunderbird/7.0.1`},
		{"Thunderbird",`Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9.2.8) Gecko/20100802 Lightning/1.0b2 Thunderbird/3.1.2 ThunderBrowse/3.3.2`},
		{"Live Mail",`Outlook-Express/7.0 (MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; HPDTDF; .NET4.0C; BRI/2; AskTbLOL/5.12.5.17640; TmstmpExt)`},
		{"Yahoo",`YahooMobileMail/1.0 (Android Mail; 1.3.10) (supersonic;HTC;PC36100;2.3.5/GRJ90)`},

		// Bots
		{"",`Googlebot/2.1 (+http://www.google.com/bot.html)`},
		{"",`Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)`},
		{"",`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.75 Safari/537.36 Google Favicon`},
		{"",`Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)`},
		{"",`Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)`},
		{"",`Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)`},
		{"",`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36; 360Spider`},
		{"",`Mozilla/5.0 (compatible; YandexMetrika/2.0; +http://yandex.com/bots yabs01)`},
		{"",`Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)`},
		{"",`Mozilla/5.0 (compatible; DuckDuckGo-Favicons-Bot/1.0; +http://duckduckgo.com)`},
	}
)

func main() {
	for _, uas := range userAgents {
		fmt.Printf("%s\t\t%s\n", uas[0], uas[1])
		//mssola(uas)
		//fmt.Printf("\n\n")
		uap(uas[1])
		//fmt.Printf("\n\n")
		//uasurferF(uas)
		fmt.Printf("\n\n")
	}
}

func mssola(uas string) {
	// The "New" function will create a new UserAgent object and it will parse
	// the given string. If you need to parse more strings, you can re-use
	// this object and call: ua.Parse("another string")
	ua := user_agent.New(uas)

	fmt.Printf("%v\n", ua.Mobile())  // => false
	fmt.Printf("%v\n", ua.Bot())     // => false
	fmt.Printf("%v\n", ua.Mozilla()) // => "5.0"

	fmt.Printf("%v\n", ua.Platform()) // => "X11"
	fmt.Printf("%v\n", ua.OS())       // => "Linux x86_64"
	fmt.Printf("%v\n", ua.OSInfo())

	name, version := ua.Engine()
	fmt.Printf("%v\n", name)    // => "AppleWebKit"
	fmt.Printf("%v\n", version) // => "537.11"

	name, version = ua.Browser()
	fmt.Printf("%v\n", name)    // => "Chrome"
	fmt.Printf("%v\n", version) // => "23.0.1271.97"

	fmt.Printf("%v\n", ua.Localization())
}

func uap(uas string) {
	parser, err := uaparser.New("./vendor/github.com/ua-parser/uap-go/uap-core/regexes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client := parser.Parse(uas)

	fmt.Println(client.UserAgent.Family) // "Amazon Silk"
	// fmt.Println(client.UserAgent.Major, client.UserAgent.Minor, client.UserAgent.Patch)  // "1" "1" "0-80"
	fmt.Println(client.Os.Family) // Android
	//fmt.Println(client.Os.Major, client.Os.Minor, client.Os.Patch, client.Os.PatchMinor)
	fmt.Println(client.Device.Family) // "Kindle Fire"
	fmt.Println(client.Device.Brand)
	fmt.Println(client.Device.Model)
}

func uasurferF(uas string) {
	// Parse() returns all attributes, including returning the full UA string last
	fmt.Printf("%+v\n", uasurfer.Parse(uas))
}
