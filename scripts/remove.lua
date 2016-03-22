remove = {
    name = "remove",
    typ = "GET",
    needArgs = "this is remove",
    handler = function(req)
        return DelRow(remove,req)
    end
}

