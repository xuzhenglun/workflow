modify = {
    name = "modify",
    typ = "POST",
    needArgs = "name address",
    handler = function(req)
        log("this is modifier")
        return ModRow(modify,req)
    end
}
