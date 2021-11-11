package models

import (
	"io/ioutil"
	"strings"
	"zonst/qipai/api/configapisrv/config"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"

	"github.com/go-xweb/log"
)

type HtmlPageTable struct {
	ID                int    `gorm:"column:id" json:"id"`
	GameID            int    `gorm:"column:game_id" json:"game_id"`
	GameName          string `gorm:"column:game_name" json:"game_name"`
	EnName            string `gorm:"column:en_name" json:"en_name"`
	UrlName           string `gorm:"column:url_name" json:"url_name"`
	AndroidUrl        string `gorm:"column:android_url" json:"android_url"`
	IosUrl            string `gorm:"column:ios_url" json:"ios_url"`
	IsAndroidRedirect bool   `gorm:"column:is_android_redirect" json:"is_android_redirect"`
	IsIosRedirect     bool   `gorm:"column:is_ios_redirect" json:"is_ios_redirect"`
	HeadImg           string `gorm:"column:head_img" json:"head_img"`
	BodyImg           string `gorm:"column:body_img" json:"body_img"`
	AdImg             string `gorm:"column:ad_img" json:"ad_img"`
	AdUrl             string `gorm:"column:ad_url" json:"ad_url"`
	IsAdShow          bool   `gorm:"column:is_ad_show" json:"is_ad_show"`
	// AndroidID         int                  `gorm:"column:android_id" json:"android_id"`
	// IosID             int                  `gorm:"column:ios_id" json:"ios_id"`
	AndroidList []AndriodPackageList `json:"android_list"`
	IosList     []IosPackageList     `gorm:"column:ios_list" json:"ios_list"`
	Close       bool                 `gorm:"column:close" json:"close"`
}

type AndriodPackageList struct {
	PackageURL string `gorm:"column:package_url" json:"package_url"`
	Url        string `gorm:"column:url" json:"url"`
	Remark     string `gorm:"column:remark" json:"remark"`
}

type IosPackageList struct {
	PackageURL string `gorm:"column:package_url" json:"package_url"`
	Url        string `gorm:"column:url" json:"url"`
	Remark     string `gorm:"column:remark" json:"remark"`
}

