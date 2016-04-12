{{template "head"}}

<div class="container">

  <form method="POST" action="/login" class="form-signin">
    <h1 class="form-signin-heading" style="padding:20px">Please sign in <small>Input the License</small></h1>
    <label for="inputLicense" class="sr-only">License</label>
    <textarea  class="form-control" rows="10" id="auth" name="auth" ></textarea>
    <div class="checkbox">
      <label>
        <input type="checkbox" value="on" id="long" name="long"> Remember me</input>
      </label>
    </div>
    <button class="btn btn-lg btn-primary btn-block" type="submit">Sign in</button>
  </form>

</div>
<!-- /container -->
