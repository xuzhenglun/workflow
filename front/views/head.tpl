{{define "head"}}

   <ul class="nav navbar-nav navbar-right">
        {{if .}}
          <li class="dropdown">
    		<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button">{{.Name}}<span class="caret"></span></a>
    		<ul class="dropdown-menu">
        	<li role="separator" class="divider"></li>
        	<li class="dropdown-header">SignedTime</li>
                <li><a>{{.SignTimeString}}</a></li>
        	<li class="dropdown-header">ExpireTime</li>
                <li><a>{{.ExpireTimeString}}</a></li>
                <li><a href="/logout">Log out</a></li>
        {{else}}
                <li><a href="/login">Log In</a></li>
                <li><a href="/Sign">Sign Up</a></li>
        {{end}}
    </ul>

</div>
<!--/.nav-collapse -->
</div>
</nav>

{{end}}