// CreateHTML 下载页面配置生成html
func CreateHTML(c *gin.Context, html string, gameName string, isAndroidRedirect bool, isIosRedirect bool, isAdShow bool, adURL string, adImg string, androidURL string, headImg string, bodyImg string, iosURL string, gameID, htmlURL string) bool {
	flag := strings.Contains(htmlURL, "4417.com")
	cfg := c.MustGet("config").(*config.Config)
	staticUrl := cfg.StaticURL
	licenceNum1 := cfg.LicenceNum1
	licenceNum2 := cfg.LicenceNum2
	pageIndex := cfg.PageIndex
	temp := ""
	temp = `<!DOCTYPE html>
		<html><head>
				<meta http-equiv="content-type" content="text/html; charset=UTF-8">
				<meta charset="utf-8">
				<meta content="width=device-width,user-scalable=no" name="viewport">
				<meta name="HandheldFriendly" content="true">
				<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
                <meta http-equiv="Pragma" content="no-cache" />
                <meta http-equiv="Expires" content="0" />
				<meta http-equiv="x-rim-auto-match" content="none">
				<meta name="format-detection" content="telephone=no">
				<title>` + gameName + `</title>
                `
	if flag {
		temp += `<meta name="keywords" content="` + gameName + `,中至麻将下载,中至四川麻将下载,闲来湖南麻将,,闲来上饶麻将,皮皮江西麻将,土豪金麻将,呱呱麻将,微乐江西棋牌,同城游上饶麻将,">
				<meta name="description" content="《` + gameName + `》下载官方网站，` + gameName + `是一款在江西非常流行的麻将游戏，因玩法易上手，耐玩性和趣味性都不错，所以在全国各地都比较流行。时尚简约的棋牌画面，真人实时对战玩法，陪伴您快乐每一天。">
				<meta name="author" content="中至游戏/" />`
	}
	temp +=
		`
				<link rel="stylesheet" href="` + htmlURL + `js/notchange/common_layout.css">
				<link rel="stylesheet" href="` + htmlURL + `js/notchange/app_push.css">
				<link rel="stylesheet" href="js/notchange/main.css">
				    <style>

        #load{
            position: relative;
        }

        .chuangyu-icon{
         margin: 0 auto;
          padding: 20px 0;
        width: 100px;
        height: auto;

       }
       .chuangyu-icon.hide{
         display: none;
       }
       .chuangyu-icon img{
       display: block;
       width: 100%;
       max-width: 100%;

       }


    </style>

				<script src="` + staticUrl + `/3w4417com/js/jquery.js"></script>
				<script src="` + staticUrl + `/3w4417com/js/app_push.js"></script>
				<script>
		// function is_weixin() {
		// 	var ua = navigator.userAgent.toLowerCase();
		// 	if (ua.match('micromessenger') == "micromessenger" ) {
		// 		return true;
		// 	} else {
		// 		return false;
		// 	}
        //
		// };


		function is_ios(){
			// var u = navigator.userAgent;
			// var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
			// if (isIOS == true) {
			// 	return true
			// }
			// return false
		 var str = navigator.userAgent.toLowerCase();
		 var ver = str.match(/cpu iphone os (.*?) like mac os/) ? str.match(/cpu iphone os (.*?) like mac os/) : str.match(/cpu os (.*?) like mac os/) ? str.match(/cpu os (.*?) like mac os/) : str.match(/intel mac os x (.*?)/);
		 return ver
		}
	function getSearchParam(name) {
    var  reg = new RegExp('(^|&)'+name+'=([^&]*)(&|$)', 'i');
    var r = window.location.search.substr(1).match(reg);
    if (r != null) {
        return decodeURIComponent(r[2]);
    }
    return null;
}

      /**
       * 概率计算
       * @param  {number} int 输入百分之几  （不需要%符号）
       * @returns {boolean}
       */
      function probabilityHappend(int) {
        if (typeof int !== 'number') {
          alert('请输入数字');
          return;
        }
        if (parseInt(int, 10) !== Number(int)) {
          alert('请输入整数');
          return;
        }
        if (int < 0) {
          alert('请输入正整数');
          return;
        }
        if (int === 0) {
          return false;
        }
        if (Math.round(Math.random() * 100) <= int) {
          return true;
        }
        return false;
      }

      // var isHappend = 0;
      // var notHappend = 0;
      //  var times=1000;
      //  var probabilityInt=30;
      //
      //  function probabilityHappend0(int) {
      //   if (typeof int !== 'number') {
      //     alert('请输入数字');
      //     return;
      //   }
      //   if (parseInt(int, 10) !== Number(int)) {
      //     alert('请输入整数');
      //     return;
      //   }
      //   if (int < 0) {
      //     alert('请输入正整数');
      //     return;
      //   }
      //   if (int === 0) {
      //       notHappend=notHappend+1
      //     return false;
      //   }
      //   if (Math.round(Math.random() * 100) <= int) {
      //       isHappend=isHappend+1
      //     return true;
      //   }
      //   notHappend=notHappend+1
      //   return false;
      // }
      //
      // for (let i=0;i<times;i++){
      //      probabilityHappend0(probabilityInt)
      // }
      // console.log('相加的',(isHappend/(isHappend+notHappend))*100)
      // console.log('times的',(isHappend/times)*100)


    var curURLToArray=window.location.host.split('.');
    // 注意：当前页面的网址必须类似的xxx.xxxxxx.xxxx，否则将出现错误！！！
    var domain=curURLToArray[1]+'.'+curURLToArray[2];
    var gameIdAndPathname=(window.location.pathname).replace('/','');
    var gameId=Number(gameIdAndPathname.split('.')[0]);
 //    // 赣州 余干和广昌 AppStore下载概率切到50%
 //    if(Number(gameId)===17 || Number(gameId)===28 || Number(gameId)===58){
 //    if(is_ios() && !getSearchParam('still')){      
 //         // 万年 app store 50%的概率走苹果官网
 //          if(probabilityHappend(50)){
 //            window.location.href='https://apps.apple.com/cn/app/%E4%B8%AD%E8%87%B3%E9%BA%BB%E5%B0%86/id1521388737'
 //           }else{
 //             window.location.href='https://h5.'+domain+'/ios/getChannel.html?game_id='+gameId
 //            }
 //     }
 //   }else{
 //        if(is_ios() && (!getSearchParam('still'))){
 //      window.location.href='https://h5.'+domain+'/ios/getChannel.html?game_id='+gameId
 //
 //     }
 //
 // }
 
 function isUseEnterprise(gameId) {
     var useEnterpriseGameList=[77,72,75,25,27,18,39,22,47,48,30,36,38,80,55,16,51,19,17,13,44,57,66,68,90,91,92];
     return useEnterpriseGameList.find(function(item) {
       return gameId ===item
     })
 }
 // 1、没有still 的 有企业签 2、没有still 没有企业签 3、有still 的 没有企业签 4、有still 的 有企业签
 
         if(is_ios()){
             if(!getSearchParam('still') || (getSearchParam('still') && !isUseEnterprise(gameId))){
                     window.location.href='https://h5.'+domain+'/ios/getChannel.html?game_id='+gameId
             }
     }






		// function reloadInWeixin() {
		//   if (is_weixin()) {
		//     if(getSearchParam('getTimeForReload')){
		//         return '';
		//     }else {
		//      var joiner=window.location.href.indexOf('?')===-1?'?':'&';
         //    window.location.href = window.location.href+joiner+'getTimeForReload='+ new Date().getTime();
		//     }
		//  }
		//
		// }
		// reloadInWeixin();


		`
	// if isAndroidRedirect && !isIosRedirect {
	// 	temp = temp + `$(document).ready(function(){
	// 		if (is_weixin() && is_android()) {
	// 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
	// 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
	// 				}
	// 				$('#guide').show();
	// 				$('#mask').show();}
	// 		});`
	// } else if !isAndroidRedirect && isIosRedirect {
	// 	temp = temp + `$(document).ready(
	// 		function(){
	// 			if (is_weixin() && is_ios()) {
	// 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
	// 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
	// 				}
	// 				$('#guide').show();
	// 				$('#mask').show();}
	// 		});`
	// } else {
	// 	temp = temp + `$(document).ready(
	// 		function(){
	// 			if (is_weixin() && (is_ios() || is_android())) {
	// 				if( /iPhone|iPad/i.test(navigator.userAgent)) {
	// 					document.getElementById("guide").src="https://static.xq5.com/3w4417com/img/guide_android.png";
	// 				}
	// 				$('#guide').show();
	// 				$('#mask').show();}
	// 		});`
	// }
	temp = temp + `
		</script>
			</head>
			<body>
				<div id="wrap">

						<div class="bg1">
						<img id='mask' style="float:right;position:fixed;top:0;display:none;height:100vh" src="` + staticUrl + `/3w4417com/img/black.png" >
						<img id='guide-android' style="width:78%;float:right;position:absolute;top:0;display:none;left:20%" src="` + staticUrl + `/3w4417com/img/guide_android.png" >
						<img id='guide-ios' style="width:78%;float:right;position:absolute;top:0;display:none;left:20%" src="` + htmlURL + `js/notchange/guide_ios.png" >
						</div>
					<div class="bg1">
					`
	if isAdShow {
		temp = temp + `<a class="getpack" id="ad" href="` + adURL + `">
						<img style="`
		if adImg == "" {
			temp = temp + "display:none"
		}
		temp = temp +
			`" src="` + adImg + `">
						</a>`
	}

	temp = temp +
		`
					<a class="getpack" id="load" onclick="alert_msg()" href="` + androidURL + `">
					<img src="` + headImg + `">
					</a>
					<a class="getpack">
					<div class="bg1">
						<img src="` + bodyImg + `"></div>
						<img id='xinren' style="float:right;position:absolute;bottom:0;display:none;" src="` + htmlURL + `xin.gif" >
					<iframe id="ifrm" height="0" width="0"></iframe>
					<form id="form1" action="insertDate.action" method="post">
						<input id="dUrl" name="dUrl" value="" type="hidden">
						<input id="appid" name="appid" value="100286" type="hidden">
						<input id="packetid" name="packetid" value="11058" type="hidden">
						<input id="from" name="from" value="thnnly" type="hidden">
					</form>
					</a>
					</div>
				

				 </div>


		<script>
		var _hmt = _hmt || [];
		(function() {
			var hm = document.createElement("script");
			hm.src = "https://hm.baidu.com/hm.js?29771296674c12a7bba9307344a0ed53";
			var s = document.getElementsByTagName("script")[0];
			s.parentNode.insertBefore(hm, s);
		})();
		// 此处 function is_ios 为重复代码 待删除 测试是否更新2020.4.17
		// function is_ios(){
		// 	var u = navigator.userAgent;
		// 	var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
		// 	if (isIOS == true) {
		// 		return true
		// 	}
		// 	return false
		// }

		 function is_weixin() {
			var ua = navigator.userAgent.toLowerCase();
			if (ua.match('micromessenger') == "micromessenger" ) {
				return true;
			} else {
				return false;
			}

		}

	     function is_android() {
			var u = navigator.userAgent;
			var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
			if (isAndroid == true) {
				return true
			}
			return false
		}
        //
		//   function is_ios(){
		// 	var u = navigator.userAgent;
		// 	var isIOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
		// 	if (isIOS == true) {
		// 		return true
		// 	}
		// 	return false
		// }
	   function is_ios(){
			// var u = navigator.userAgent;
			// var isIOS = !!u.match(/(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端
			// if (isIOS == true) {
			// 	return true
			// }
			// return false
		 var str = navigator.userAgent.toLowerCase();
		 var ver = str.match(/cpu iphone os (.*?) like mac os/) ? str.match(/cpu iphone os (.*?) like mac os/) : str.match(/cpu os (.*?) like mac os/) ? str.match(/cpu os (.*?) like mac os/) : str.match(/intel mac os x (.*?)/);
		 return ver
		}

	      function alert_msg(){
			if (!is_weixin()) {
				if (is_ios()){
            		// setTimeout("alert('正在安装，请返回桌面查看!')",2000);
            		// alert('执行了11111')

            		// document.getElementById("xinren").style.display="block";
            		console.log('ios系统')

            	}else if(!is_android()){
                       alert('请在浏览器设置中将【浏览器标识】或【浏览器UA标识】设置为您当前手机系统对应的选项或默认选项')
                }
			}
		}



		  // 在微信内置浏览器中  引导用户采用在微信外打开的方式
    $(function() {
        if (is_weixin()) {
            $('#mask').show();
            $('#mask').css({
                'z-index':8
            })
            if(is_android()){
                $('#guide-android').css({
                    'z-index':9
                });
                $('#guide-android').show();
            }else {
                $('#guide-ios').css({
                    'z-index':9
                });
                $('#guide-ios').show()

            }

        }

        		 // 在ios 微信内置浏览器中 展示 信任设置
		// if(is_ios() && is_weixin()){
		// 		document.getElementById("xinren").style.display="block";
		// }
    })
    // // 信任按钮 事件
	//     window.onload=function(){
    //               document.getElementById("xinren").onclick=function(){
    //               if((is_weixin())){
    //               alert('请点击右上角... 选择【在Safari中打开】后在Safari中点击该按钮')
    //               }else{
    //               var url="https://www.4417.com/js/notchange/guixiDev.mobileprovision";
    //                            window.open(url);
    //                            return true;
    //               }
    //            }
    //         }


            // 赋值app下载地址
		var u = navigator.userAgent;
		var isAndroid = u.indexOf('Android') > -1 || u.indexOf('Adr') > -1; //android终端
		// var isiOS = !!u.match(/\(i[^;]+;( U;)? CPU.+Mac OS X/); //ios终端


		if(!is_weixin()){
        if (is_ios()) {
        var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;'];
        document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];

		}else {
		 var tempArrAndroid = "` + androidURL + `" ? "` + androidURL + `".split('|') : ['javascript:;'];
        document.getElementById("load").href=tempArrAndroid[Math.floor(Math.random()*tempArrAndroid.length)];

		}
		}else{
		if(is_ios()){
			var tempArrIos = "` + iosURL + `" ? "` + iosURL + `".split('|') : ['javascript:;']
        			document.getElementById("load").href=tempArrIos[Math.floor(Math.random()*tempArrIos.length)];
		}else{
			document.getElementById("load").href="javascript:void(0);"
		}

		}


	  if(!is_weixin()) {
			setTimeout(function () {
				document.querySelector('#load').click();//直接触发下载链接
			},1000)
	}

		</script>
		</body>
		<footer class="footera">
			<div style="width:100%;height:30px">
				<ul class="footerUl">
					<li>
						<a href="index.html">首页</a>
					</li>
					<li>
						<a href="custody.html">家长监护</a>
					</li>
					<li>
						<a href="prevent.html">防沉迷</a>
					</li>
					<li>
						<a href="login.html">登录</a>
					</li>
					<li>
						<a href="register.html">注册</a>
					</li>
				</ul>
			</div>
			<p class="beian">
				<span>` + licenceNum1 + `</span>`
	if flag {
		temp +=
			`<span>
					<a href="img/增值电信业务经营许可证.jpg">赣B2—20150021</a>
				</span>
				<br>
				<span>
					<a href="img/网络文化经营许可证.jpg">赣网文〔2015〕1354-007号</a>
				</span>`
	}
	temp +=
		`<span>` + licenceNum2 + `</span>
			</p>
			<p class="copyright">` + pageIndex + `</p>
			    <div class="chuangyu-icon"><a target="cyxyv" href="https://v.yunaq.com/certificate?domain=www.4417.com&from=label&code=90030"> <img src="https://aqyzmedia.yunaq.com/labels/label_sm_90030.png"></a></div>
		</footer>
		<script >
		// 不是4417 隐藏创宇
		if(domain !=='4417.com'){
		     $('.chuangyu-icon').addClass('hide')
		}

		if(!getSearchParam('lcreload')){
	   // 页面可见时进行强制刷新

	    // 设置隐藏属性和改变可见属性的事件的名称
    var hidden, visibilityChange;
    if (typeof document.hidden !== "undefined") { // Opera 12.10 and Firefox 18 and later support
        hidden = "hidden";
        visibilityChange = "visibilitychange";
    } else if (typeof document.msHidden !== "undefined") {
        hidden = "msHidden";
        visibilityChange = "msvisibilitychange";
    } else if (typeof document.webkitHidden !== "undefined") {
        hidden = "webkitHidden";
        visibilityChange = "webkitvisibilitychange";
    }


    // 如果浏览器不支持addEventListener 或 Page Visibility API 给出警告
    if (typeof document.addEventListener === "undefined" || typeof document[hidden] === "undefined") {
         alert('浏览器版本过低,建议您升级浏览器或更换浏览器')
    } else {
        // 处理页面可见属性的改变
        function handleVisibilityChange() {
            if (document[hidden]) {
               console.log('隐藏了')
            } else {
                // alert('即将刷新')
                setTimeout(function() {
                  window.location.reload(true)
                },5000)
            }
        }
        document.addEventListener(visibilityChange, handleVisibilityChange, false);
    }

		}



        </script>
		</html>`

	var d1 = []byte(temp)
	if err := ioutil.WriteFile(html, d1, 0666); err != nil {
		log.Errorf("CreateHTML: err: %v\n", err.Error())
		return false
	}
	return true
}

type OnDownloadPageOpRequest struct {
	GameID int  `json:"game_id" binding:"required"`
	Close  bool `json:"close"`
}

func GetAllPage(db *gorm.DB) (*[]HtmlPageTable, error) {
	var pages []HtmlPageTable
	err := db.Table("html_page").Where("game_id != 0").Find(&pages).Error
	return &pages, err
}

type GameStatus struct {
	GameID int32 `json:"game_id" gorm:"game_id"`
	Close  bool  `json:"close" gorm:"close"`
}

func GetGameIDListCloseStatus(db *gorm.DB, gameIDList []int32) (statusList []GameStatus, err error) {
	err = db.Table("html_page").Select("game_id,close").
		Where("game_id in (?)", gameIDList).Find(&statusList).Error
	return
}

func GetHtmlPageByGameID(qipaidb *gorm.DB, gameID int) (*HtmlPageTable, error) {
	htmlPageConfig := &HtmlPageTable{}
	if err := qipaidb.Table("html_page").Where("game_id = ?", gameID).First(htmlPageConfig).Error; err != nil {
		log.Errorf("OnDownloadPageSelectRequest: err: %v\n", err.Error())
		return nil, err
	}
	return htmlPageConfig, nil

}
