nomodify = {
    father = "start",
    name = "modify",
    typ = "POST",
    pass = "true",
    needArgs = "name address",
    handler = function(req)
        log("this is modifier")
        log("this is hotplug")
        return ModRow(modify,req)
    end
}
