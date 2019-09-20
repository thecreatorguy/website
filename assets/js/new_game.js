"use strict";

class Game {
    constructor(width, height) {
        this.canvas = document.createElement("canvas");
        this.canvas.width = width;
        this.canvas.height = height;
        document.getElementById("game_container").append(this.canvas);
        this.context = this.canvas.getContext("2d");



        this.gameLoop = new Worker("assets/js/game_loop.js");
        console.log(this.gameLoop.onmessage);
        this.gameLoop.onmessage = (e) => {
            console.log(e.data);
        }

        this.gameLoop.onerror = (e) => console.log(e);
    }

    start() {
        console.log("message sent");
        this.gameLoop.postMessage("hi");
    }
}

let game;
function startGame() {
    game = new Game(400, 400);
    game.start();
}
document.addEventListener("DOMContentLoaded", startGame);

