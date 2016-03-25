remove = {
    name = "remove",
    typ = "GET",
    father = "start",
    needArgs = "this is remove",
    handler = function(req)
        return DelRow(remove,req)
    end
}

