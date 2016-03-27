check = {
    father = "start",
    name = "check",
    typ = "GET",
    needArgs = "address",
    handler = function(req)
        log("this is checker")
        a=FindRow(check,req)
        return a
    end
}
