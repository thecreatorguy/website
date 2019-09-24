"use strict";

/** Map of all colors in game, assigned here to make it easier to modify */
let colorMap = new Map();
colorMap.set("border",          "rgb(84, 84, 69)");
colorMap.set("floor",           "rgb(119, 119, 98)");
colorMap.set("wall",            "rgb(255, 255, 255)");
colorMap.set("lava",            "rgb(237, 60, 21)");
colorMap.set("goal",            "rgb(80, 244, 66)");
colorMap.set("switch1a",        "rgb(226, 180, 93)");
colorMap.set("switch1b",        "rgb(244, 213, 154)");
colorMap.set("disappearing1a",  "rgb(255, 182, 0)");
colorMap.set("disappearing1d",  "rgb(142, 101, 0)");
colorMap.set("switch2a",        "rgb(56, 202, 216)");
colorMap.set("switch2b",        "rgb(146, 224, 232)");
colorMap.set("disappearing2a",  "rgb(0, 255, 242)");
colorMap.set("disappearing2d",  "rgb(0, 96, 91)");
colorMap.set("switch3a",        "rgb(92, 119, 219)");
colorMap.set("switch3b",        "rgb(153, 144, 237)");
colorMap.set("disappearing3a",  "rgb(0, 55, 255)");
colorMap.set("disappearing3d",  "rgb(0, 24, 112)");
colorMap.set("tp1a",            "rgb(255, 137, 215)");
colorMap.set("tp1b",            "rgb(237, 68, 180)");
colorMap.set("tp2a",            "rgb(185, 56, 255)");
colorMap.set("tp2b",            "rgb(124, 1, 191)");
colorMap.set("hero",            "rgb(245, 249, 0)");   
colorMap.set("pause",           "rgba(0, 0, 0, 0.4");

/** Number of levels that exist */
const maxLevel = 12;

/** 
 * Class to deconstruct loaded txt file from server into useful level information
 */
class Level {
    /**
     * Creates the level from server information. Important variables include:
     * 
     * title - Title of level
     * grid - 2D array of level information, including floors, walls, lava, goal, etc.
     * spawn - Coordinate pair of where hero is to be spawned
     * toggleTilesState - 3 element array of booleans depicting starting state of each switch tile
     * tp[1, 2]Locations - 2 element arrays of coordinate pairs for the two teleporter types
     * 
     * @param {String} text text information from server
     */
    constructor(text) {
        let lines = text.split("\n");
        let lineIter = lines[Symbol.iterator]();
        let line = "";

        this.grid = [];
        this.tp1Locations = [];
        this.tp2Locations = [];
        for (let i = 0; i < 25; i++) {
            this.grid.push([]);
            line = lineIter.next().value;
            for (let j = 0; j < 25; j++) {
                let tileId = Number(line.slice(j*2, j*2+2));
                this.grid[i].push(tileId);
                if (tileId == 10) {
                    this.tp1Locations.push({x: j, y: i})
                }
                if (tileId == 11) {
                    this.tp2Locations.push({x: j, y: i})
                }
            }
        }
        line = lineIter.next().value;
        this.spawn = {x: Number(line.slice(0, 2)), 
                      y: Number(line.slice(2, 4))};

        this.toggleTilesState = [];
        for (let i = 0; i < 3; i++) {
            line = lineIter.next().value;
            this.toggleTilesState.push(Number(line));
        }
        
        this.title = lineIter.next().value;
    }
}

/**
 * Generic class that structures how a game is to be made and run.
 * Requires a div id "game_container" to place the canvas at the end of.
 */
class GameArea {
    /**
     * Creates the game area, which includes a canvas.
     * 
     * @param {Number} width width of canvas in px
     * @param {Number} height height of canvas in px
     * @param {Number} fps times per second to run the main loop. includes game logic as well as frames
     */
    constructor(width, height, fps, gameUpdateFunc, renderFrameFunc) {
        this.canvas = document.createElement("canvas");
        this.canvas.width = width;
        this.canvas.height = height;
        document.getElementById("game_container").append(this.canvas);
        this.context = this.canvas.getContext("2d");
        this.fps = fps;
    }

    /**
     * Set the functions to update the game
     * 
     * @param {Function} gameUpdateFunc function to run at the beginning of each iteration of the main loop
     * @param {Function} renderFrameFunc function to run after gameUpdateFunc in each iteration of the main loop
     */
    setUpdateFunctions(gameUpdateFunc, renderFrameFunc) {
        this._gameUpdate = gameUpdateFunc;
        this._renderFrame = renderFrameFunc;
    }

