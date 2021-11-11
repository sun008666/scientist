// package models

// import (
//     "io/ioutil"

// "github.com/go-xweb/log"
// )

// type HtmlPageTable struct {
//     ID                int    `gorm:"column:id" json:"id"`
//     GameID            int    `gorm:"column:game_id" json:"game_id"`
//     GameName          string `gorm:"column:game_name" json:"game_name"`
//     EnName            string `gorm:"column:en_name" json:"en_name"`
//     UrlName           string `gorm:"column:url_name" json:"url_name"`
//     AndroidUrl        string `gorm:"column:android_url" json:"android_url"`
//     IosUrl            string `gorm:"column:ios_url" json:"ios_url"`
//     IsAndroidRedirect bool   `gorm:"column:is_android_redirect" json:"is_android_redirect"`
//     IsIosRedirect     bool   `gorm:"column:is_ios_redirect" json:"is_ios_redirect"`
//     HeadImg           string `gorm:"column:head_img" json:"head_img"`
//     BodyImg           string `gorm:"column:body_img" json:"body_img"`
//     AdImg             string `gorm:"column:ad_img" json:"ad_img"`
//     AdUrl             string `gorm:"column:ad_url" json:"ad_url"`
//     IsAdShow          bool   `gorm:"column:is_ad_show" json:"is_ad_show"`
//     // AndroidID         int                  `gorm:"column:android_id" json:"android_id"`
//     // IosID             int                  `gorm:"column:ios_id" json:"ios_id"`
//     AndroidList []AndriodPackageList `json:"android_list"`
//     IosList     []IosPackageList     `gorm:"column:ios_list" json:"ios_list"`
//     Close       bool                 `gorm:"column:close" json:"close"`
// }

// type AndriodPackageList struct {
//     PackageURL string `gorm:"column:package_url" json:"package_url"`
//     Url        string `gorm:"column:url" json:"url"`
//     Remark     string `gorm:"column:remark" json:"remark"`
// }

// type IosPackageList struct {
//     PackageURL string `gorm:"column:package_url" json:"package_url"`
//     Url        string `gorm:"column:url" json:"url"`
//     Remark     string `gorm:"column:remark" json:"remark"`
// }

// // CreateHTML 下载页面配置生成html
// func CreateHTML(html string, gameName string, isAndroidRedirect bool, isIosRedirect bool, isAdShow bool, adURL string, adImg string, androidURL string, headImg string, bodyImg string, iosURL string, gameID string) bool {

//     temp := ""

//     if gameID == "9" || gameID == "29" {

//         temp = `<!DOCTYPE html>
// 		<html><head>
// 				<meta http-equiv="content-type" content="text/html; charset=UTF-8">
// 				<meta charset="utf-8">
// 				<meta content="width=device-width,user-scalable=no" name="viewport">
// 				<meta name="HandheldFriendly" content="true">
// 				<meta http-equiv="x-rim-auto-match" content="none">
// 				<meta name="format-detection" content="telephone=no">
// 				<title>` + gameName + `</title>
// 				<meta name="keywords" content="` + gameName + `,中至麻将下载,中至四川麻将下载,闲来湖南麻将,,闲来上饶麻将,皮皮江西麻将,土豪金麻将,呱呱麻将,微乐江西棋牌,同城游上饶麻将,">
// 				<meta name="description" content="《` + gameName + `》下载官方网站，` + gameName + `是一款在江西非常流行的麻将游戏，因玩法易上手，耐玩性和趣味性都不错，所以在全国各地都比较流行。时尚简约的棋牌画面，真人实时对战玩法，陪伴您快乐每一天。">
// 				<meta name="author" content="中至游戏/" />
// 				<link rel="stylesheet" href="https://www.4417.com/css/common_layout.css">
// 				<link rel="stylesheet" href="https://www.4417.com/css/app_push.css">
// 				<link rel="stylesheet" href="css/main.css">
// 				<script src="https://static.xq5.com/3w4417com/js/jquery.js"></script>
// 				<script src="https://static.xq5.com/3w4417com/js/app_push.js"></script>
// 				<script>
// 		function is_weixin() {
// 			var ua = navigator.userAgent.toLowerCase();
// 			if (ua.match('micromessenger') == "micromessenger" ) {
// 				return true;
// 			} else {
// 				return false;
// 			}

// 		};

// 		function is_android() {
// 			var u = navigator.userAgent;
// 			var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
// 			if (isAndroid == true) {
// 				return true
// 			}
// 			return false
// 		}

