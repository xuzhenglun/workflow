<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
  <meta name="description" content="">
  <meta name="author" content="">
  <link rel="icon" href="../../favicon.ico">

  <title>Index -Workflow</title>

  <!-- Bootstrap core CSS -->
  <link href="/static/css/bootstrap.min.css" rel="stylesheet">

  <!-- Custom styles for this template -->
  <link href="/static/css/cover.css" rel="stylesheet">

  <!-- Just for debugging purposes. Don't actually copy these 2 lines! -->
  <!--[if lt IE 9]><script src="../../assets/js/ie8-responsive-file-warning.js"></script><![endif]-->
  <script src="/static/js/ie-emulation-modes-warning.js"></script>

  <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
  <!--[if lt IE 9]>
      <script src="//cdn.bootcss.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="//cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>

<body>

  <div class="site-wrapper">

    <div class="site-wrapper-inner">

      <div class="cover-container">

        <div class="masthead clearfix">
          <div class="inner">
            <h3 class="masthead-brand">Workflow</h3>
            <nav>
              <ul class="nav masthead-nav">
                <li class="active"><a href="#">Home</a></li>
                {{if .User}}
                <li class="dropdown">
                  <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button">{{.User.Name}}<span class="caret"></span></a>
                  <ul class="dropdown-menu">
                    <li role="separator" class="divider"></li>
                    <li class="dropdown-header">SignedTime</li>
                    <li><a>{{.User.SignTimeString}}</a></li>
                    <li class="dropdown-header">ExpireTime</li>
                    <li><a>{{.User.ExpireTimeString}}</a></li>
                    <li><a href="/logout">Log out</a></li>
                    {{else}}
                    <li><a href="/login">Log In</a></li>
                    <li><a href="/Sign">Sign Up</a></li>
                    {{end}}

                  </ul>
            </nav>
          </div>
        </div>

        <div class="inner cover">
          <h1 class="cover-heading">Default Index Page</h1>
          <p class="lead">Welcome to Workflow
            <br />This is a default index page</p>
          <p class="lead">
            <a href="/start" class="btn btn-lg btn-default">Let's Begin</a>
          </p>
        </div>

        <div class="mastfoot">
          <div class="inner">
            <p>Page design</a> by <a href="https://twitter.com/mdo">@mdo</a>.</p>
          </div>
        </div>

      </div>

    </div>

  </div>

  <!-- Bootstrap core JavaScript
    ================================================== -->
  <!-- Placed at the end of the document so the pages load faster -->
  <script src="/static/js/jquery.min.js"></script>
  <script src="/static/js/bootstrap.min.js"></script>
  <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
  <script src="/static/js/ie10-viewport-bug-workaround.js"></script>
</body>

</html>