    /**
     * Starts the game loop if not already started, with an interval specified by the frames per second
     */
    start() {
        if (!this.interval) {
            this.interval = setInterval(this._update.bind(this), 1000 / this.fps);
        }
    }

    /**
     * Regularly called function to update game components and then to render a new frame 
     */
    _update() {
        // Update Game
        this._gameUpdate();

        // Render Frame
        this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.context.fillStyle = "black";
        this.context.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this._renderFrame();
    } 
}

// Pixels to move the player by in each update of the frame
const playerSpeed = 5;

/** 
 * Specific class to run the Slider Game. The game is a puzzle game
 * where the player's goal is to slide their square to the green
 * goal tile while avoiding obstacles and using other game
 * mechanics to their advantage. The player should beware, however:
 * Once started, the square will not stop moving until it hits a
 * solid block! Early levels are designed to teach concepts, while
 * later levels are designed to be more of a challenge. Average
 * playtime for a new player is between 20 and 30 minutes.
 */
class SliderGame extends GameArea {
    /**
     * Creates the Slider Game, including the canvas.
     * Requires div id "level_number" to write level number to,
     * div id "level_title" to write title to, and
     * div id "timer" to write time data to.
     */
    constructor() {
        // Slider Game is run at 50 fps, which is intense for some computers, but at fewer frames
        // runs fairly choppily. I will leave this for now, and maybe revise later
        super(500, 500, 50);
        super.setUpdateFunctions(this._sliderGameUpdate.bind(this), this._sliderRenderFrame.bind(this));

        // All local variables
        this.level = 1;
        this.levelGrid = [];
        this.player = {x: 0, y: 0, xVel: 0, yVel: 0, dir: "none", startedSliding: false};
        this.toggleTilesState = [];
        this.accumulatedTime = 0;
        this.lastDate = Date.now();
        this.startHours = new Date(0).getHours();
        this.paused = false;

        // References to document elements
        this.levelNumber = document.getElementById("level_number");
        this.levelTitle = document.getElementById("level_title");
        this.timer = document.getElementById("timer");

        // Add controls variables and event listeners
        this.unprocessedKeys = new Map();
        document.addEventListener("keydown", (event) => {
            this.unprocessedKeys.set(event.code, true);
        });
        document.addEventListener("keyup", (event) => {
            this.unprocessedKeys.set(event.code, false);
        });

        this.mouseCommand = "none";
        this.mouseDown = false;
        document.addEventListener("click", (event) => {
            let dx = event.clientX - (this.canvas.getBoundingClientRect().left + this.player.x);
            let dy = event.clientY - (this.canvas.getBoundingClientRect().top + this.player.y);
            if (Math.abs(dx) > Math.abs(dy)) {
                if (dx > 15 && dx < 200) {
                    this.mouseCommand = "right";
                }
                else if (dx < -15 && dx > -200) {
                    this.mouseCommand = "left";
                }
            }
            else {
                if (dy > 15 && dy < 200) {
                    this.mouseCommand = "down";
                }
                else if (dy < -15 && dy > -200) {
                    this.mouseCommand = "up";
                }
            }
        });

        // Initialize level variables
        this._setLevel();
    }

    /**
     * Updates the slider game components, including the paused state, 
     * time and player logic
     */
    _sliderGameUpdate() {
        // Pause or Unpause game
        if (this._keyPressed("Space")) {
            this.paused = !this.paused;
        }

        // Update time
        this._updateTime();
        
        // If the game is won, check if the player wants to restart game
        if (this.level > maxLevel) {
            if (this._keyPressed("KeyR")) {
                this.accumulatedTime = 0;
                this.lastDate = Date.now();
                this.startHours = new Date(0).getHours();
                this.level = 1;
                this._setLevel();
                return;
            }
        }
        else if (!this.paused) { // Otherwise, update the player logic
            if (this._updatePlayer()) {
                this._setLevel();
            } 
        }
    }