// 		function is_ios(){
// 			var u = navigator.userAgent;
// 			var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 			if (isIOS == true) {
// 				return true
// 			}
// 			return false
// 		}
// 		`
//         // if isAndroidRedirect && !isIosRedirect {
//         // 	temp = temp + `$(document).ready(function(){
//         // 		if (is_weixin() && is_android()) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // } else if !isAndroidRedirect && isIosRedirect {
//         // 	temp = temp + `$(document).ready(
//         // 		function(){
//         // 			if (is_weixin() && is_ios()) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // } else {
//         // 	temp = temp + `$(document).ready(
//         // 		function(){
//         // 			if (is_weixin() && (is_ios() || is_android())) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // }
//         temp = temp + `
// 		</script>
// 			</head>
// 			<body>
// 				<div id="wrap">

// 						<div class="bg1">
// 						<img id='mask' style="float:right;position:absolute;top:0;display:none;" src="https://static.xq5.com/3w4417com/img/black.png" >
// 						<img id='guide' style="width:78%;float:right;position:absolute;top:0;display:none;left:20%" src="https://static.xq5.com/3w4417com/img/guide_android.png" >
// 						</div>
// 					<div class="bg1">
// 					`
//         if isAdShow {
//             temp = temp + `<a class="getpack" id="ad" href="` + adURL + `">
// 						<img style="`
//             if adImg == "" {
//                 temp = temp + "display:none"
//             }
//             temp = temp +
//                 `" src="` + adImg + `">
// 						</a>`
//         }

//         temp = temp +
//             `
// 					<a class="getpack" id="load" onclick="alert_msg()" href="` + androidURL + `">
// 					<img src="` + headImg + `">
// 					</a>
// 					<a class="getpack">
// 					<div class="bg1">
// 						<img src="` + bodyImg + `"></div>
// 						<img id='xinren' style="float:right;position:absolute;bottom:0;display:none;" src="https://www.4417.com/images/xin.gif" >
// 					<iframe id="ifrm" height="0" width="0"></iframe>
// 					<form id="form1" action="insertDate.action" method="post">
// 						<input id="dUrl" name="dUrl" value="" type="hidden">
// 						<input id="appid" name="appid" value="100286" type="hidden">
// 						<input id="packetid" name="packetid" value="11058" type="hidden">
// 						<input id="from" name="from" value="thnnly" type="hidden">
// 					</form>
// 					</a>
// 					</div>
// 					</div>

// 		<script>
// 		var _hmt = _hmt || [];
// 		(function() {
// 			var hm = document.createElement("script");
// 			hm.src = "https://hm.baidu.com/hm.js?29771296674c12a7bba9307344a0ed53";
// 			var s = document.getElementsByTagName("script")[0];
// 			s.parentNode.insertBefore(hm, s);
// 		})();
// 		function is_ios(){
// 			var u = navigator.userAgent;
// 			var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 			if (isIOS == true) {
// 				return true
// 			}
// 			return false
// 		}
// 		function alert_msg(){
// 			if (is_weixin() && (is_android())) {
// 				$('#guide').show();
// 				$('#mask').show();
// 			} else {
// 				if (is_ios()){
// 					setTimeout("alert('正在安装，请返回桌面查看!')",2000);
// 					// alert('执行了11111')

// 					// document.getElementById("xinren").style.display="block";

// 				}else if(!is_android()){
// 				 alert('请在浏览器设置中将【浏览器标识】或【浏览器UA标识】设置为您当前手机系统对应的选项或默认选项'+is_ios())
// 				}
// 			}
// 			}
// 			// window.onload=function(){
// 			// 		document.getElementById("xinren").onclick=function(){
// 			// 				var url="https://www.4417.com/XC_iOS_comzonstytmjhouse.mobileprovision";
// 			// 				window.open(url);
// 			// 				return true;
// 			// 		}
// 			// }

// 		</script>
// 		<script type="text/javascript">
// 		var u = navigator.userAgent;
// 		var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
// 		var isiOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 		if(!is_weixin()){
// 		        if (isiOS) {
// 		        	var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;'];
//                         			document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];
//         		}else {
//         		  var tempArrAndroid = "` + androidURL + `" ? "` + androidURL + `".split('|') : ['javascript:;'];
//                  document.getElementById("load").href=tempArrAndroid[Math.floor(Math.random()*tempArrAndroid.length)];

//         		}
// 		}else{
// 		  if(isiOS){
// 				var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;']
//                document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];
// 		  }else{
// 		  document.getElementById("load").href="javascript:void(0);"
// 		}
// 	}

