modify = {
    father = "start",
    name = "modify",
    typ = "POST",
    needArgs = "pass",
    handler = function(req)
        return ModRow(modify,req)
    end
}
