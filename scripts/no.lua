no = {
    father = "modify",
    name = "no",
    typ = "POST",
    needArgs = "name address",
    handler = function(req)
        return ModRow(modify,req)
    end
}
