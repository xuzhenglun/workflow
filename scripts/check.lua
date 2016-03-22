check = {
    name = "check",
    typ = "GET",
    needArgs = "",
    handler = function(req)
        log("this is checker")
        a=FindRow(check,req)
        return a
    end
}