// 		</script>
// 		<script>
// 		if(is_weixin() == false) {
// 			setTimeout(function () {
// 				document.querySelector('#load').click();//直接触发下载链接
// 			},1000)
// 		}
// 		</script>
// 		</body>
// 		<footer class="footera">
// 			<div style="width:100%;height:30px">
// 				<ul class="footerUl">
// 					<li>
// 						<a href="index.html">首页</a>
// 					</li>
// 					<li>
// 						<a href="custody.html">家长监护</a>
// 					</li>
// 					<li>
// 						<a href="prevent.html">防沉迷</a>
// 					</li>
// 					<li>
// 						<a href="login.html">登录</a>
// 					</li>
// 					<li>
// 						<a href="register.html">注册</a>
// 					</li>
// 				</ul>
// 			</div>
// 			<p class="beian">
// 				<span>赣ICP备15001426号-12</span>
// 				<span>
// 					<a href="images/增值电信业务经营许可证.jpg">赣B2—20150021</a>
// 				</span>
// 				<br>
// 				<span>
// 					<a href="images/网络文化经营许可证.jpg">赣网文〔2015〕1354-007号</a>
// 				</span>
// 				<span>赣公网安备36010802000031号</span>
// 			</p>
// 			<p class="copyright">中至集团 版权所有</p>
// 		</footer>
// 		</html>`

//     } else {
//         temp = `<!DOCTYPE html>
// 		<html><head>
// 				<meta http-equiv="content-type" content="text/html; charset=UTF-8">
// 				<meta charset="utf-8">
// 				<meta content="width=device-width,user-scalable=no" name="viewport">
// 				<meta name="HandheldFriendly" content="true">
// 				<meta http-equiv="x-rim-auto-match" content="none">
// 				<meta name="format-detection" content="telephone=no">
// 				<title>` + gameName + `</title>
// 				<meta name="keywords" content="` + gameName + `,中至麻将下载,中至四川麻将下载,闲来湖南麻将,,闲来上饶麻将,皮皮江西麻将,土豪金麻将,呱呱麻将,微乐江西棋牌,同城游上饶麻将,">
// 				<meta name="description" content="《` + gameName + `》下载官方网站，` + gameName + `是一款在江西非常流行的麻将游戏，因玩法易上手，耐玩性和趣味性都不错，所以在全国各地都比较流行。时尚简约的棋牌画面，真人实时对战玩法，陪伴您快乐每一天。">
// 				<meta name="author" content="中至游戏/" />
// 				<link rel="stylesheet" href="https://www.4417.com/css/common_layout.css">
// 				<link rel="stylesheet" href="https://www.4417.com/css/app_push.css">
// 				<link rel="stylesheet" href="css/main.css">
// 				<script src="https://static.xq5.com/3w4417com/js/jquery.js"></script>
// 				<script src="https://static.xq5.com/3w4417com/js/app_push.js"></script>
// 				<script>
// 		function is_weixin() {
// 			var ua = navigator.userAgent.toLowerCase();
// 			if (ua.match('micromessenger') == "micromessenger" ) {
// 				return true;
// 			} else {
// 				return false;
// 			}

// 		};

// 		function is_android() {
// 			var u = navigator.userAgent;
// 			var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
// 			if (isAndroid == true) {
// 				return true
// 			}
// 			return false
// 		}

// 		function is_ios(){
// 			var u = navigator.userAgent;
// 			var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 			if (isIOS == true) {
// 				return true
// 			}
// 			return false
// 		}
// 		`
//         // if isAndroidRedirect && !isIosRedirect {
//         // 	temp = temp + `$(document).ready(function(){
//         // 		if (is_weixin() && is_android()) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // } else if !isAndroidRedirect && isIosRedirect {
//         // 	temp = temp + `$(document).ready(
//         // 		function(){
//         // 			if (is_weixin() && is_ios()) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // } else {
//         // 	temp = temp + `$(document).ready(
//         // 		function(){
//         // 			if (is_weixin() && (is_ios() || is_android())) {
//         // 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
//         // 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
//         // 				}
//         // 				$('#guide').show();
//         // 				$('#mask').show();}
//         // 		});`
//         // }
//         temp = temp + `
// 		</script>
// 			</head>
// 			<body>
// 				<div id="wrap">

// 						<div class="bg1">
// 						<img id='mask' style="float:right;position:absolute;top:0;display:none;" src="https://static.xq5.com/3w4417com/img/black.png" >
// 						<img id='guide' style="width:78%;float:right;position:absolute;top:0;display:none;left:20%" src="https://static.xq5.com/3w4417com/img/guide_android.png" >
// 						</div>
// 					<div class="bg1">
// 					`
//         if isAdShow {
//             temp = temp + `<a class="getpack" id="ad" href="` + adURL + `">
// 						<img style="`
//             if adImg == "" {
//                 temp = temp + "display:none"
//             }
//             temp = temp +
//                 `" src="` + adImg + `">
// 						</a>`
//         }

