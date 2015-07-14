var util = {
    parseQueryString: function(route, queryString) {
        var params = {}, queries, temp, i, l;
        // Split into key/value pairs
        queries = queryString.replace("#!/" + route + "/", "");
        queries = decodeURIComponent(queries).split("&");
        // Convert the array of strings into an object
        for (i = 0, l = queries.length; i < l; i++) {
            temp = queries[i].split('=');
            params[temp[0]] = decodeURIComponent(temp[1]);
        }
        return params;
    }
}
