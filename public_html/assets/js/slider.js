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

// Pixels to move the player by in each update of the frame
const PLAYER_SPEED = 5;
const MAX_LEVEL = 11;
const FPS = 50;

/** 
 * Specific class to run the Slider Game. The game is a puzzle game
 * where the player's goal is to slide their square to the green
 * goal tile while avoiding obstacles and using other game
 * mechanics to their advantage. The player should beware, however:
 * Once started, the square will not stop moving until it hits a
 * solid block! Early levels are designed to teach concepts, while
 * later levels are designed to be more of a challenge. Average
 * playtime for a new player is between 20 and 30 minutes.
 * 
 * Requires a div id "game_container" to place the game inside of.
 */
class SliderGame {
    /**
     * Creates the Slider Game, including the canvas
     */
    constructor(levelData) {
        // Slider Game is run at 50 fps, which is intense for some computers, but at fewer frames
        // runs fairly choppily. I will leave this for now, and maybe revise later

        // Add game document elements
        this._initGameElements();

        // All local variables
        this.levelData = levelData;
        this.levelNum = 0;
        this.accumulatedTime = 0;
        this.lastDate = Date.now();
        this.startHours = new Date(0).getHours();
        this.paused = false;
        this.player = {
            get x() {
                return this._x;
            },
            set x(x) {
                this._x = x; 
                this.display.style.left = x + 'px';
            },
            get y() {
                return this._y;
            },
            set y(y) {
                this._y = y; 
                this.display.style.top = y + 'px';
            }, 
            xVel: 0,
            yVel: 0, 
            dir: 'none', 
            startedSliding: false,
            display: this.playerDisplay
        };
        this.player.x = 0;
        this.player.y = 0;

        // Add control variables and event listeners
        this.unprocessedKeys = new Map();
        document.addEventListener('keydown', (e) => {
            e.preventDefault();
            this.unprocessedKeys.set(e.code, true);
        });
        document.addEventListener('keyup', (e) => {
            e.preventDefault();
            this.unprocessedKeys.set(e.code, false);
        });

        // Initialize level variables
        this._setLevel(this.levelNum);
    }

    /**
     * Adds the DOM elements needed for the game into the game-container element
     */
    _initGameElements() {
        let gameContainer = document.getElementById('game-container');

        // Create the info bar that contains information that the user needs like level info
        // and the current time, and keep track of the elements to modify later
        let infobar = document.createElement('div');
        infobar.id = 'info-bar';
        infobar.append(this.levelNumberDisplay = document.createElement('div'));
        infobar.append(this.levelTitleDisplay = document.createElement('div'));
        infobar.append(this.timerDisplay = document.createElement('div'));
        gameContainer.append(infobar);

        let levelContainer = document.createElement('div');
        levelContainer.id = 'level-container';
        gameContainer.append(levelContainer);

        // Next set up the canvas and context, where the game will be drawn
        this.canvas = document.createElement('canvas');
        this.canvas.id = 'level';
        this.canvas.width = 500;
        this.canvas.height = 500;
        levelContainer.append(this.canvas);
        this.context = this.canvas.getContext('2d');

        // Put the player in the canvas- a 20px yellow square
        this.playerDisplay = document.createElement('div');
        this.playerDisplay.id = 'player'
        levelContainer.append(this.playerDisplay);

        // Set up pausing layer that will go overtop the player and level when paused
        this.pauseLayer = document.createElement('div');
        this.pauseLayer.id = 'pause';
        this.pauseLayer.style.visibility = 'hidden';
        levelContainer.append(this.pauseLayer);
    }

    /**
     * Starts the game loop if not already started, with an interval specified by the frames per second
     */
    start() {
        if (!this.interval) {
            this.interval = setInterval(this._update.bind(this), 1000 / FPS);
        }
    }

    /**
     * Updates the slider game components, including the paused state, 
     * time and player logic
     */
    _update() {
        // Pause or Unpause game
        if (this._keyPressed("Space")) {
            this.paused = !this.paused;
            this.pauseLayer.style.visibility = (this.paused) ? 'visible' : 'hidden';
        }

        this._updateTime();
        
        // If the game is won, check if the player wants to restart game
        if (this.levelNum > MAX_LEVEL) {
            this.pauseLayer.style.visibility = 'visible';
            if (this._keyPressed("KeyR")) {
                this.pauseLayer.style.visibility = 'hidden';
                this.accumulatedTime = 0;
                this.lastDate = Date.now();
                this.startHours = new Date(0).getHours();
                this.levelNum = 0;
                this._setLevel();
                return;
            }
        }
        else if (!this.paused) { // Otherwise, update the player logic
            if (this._updatePlayer()) {
                if (this.levelNum <= MAX_LEVEL) {
                    this._setLevel();
                }
            } 
        }
    }

