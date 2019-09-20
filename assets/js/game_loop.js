"use strict";
console.log("game_loop start " + Date.now());

function onmessage(e) {
    console.log(e);
    postMessage("recieved");
}