//         temp = temp +
//             `
// 					<a class="getpack" id="load" onclick="alert_msg()" href="` + androidURL + `">
// 					<img src="` + headImg + `">
// 					</a>
// 					<a class="getpack">
// 					<div class="bg1">
// 						<img src="` + bodyImg + `"></div>
// 						<img id='xinren' style="float:right;position:absolute;bottom:0;display:none;" src="https://www.4417.com/images/xin.gif" >
// 					<iframe id="ifrm" height="0" width="0"></iframe>
// 					<form id="form1" action="insertDate.action" method="post">
// 						<input id="dUrl" name="dUrl" value="" type="hidden">
// 						<input id="appid" name="appid" value="100286" type="hidden">
// 						<input id="packetid" name="packetid" value="11058" type="hidden">
// 						<input id="from" name="from" value="thnnly" type="hidden">
// 					</form>
// 					</a>
// 					</div>
// 					</div>

// 		<script>
// 		var _hmt = _hmt || [];
// 		(function() {
// 			var hm = document.createElement("script");
// 			hm.src = "https://hm.baidu.com/hm.js?29771296674c12a7bba9307344a0ed53";
// 			var s = document.getElementsByTagName("script")[0];
// 			s.parentNode.insertBefore(hm, s);
// 		})();
// 		function is_ios(){
// 			var u = navigator.userAgent;
// 			var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 			if (isIOS == true) {
// 				return true
// 			}
// 			return false
// 		}
// 		function alert_msg(){
// 			if (is_weixin() && (is_android())) {
// 				$('#guide').show();
// 				$('#mask').show();
// 			} else {
// 				if (is_ios()){
// 					// setTimeout("alert('正在安装，请返回桌面查看!')",2000);
// 					// alert('执行了11111')

// 					document.getElementById("xinren").style.display="block";

// 				}else if(!is_android()){
//                 alert('请在浏览器设置中将【浏览器标识】或【浏览器UA标识】设置为您当前手机系统对应的选项或默认选项'+is_ios())
//             }
// 			}
// 			}
// 			window.onload=function(){
//                   document.getElementById("xinren").onclick=function(){
//                   if((is_weixin())){
//                   alert('请点击右上角... 选择【在Safari中打开】后在Safari中点击该按钮')
//                   }else{
//                     var url="https://www.4417.com/XC_iOS_comzonstytmjhouse.mobileprovision";
//                                window.open(url);
//                                return true;
//                   }
//                }
//             }
// 			if(is_ios() && is_weixin()){
// 				document.getElementById("xinren").style.display="block";
// 				}

// 		</script>
// 		<script type="text/javascript">
// 		var u = navigator.userAgent;
// 		var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
// 		var isiOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
// 		if(!is_weixin()){
//         if (isiOS) {
//         var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;'];
//         document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];

// 		}else {
// 		 var tempArrAndroid = "` + androidURL + `" ? "` + androidURL + `".split('|') : ['javascript:;'];
//         document.getElementById("load").href=tempArrAndroid[Math.floor(Math.random()*tempArrAndroid.length)];

// 		}
// 		}else{
// 		if(isiOS){
// 			var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;']
//         			document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];
// 		}else{
// 			document.getElementById("load").href="javascript:void(0);"
// 		}

// 		}

// 		</script>
// 		<script>
// 		if(is_weixin() == false) {
// 			setTimeout(function () {
// 				document.querySelector('#load').click();//直接触发下载链接
// 			},1000)
// 		}
// 		</script>
// 		</body>
// 		<footer class="footera">
// 			<div style="width:100%;height:30px">
// 				<ul class="footerUl">
// 					<li>
// 						<a href="index.html">首页</a>
// 					</li>
// 					<li>
// 						<a href="custody.html">家长监护</a>
// 					</li>
// 					<li>
// 						<a href="prevent.html">防沉迷</a>
// 					</li>
// 					<li>
// 						<a href="login.html">登录</a>
// 					</li>
// 					<li>
// 						<a href="register.html">注册</a>
// 					</li>
// 				</ul>
// 			</div>
// 			<p class="beian">
// 				<span>赣ICP备15001426号-12</span>
// 				<span>
// 					<a href="images/增值电信业务经营许可证.jpg">赣B2—20150021</a>
// 				</span>
// 				<br>
// 				<span>
// 					<a href="images/网络文化经营许可证.jpg">赣网文〔2015〕1354-007号</a>
// 				</span>
// 				<span>赣公网安备36010802000031号</span>
// 			</p>
// 			<p class="copyright">中至集团 版权所有</p>
// 		</footer>
// 		</html>`
//     }

//     var d1 = []byte(temp)
//     if err := ioutil.WriteFile(html, d1, 0666); err != nil {
//         log.Errorf("CreateHTML: err: %v\n", err.Error())
//         return false
//     }
//     return true
// }

// type OnDownloadPageOpRequest struct {
//     GameID int  `json:"game_id" binding:"required"`
//     Close  bool `json:"close"`
// }