    /**
     * Update the time value in the "timer" div
     */
    _updateTime() {
        // Determine what to append to time text content, and update time if required
        let append = "";
        if (this.levelNum > MAX_LEVEL) {
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
            this.timerDisplay.textContent = `Time: ${nowDate.getHours() - this.startHours}:`+
                                            `${twoDigitFormat.format(nowDate.getMinutes())}:` +
                                            `${twoDigitFormat.format(nowDate.getSeconds())}${append}`;
        }
        else {
            this.timerDisplay.textContent = `Time: ${twoDigitFormat.format(nowDate.getMinutes())}:` +
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
        if (this.player.dir == 'none') {
            this.player.xVel = 0;
            this.player.yVel = 0;  
            if (this._hasCommand('left')) {
                this.player.dir = 'left';
                this.player.startedSliding = true;
                this.player.xVel = -PLAYER_SPEED;
            }
            else if (this._hasCommand('down')) {
                this.player.dir = 'down';
                this.player.startedSliding = true;
                this.player.yVel = PLAYER_SPEED;
            }
            else if (this._hasCommand('right')) {
                this.player.dir = 'right';
                this.player.startedSliding = true;
                this.player.xVel = PLAYER_SPEED;
            }
            else if (this._hasCommand('up')) {
                this.player.dir = 'up';
                this.player.startedSliding = true;
                this.player.yVel = -PLAYER_SPEED;
            }
        }

        // This block is triggered if the player has perfectly passed over a grid square.
        // Any functionality that is dependent on landing on a square is triggered.
        // This is done before stopping because velocity is conserved in the event that
        // this functionality causes whatever would stop velocity to no longer stop it,
        // i.e. a button switches off a wall, or the player is teleported to a new location.
        if (!this.player.startedSliding && 
            (this.player.x / 20) % 1 == 0 && (this.player.y / 20) % 1 == 0 && 
            this.player.dir != 'none') {
            switch (this.level.grid[this.player.y / 20][this.player.x / 20]) {
                // 2 is lava- returning true resets the level
                case 2: 
                    return true;
                // 4 is orange buttons- toggle the tile state.
                case 4: 
                    this.level.toggleTilesState['orange'] = !this.level.toggleTilesState['orange'];
                    this._renderLevel();
                    break;
                // 6 is teal buttons- toggle the tile state
                case 6: 
                    this.level.toggleTilesState['teal'] = !this.level.toggleTilesState['teal'];
                    this._renderLevel();
                    break;
                // 8 is blue buttons- toggle the tile state
                case 8: 
                    this.level.toggleTilesState['blue'] = !this.level.toggleTilesState['blue'];
                    this._renderLevel();
                    break;
                // 10 is pink teleporters- move player to other teleporter location
                case 10: 
                    if (this.player.x / 20 == this.level.tpLocations['pink'][0].x &&
                        this.player.y / 20 == this.level.tpLocations['pink'][0].y) {
                        this.player.x = this.level.tpLocations['pink'][1].x * 20;
                        this.player.y = this.level.tpLocations['pink'][1].y * 20;
                    } else {
                        this.player.x = this.level.tpLocations['pink'][0].x * 20;
                        this.player.y = this.level.tpLocations['pink'][0].y * 20;
                    }
                    break;
                // 11 is purple teleporters- move player to other teleporter location
                case 11: 
                    if (this.player.x / 20 == this.level.tpLocations['purple'][0].x &&
                        this.player.y / 20 == this.level.tpLocations['purple'][0].y) {
                        this.player.x = this.level.tpLocations['purple'][1].x * 20;
                        this.player.y = this.level.tpLocations['purple'][1].y * 20;
                    } else {
                        this.player.x = this.level.tpLocations['purple'][0].x * 20;
                        this.player.y = this.level.tpLocations['purple'][0].y * 20;
                    }
                    break;
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
        if (this.player.dir == 'right') {
            xIndex = Math.floor((this.player.x + this.player.xVel + (20 - PLAYER_SPEED)) / 20);
        }
        else {
            xIndex = Math.floor((this.player.x + this.player.xVel) / 20);
        }
        if (this.player.dir == 'down') {
            yIndex = Math.floor((this.player.y + this.player.yVel + (20 - PLAYER_SPEED)) / 20);
        }
        else {
            yIndex = Math.floor((this.player.y + this.player.yVel) / 20);
        }

        // If the player is moving and the location they are moving to 
        // is valid, stop; otherwise update location
        if (this.player.dir != 'none') {
            if (xIndex >= 0 && xIndex < 25 && yIndex >= 0 && yIndex < 25) {
                if (this.level.grid[yIndex][xIndex] == 1 ||
                        this.level.toggleTilesState['orange']  && this.level.grid[yIndex][xIndex] == 5 ||
                        this.level.toggleTilesState['teal']    && this.level.grid[yIndex][xIndex] == 7 ||
                        this.level.toggleTilesState['blue']    && this.level.grid[yIndex][xIndex] == 9 ||
                        !this.level.toggleTilesState['orange'] && this.level.grid[yIndex][xIndex] == 12 ||
                        !this.level.toggleTilesState['teal']   && this.level.grid[yIndex][xIndex] == 13 ||
                        !this.level.toggleTilesState['blue']   && this.level.grid[yIndex][xIndex] == 14 ) {
                    this.player.dir = 'none';
                }
                else {
                    this.player.x += this.player.xVel;
                    this.player.y += this.player.yVel;
                }
            }
            else {
                this.player.dir = 'none';
            }
        }

        // The final check is for the goal. This is done here because the goal can only be said
        // to be reached if the player is stopped on it.
        if (this.player.dir === 'none' && this.level.grid[this.player.y / 20][this.player.x / 20] === 3) {
            this.levelNum++;
            return true;
        }

        // If we've made it this far, we don't need to advance the level
        return false;
    }

    /**
     * Renders an entire frame of the slider game by rendering each background tile,
     * and then the player on top of it, and then the foreground greying out if
     * the game has been won or paused.
     */
    _renderLevel() {
        // Render the background
        let left = 0;
        let top = 0;
        for (let row = 0; row < 25; row++) {
            for (let col = 0; col < 25; col++) {
                this.context.fillStyle = colorMap.get("border");
                this.context.fillRect(left + col * 20, top + row * 20, 20, 20);
                switch (this.level.grid[row][col]) {
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
                        if (this.level.toggleTilesState['orange']) {
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
                        if (this.level.toggleTilesState['teal']) {
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
                        if (this.level.toggleTilesState['blue']) {
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
                        if (!this.level.toggleTilesState['orange']) {
                            this.context.fillStyle = colorMap.get("disappearing1a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing1d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 13:
                        if (!this.level.toggleTilesState['teal']) {
                            this.context.fillStyle = colorMap.get("disappearing2a");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        else {
                            this.context.fillStyle = colorMap.get("disappearing2d");
                            this.context.fillRect(left + col * 20 + 1, top + row * 20 + 1, 18, 18);
                        }
                        break;
                    case 14:
                        if (!this.level.toggleTilesState['blue']) {
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
    }

    /**
     * Sets (or resets) all level information, including the player's location and velocity.
     */
    _setLevel() {
        // If the current level is invalid, don't try to set the level
        if (this.levelNum < 0 || this.levelNum > MAX_LEVEL) {
            return;
        }

        // (Re)set the level data and canvas
        this.level = JSON.parse(JSON.stringify(this.levelData[this.levelNum]));
        this._renderLevel();

        // (Re)set the player position and movement data
        this.player.x = this.level.spawn.x * 20;
        this.player.y = this.level.spawn.y * 20;
        this.player.xVel = 0;
        this.player.yVel = 0;
        this.player.dir = 'none';

        // Update the level displays
        this.levelNumberDisplay.textContent = `Level: ${this.levelNum + 1}`;
        this.levelTitleDisplay.textContent = this.level.title;
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
     * @param {String} direction direction to check for a command. Should be 'up', 'down', 'left', or 'right'
     * @returns {Boolean} whether or not the specified direction has a command waiting
     */
    _hasCommand(direction) {
        
        if (direction == 'left') {
            return this._keyPressed("KeyA") || this._keyPressed("ArrowLeft");
        }
        else if (direction == 'right') {
            return this._keyPressed("KeyD") || this._keyPressed("ArrowRight");
        }
        else if (direction == 'down') {
            return this._keyPressed("KeyS") || this._keyPressed("ArrowDown");
        }
        else if (direction == 'up') {
            return this._keyPressed("KeyW") || this._keyPressed("ArrowUp");
        }
    }
}

//////////////////////////
// Main Execution Point //
//////////////////////////

document.addEventListener("DOMContentLoaded", () => {
    let game = new SliderGame(JSON.parse(document.getElementById('level-data').innerHTML));
    game.start();
});
