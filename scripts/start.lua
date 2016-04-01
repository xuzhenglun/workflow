start = {
    name = "start",
    typ = "POST",
    groups = "start",
    needArgs = "name address idcard",
    handler = function(req)
        log("this is starter")
        return AddRow(start,req)
    end
}

