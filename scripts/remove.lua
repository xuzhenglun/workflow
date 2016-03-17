remove = {
    name = "remove",
    typ = "GET",
    helper = "this is remove",
    handler = function(n)
        log(n)
        DelRow(n)
        return "Done"
    end
}

