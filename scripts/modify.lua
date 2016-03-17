modify = {
    name = "modify",
    typ = "POST",
    helper = "name address",
    handler = function(n)
        log("this is modifier")
        ModRow(n,str)
        return "Done"
    end
}
