# OAuth2-app
an OAuth2-app including server and client.

OAuth2.0授权服务实践

1.OAuth简单介绍：
OAuth（开放授权）是一个开放标准，允许用户让第三方应用访问该用户在某一网站上存储的私密的资源（如照片，视频，联系人列表），而无需将用户名和密码提供给第三方应用。
    OAuth  允许用户提供一个令牌，而不是用户名和密码来访问他们存放在特定服务提供者的数据。每一个令牌授权一个特定的网站（例如，视频编辑网站)在特定的时段（例如，接下来的2小时内）内访问特定的资源（例如仅仅是某一相册中的视频）。
可以联想一下微信公众平台开发，在微信公众平台中当我们访问某个页面，页面可能弹出一个提示框，提示某第三方应用需要获取我们的个人信息问是否允许，点确认其实就是授权第三方应用获取我们在微信公众平台的个人信息。这里微信网页授权就是使用的OAuth2.0。
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/1.png) 

2.认证授权过程
在认证和授权的过程中涉及的三方包括：
1、服务提供方，用户使用服务提供方来存储受保护的资源，如照片，视频，联系人列表。
2、用户，存放在服务提供方的受保护的资源的拥有者。
3、客户端，要访问服务提供方资源的第三方应用，通常是网站，如提供照片打印服务的网站。在认证过程之前，客户端要向服务提供者申请客户端标识。

使用OAuth进行认证和授权的主要过程如下图：
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/2.png) 

上面的client就是第三方网站，Resource Server资源服务器，Authorization Server授权服务器。其中资源服务器和授权服务器可以是同一个。
Access Token每个Client有一个并且有时效性。

3.几种授权模式
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/3.png) 

1、授权码模式，是功能最完整、流程最严密的授权模式。
2、简化模式，不通过第三方应用的服务器，直接在浏览器中向认证服务器申请令牌，跳过了“授权码”这个步骤。
3、密码模式，需要用户向客户端提供自己的用户名和密码。
4、客户端模式，指客户端以自己的名义，而不是以用户的名义，向“服务提供商”进行认证。

4.OAuth2-demo实践
1.项目工程结构如下：

项目包括授权服务器实现和客户端实现，分别运行server和client，服务端和客户端分别监听在9000和9001端口上。然后在网页浏览器进行操作即可。

2.网页上操作流程：
A. 首先输入localhost:9001/login地址，请求客户端的服务：
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/5.png) 

B. 请求上述服务之后，会重定向到服务端去认证，认证会经历 登录 - 授权 的流程。服务端登录，页面如下：
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/6.png) 

C. 服务端授权，页面如下：
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/7.png) 

D.点击授权进行授权之后，服务端会生成授权码code，并重定向到客户端的callback页面。callback过程中客户端会拿着这个code去发送POST请求给服务端，申请access_token等信息。
![介绍](https://github.com/guobingithub/OAuth2-app/blob/master/image/8.png) 

客户端拿到access_token信息之后，可以凭借它在特定的时间内，去访问资源服务器的特定资源。

Over!
