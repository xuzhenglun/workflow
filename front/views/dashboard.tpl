{{template "head" .User}}


<div class="container-fluid">
  <div class="row">
    <div class="col-sm-3 col-md-2 sidebar">
      <ul class="nav nav-sidebar">
        {{$name := .Name}}
        {{range .Activities}}
        <li {{if eq . $name}}class="active" {{end}}><a href="/{{.}}">{{.}} <span class="sr-only">(current)</span></a></li>
        {{end}}
      </ul>
    </div>
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <h1 class="sub-header">Section title</h1>
      <div class="table-responsive">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>#事件ID</th>
              <th>动作</tr>
          </thead>
          <tbody>
            {{range .ids}}
            <tr>
              <td>{{.}}</td>
              <td><a href="/{{$name}}/{{.}}">处理</a></td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </div>
    </div>
  </div>
</div>
