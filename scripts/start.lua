start = {
    name = "start",
    typ = "POST",
    needArgs = "name address",
    handler = function(req)
        log("this is starter")
        return AddRow(start,req)
    end
}

