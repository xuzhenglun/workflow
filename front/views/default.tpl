{{template "head" .User}}


<div class="container">
  <div class="row">
    <div class="col-sm-3 col-md-2 sidebar">
      <ul class="nav nav-sidebar">
        {{$name := .Name}} {{range .Activities}}
        <li {{if eq . $name}}class="active" {{end}}><a href="/{{.}}">{{.}} <span class="sr-only">(current)</span></a></li>
        {{end}}
      </ul>
    </div>

    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <div class="start-template">
        <div class="activity-name">
          <h1>{{.Name}}</h1>
        </div>


        <div class="bs-example">
          <table class="table table-bordered">
            <tbody>
              {{range $index, $elem := .Info}}
              <tr>
                <th scope="row" style="width: 15%">{{$index}}</th>
                <td>{{$elem}}</td>
              </tr>
              {{end}}
            </tbody>
          </table>
        </div>

        <form class="form-group" action="" method="post" id="FORM" target="id_iframe">
          <div class="form-group">
            <input type="hidden" class="form-control" id="pass" name="pass" value="false" />
            {{range .Args}}
            <label class="pull-left">{{.}}</label>
            <input type="text" class="form-control" id={{.}}>
			{{end}}
          </div>
        </form>
          <div class="button-group">
            <p style="line-height: 70px; text-align: center;">
              {{if .IsPass}}
              <button class="btn btn-danger btn-lg" id="back" onclick="post(`false`);">拒绝</button>
              <button class="btn btn-success btn-lg" id="back" onclick="post(`true`);">允许</button>
              {{else}}
              <button class="btn btn-primary btn-lg" id="back" onclick="post(`false`);">确定提交</button>
              {{end}}
            </p>
          </div>
          <script type="text/javascript">
            function post(arg) {
			  document.getElementById("FORM").pass.value=arg
			  var form = $(FORM).serializeObject()
				  $.post("",form,function(data){
					var status = JSON.parse(data)
					if (status.Code == 200){
 						alert("Success");
						location.href = document.referrer
					}else{
						alert(status.Msg);
					}
				  })
            }
          </script>

        <iframe id="id_iframe" name="id_iframe" style="display:none"></iframe>
      </div>
    </div>
  </div>
</div>