    /**
     * Update the time value in the "timer" div
     */
    _updateTime() {
        // Determine what to append to time text content, and update time if required
        let append = "";
        if (this.level > maxLevel) {
            append = " < WINNING TIME!"
        }
        else if (this.paused) {
            append = " | PAUSED";
        }
        else {
            this.accumulatedTime += Date.now() - this.lastDate;
        }
        this.lastDate = Date.now();

        // Update the time value in the div
        let nowDate = new Date(this.accumulatedTime);
        let twoDigitFormat = new Intl.NumberFormat('en-US', {minimumIntegerDigits: 2});
        if (nowDate.getHours() - this.startHours > 0) {
            this.timer.textContent = `Time: ${nowDate.getHours() - this.startHours}:`+
                                            `${twoDigitFormat.format(nowDate.getMinutes())}:` +
                                            `${twoDigitFormat.format(nowDate.getSeconds())}${append}`;
        }
        else {
            this.timer.textContent = `Time: ${twoDigitFormat.format(nowDate.getMinutes())}:` +
                                            `${twoDigitFormat.format(nowDate.getSeconds())}${append}`;
        }
    }

    /**
     * Updates the player location depending on the controls, velocity, and grid.
     * 
     * @returns {Boolean} true if the level needs to be set/reset, false otherwise
     */
    _updatePlayer() {
        // If 'R' is pressed, reset the level
        if(this._keyPressed("KeyR")) {
            return true;
        }

        // If the player is stationary and has a command waiting, start moving in that direction.
        // If there is more than one, it prioritizes left > down > right > up.
        if (this.player.dir == "none") {
            this.player.xVel = 0;
            this.player.yVel = 0;  
            if (this._hasCommand("left")) {
                this.player.dir = "left";
                this.player.startedSliding = true;
                this.player.xVel = -playerSpeed;
            }
            else if (this._hasCommand("down")) {
                this.player.dir = "down";
                this.player.startedSliding = true;
                this.player.yVel = playerSpeed;
            }
            else if (this._hasCommand("right")) {
                this.player.dir = "right";
                this.player.startedSliding = true;
                this.player.xVel = playerSpeed;
            }
            else if (this._hasCommand("up")) {
                this.player.dir = "up";
                this.player.startedSliding = true;
                this.player.yVel = -playerSpeed;
            }
        }

        // This block is triggered if the player has perfectly passed over a grid square.
        // Any functionality that is dependent on landing on a square is triggered.
        // This is done before stopping because velocity is conserved in the event that
        // this functionality causes whatever would stop velocity to no longer stop it,
        // i.e. a button switches off a wall, or the player is teleported to a new location.
        if (!this.player.startedSliding && 
            (this.player.x / 20) % 1 == 0 && (this.player.y / 20) % 1 == 0 && 
            this.player.dir != "none") {
            switch (this.levelGrid[this.player.y / 20][this.player.x / 20]) {
                // 2 is lava- returning true resets the level
                case 2: 
                    return true;
                // 4 is orange buttons- toggle the tile state
                case 4: 
                    this.toggleTilesState[0] = !this.toggleTilesState[0];
                    break;
                // 6 is teal buttons- toggle the tile state
                case 6: 
                    this.toggleTilesState[1] = !this.toggleTilesState[1];
                    break;
                // 8 is blue buttons- toggle the tile state
                case 8: 
                    this.toggleTilesState[2] = !this.toggleTilesState[2];
                    break;
                // 10 is pink teleporters- move player to other teleporter location
                case 10: 
                    if (this.player.x / 20 == levels[this.level].tp1Locations[0].x &&
                        this.player.y / 20 == levels[this.level].tp1Locations[0].y) {
                        this.player.x = levels[this.level].tp1Locations[1].x * 20;
                        this.player.y = levels[this.level].tp1Locations[1].y * 20;
                    }
                    else {
                        this.player.x = levels[this.level].tp1Locations[0].x * 20;
                        this.player.y = levels[this.level].tp1Locations[0].y * 20;
                    }
                    break;
                // 11 is purple teleporters- move player to other teleporter location
                case 11: 
                    if (this.player.x / 20 == levels[this.level].tp2Locations[0].x &&
                        this.player.y / 20 == levels[this.level].tp2Locations[0].y) {
                        this.player.x = levels[this.level].tp2Locations[1].x * 20;
                        this.player.y = levels[this.level].tp2Locations[1].y * 20;
                    }
                    else {
                        this.player.x = levels[this.level].tp2Locations[0].x * 20;
                        this.player.y = levels[this.level].tp2Locations[0].y * 20;
                    }
                    break;

                default: break;
            }
        }
        else if (this.player.startedSliding) { 
            // This variable ensures that mechanics are not triggered when you
            this.player.startedSliding = false; // have just started moving
        }

        // Now we need to determine where the player would end up if allowed to continue
        // moving in the direction that it is. This allows us to stop the movement if
        // the location the hero is moving into is a solid block.
        let xIndex, yIndex;
        if (this.player.dir == "right") {
            xIndex = Math.floor((this.player.x + this.player.xVel + (20 - playerSpeed)) / 20);
        }
        else {
            xIndex = Math.floor((this.player.x + this.player.xVel) / 20);
        }
        if (this.player.dir == "down") {
            yIndex = Math.floor((this.player.y + this.player.yVel + (20 - playerSpeed)) / 20);
        }
        else {
            yIndex = Math.floor((this.player.y + this.player.yVel) / 20);
        }

        // If the player is moving and the location they are moving to 
        // is valid, stop; otherwise update location
        if (this.player.dir != "none") {
            if (xIndex >= 0 && xIndex < 25 && yIndex >= 0 && yIndex < 25) {
                if (this.levelGrid[yIndex][xIndex] == 1 ||
                        this.toggleTilesState[0] && this.levelGrid[yIndex][xIndex] == 5 ||
                        this.toggleTilesState[1] && this.levelGrid[yIndex][xIndex] == 7 ||
                        this.toggleTilesState[2] && this.levelGrid[yIndex][xIndex] == 9 ||
                        !this.toggleTilesState[0] && this.levelGrid[yIndex][xIndex] == 12 ||
                        !this.toggleTilesState[1] && this.levelGrid[yIndex][xIndex] == 13 ||
                        !this.toggleTilesState[2] && this.levelGrid[yIndex][xIndex] == 14 ) {
                    this.player.dir = "none";
                }
                else {
                    this.player.x += this.player.xVel;
                    this.player.y += this.player.yVel;
                }
            }
            else {
                this.player.dir = "none";
            }
        }
        
        // The final check is for the goal. This is done here because the goal can only be said
        // to be reached if the player is stopped on it.
        if (this.player.dir == "none" && this.levelGrid[this.player.y / 20][this.player.x / 20] == 3) {
            this.level++;
            return true;
        }

        // If the function made it this far, then the level does not need to be reset
        return false;
    }

