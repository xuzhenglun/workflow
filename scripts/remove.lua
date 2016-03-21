remove = {
    name = "remove",
    typ = "GET",
    helper = "this is remove",
    handler = function(n)
        log(n)
        return DelRow(n)
    end
}

