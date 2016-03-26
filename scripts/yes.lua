yes = {
    father = "modify",
    name = "yes",
    typ = "POST",
    pass = "true",
    needArgs = "name address",
    handler = function(req)
        return ModRow(modify,req)
    end
}
