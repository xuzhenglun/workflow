modify = {
    name = "modify",
    typ = "POST",
    helper = "name address",
    handler = function(n)
        log("this is modifier")
        return ModRow(n,str)
    end
}
