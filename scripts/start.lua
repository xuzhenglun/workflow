start = {
    name = "start",
    typ = "POST",
    helper = "name address",
    handler = function(n)
        log("this is starter")
        return AddRow(n)
    end
}

