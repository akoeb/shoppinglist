export const httpOptions = (m, obj) => {
    // this object will be returned
    var options = {
        method: m !== undefined ? m : "GET",
        mode: "cors",
        cache: "no-cache",
        credentials: "same-origin",
        headers: {
            "Content-Type": "application/json"
        },
        redirect: "follow",
    }
    if (obj !== undefined) {
        options['body'] = JSON.stringify(obj)
    }
    return options
}

export const backend = (path) => {
    // TODO: pass in environment, set it globally
    // TODO: move backend and httpOptions functions to some global place
    let env = "";
    let url = "http://localhost:8080/";
    if (env === "production") {
        url = "/";
    }

    return url + path;
}
