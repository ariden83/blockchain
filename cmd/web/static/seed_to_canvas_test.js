const cutSeed = str => {
    const words = str.split(' ');
    for (const w of words) {
        let c = document.createElement("CANVAS");
        let ctx = c.getContext("2d");
        ctx.font = "30px Calibri";
        ctx.setTransform((Math.random() / 10) + 0.9,    //scalex
            0.1 - (Math.random() / 5),      //skewx
            0.1 - (Math.random() / 5),      //skewy
            (Math.random() / 10) + 0.9,     //scaley
            (Math.random() * 3) + 3,      //transx
            0);                           //transy
        ctx.fillText(w, 10, 50);

        let div = document.createElement("DIV");
        div.className = "col-xs-3"
        div.appendChild(c);
        document.getElementById("toto").appendChild(div);
    }
}