check = {
    name = "check",
    typ = "GET",
    helper = "",
    handler = function(n)
        log("this is checker")
        a=FindRow(n)
        return a
    end
}