    /**
     * Renders an entire frame of the slider game by rendering each background tile,
     * and then the player on top of it, and then the foreground greying out if
     * the game has been won or paused.
     */
    _sliderRenderFrame() {
        // Render the background
        let left = 0;
        let top = 0;
        for (let row = 0; row < 25; row++) {
            for (let col = 0; col < 25; col++) {
                this.context.fillStyle = colorMap.get("border");
                this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                switch (this.levelGrid[row][col]) {
                    case 0: 
                        this.context.fillStyle = colorMap.get("floor");
                        this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        break;
                    case 1: 
                        this.context.fillStyle = colorMap.get("wall");
                        this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        break;
                    case 2: 
                        this.context.fillStyle = colorMap.get("lava");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        break;
                    case 3: 
                        this.context.fillStyle = colorMap.get("goal");
                        this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        break;
                    case 4: 
                        this.context.fillStyle = colorMap.get("switch1a");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        this.context.fillStyle = colorMap.get("switch1b");
                        this.context.fillRect(left + col * 20 + 10, top + row * 20, 10, 10);
                        this.context.fillRect(left + col * 20, top + row * 20 + 10, 10, 10);
                        break;
                    case 5: 
                        if (this.toggleTilesState[0]) {
                            this.context.fillStyle = colorMap.get("disappearing1a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing1d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 6: 
                        this.context.fillStyle = colorMap.get("switch2a");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        this.context.fillStyle = colorMap.get("switch2b");
                        this.context.fillRect(left + col * 20 + 10, top + row * 20, 10, 10);
                        this.context.fillRect(left + col * 20, top + row * 20 + 10, 10, 10);
                        break;
                    case 7: 
                        if (this.toggleTilesState[1]) {
                            this.context.fillStyle = colorMap.get("disappearing2a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing2d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 8: 
                        this.context.fillStyle = colorMap.get("switch3a");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        this.context.fillStyle = colorMap.get("switch3b");
                        this.context.fillRect(left + col * 20 + 10, top + row * 20, 10, 10);
                        this.context.fillRect(left + col * 20, top + row * 20 + 10, 10, 10);
                        break;
                    case 9: 
                        if (this.toggleTilesState[2]) {
                            this.context.fillStyle = colorMap.get("disappearing3a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing3d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 10: 
                        this.context.fillStyle = colorMap.get("tp1a");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        this.context.fillStyle = colorMap.get("tp1b");
                        this.context.fillRect(left + col * 20 + 10, top + row * 20, 10, 10);
                        this.context.fillRect(left + col * 20, top + row * 20 + 10, 10, 10);
                        break;
                    case 11: 
                        this.context.fillStyle = colorMap.get("tp2a");
                        this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                        this.context.fillStyle = colorMap.get("tp2b");
                        this.context.fillRect(left + col * 20 + 10, top + row * 20, 10, 10);
                        this.context.fillRect(left + col * 20, top + row * 20 + 10, 10, 10);
                        break;
                    case 12: 
                        if (!this.toggleTilesState[0]) {
                            this.context.fillStyle = colorMap.get("disappearing1a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing1d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 13:
                        if (!this.toggleTilesState[1]) {
                            this.context.fillStyle = colorMap.get("disappearing2a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing2d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 14:
                        if (!this.toggleTilesState[2]) {
                            this.context.fillStyle = colorMap.get("disappearing3a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing3d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                }
            }
        }
        // Render the player
        this.context.fillStyle = colorMap.get("hero");
        this.context.fillRect(this.player.x, this.player.y, 20, 20);

        // If game is paused or the game is won, grey out the game
        if (this.paused || this.level > maxLevel) {
            this.context.fillStyle = colorMap.get("pause");
            this.context.fillRect(0, 0, this.canvas.width, this.canvas.height);
        }
    }

    /**
     * Sets (or resets) all level information, including the player's location and velocity.
     */
    _setLevel() {
        if (this.level > maxLevel) {
            return;
        }
        this.levelGrid = levels[this.level].grid;
        this.player.x = levels[this.level].spawn.x * 20;
        this.player.y = levels[this.level].spawn.y * 20;
        this.player.xVel = 0;
        this.player.yVel = 0;
        this.player.dir = "none";
        this.toggleTilesState = levels[this.level].toggleTilesState.slice();

        this.levelNumber.textContent = "Level: " + this.level;
        this.levelTitle.textContent = levels[this.level].title;
    }
    
    /**
     * Finds if the specified key has been pressed and not processed, meaning
     * it has not been checked after being set.
     * 
     * @param {String} code the keyCode of the key being searched for
     * @returns {Boolean} whether or not the key is pressed and not processed
     */
    _keyPressed(code) {
        if (this.unprocessedKeys.get(code)) {
            this.unprocessedKeys.set(code, false);
            return true;
        }
        return false;
    }

    /**
     * Checks for a command in the specified direction.
     * 
     * @param {String} direction direction to check for a command. Should be "up", "down", "left", or "right"
     * @returns {Boolean} whether or not the specified direction has a command waiting
     */
    _hasCommand(direction) {
        let temp = this.mouseCommand;
        if (this.mouseCommand == direction) {
            this.mouseCommand = "none";
        }
        
        if (direction == "left") {
            return this._keyPressed("KeyA") || this._keyPressed("ArrowLeft") || temp == "left";
        }
        else if (direction == "right") {
            return this._keyPressed("KeyD") || this._keyPressed("ArrowRight") || temp == "right";
        }
        else if (direction == "down") {
            return this._keyPressed("KeyS") || this._keyPressed("ArrowDown") || temp == "down";
        }
        else if (direction == "up") {
            return this._keyPressed("KeyW") || this._keyPressed("ArrowUp") || temp == "up";
        }
    }
}

//////////////////////////
// Main Execution Point //
//////////////////////////

/**
 * Game and levels are global so they aren't garbage collected, 
 * and levels are loaded as soon as possible. Levels also functions 
 * as a map, i.e. level 1 is stored at index 1.
 */
let game;
let levels = [];
for (let i = 1; i <= maxLevel; i++) {
    fetch("assets/levels/Level" + i + ".txt").then(function(response) {
        response.text().then(function(text) {
            console.log("Level " + i + " recieved.");
            levels[i] = new Level(text);
        });
    });
}

/** 
 * Function to be called on when the DOM is loaded, so that game can correctly add canvas.
 * Will periodically check to see if all levels are loaded, and when they are, the game is created.
 */
function startGame() {
    setTimeout(function check() {
        if (levels.length-1 == maxLevel) {
            game = new SliderGame();
            game.start();
        }
        else {
            setTimeout(check, 500);
        }
    }, 0);
}
document.addEventListener("DOMContentLoaded", startGame());