// runner.js - Бесконечный бегун на JavaScript (Node.js)
const fs = require('fs');
const keypress = require('keypress');
const readline = require('readline');

const WIDTH = 40;
const HEIGHT = 10;
const GROUND = HEIGHT - 1;
const PLAYER = '@';
const OBSTACLE = '#';
const RECORD_FILE = 'runner_record.json';

class RunnerGame {
    constructor() {
        this.width = WIDTH;
        this.height = HEIGHT;
        this.ground = GROUND;
        this.playerX = 5;
        this.playerY = GROUND;
        this.playerVy = 0;
        this.onGround = true;
        this.obstacles = [];
        this.score = 0;
        this.speed = 1;
        this.gameOver = false;
        this.paused = false;
        this.running = true;
        this.record = this.loadRecord();
        this.lastUpdate = Date.now();
        this.timer = null;
    }

    loadRecord() {
        try {
            const data = fs.readFileSync(RECORD_FILE, 'utf8');
            return JSON.parse(data).record || 0;
        } catch { return 0; }
    }

    saveRecord() {
        fs.writeFileSync(RECORD_FILE, JSON.stringify({ record: this.record }));
    }

    jump() {
        if (this.onGround && !this.gameOver && !this.paused) {
            this.playerVy = -3.0;
            this.onGround = false;
        }
    }

    spawnObstacle() {
        this.obstacles.push(this.width - 1);
    }

    updatePhysics() {
        this.playerVy += 0.2;
        this.playerY += this.playerVy;
        if (this.playerY >= this.ground) {
            this.playerY = this.ground;
            this.playerVy = 0;
            this.onGround = true;
        }
        if (this.playerY < 0) {
            this.playerY = 0;
            this.playerVy = 0;
        }
    }

    update() {
        if (this.gameOver || this.paused) return;

        // Движение препятствий
        for (let i = this.obstacles.length - 1; i >= 0; i--) {
            this.obstacles[i]--;
            if (this.obstacles[i] < 0) this.obstacles.splice(i, 1);
        }

        // Столкновение
        const playerFloor = Math.floor(this.playerY);
        for (const x of this.obstacles) {
            if (x === this.playerX && (playerFloor >= this.ground - 1 || playerFloor === this.ground - 1)) {
                this.gameOver = true;
                if (this.score > this.record) {
                    this.record = this.score;
                    this.saveRecord();
                }
                return;
            }
        }

        this.score++;
        this.speed = 1 + Math.floor(this.score / 10);

        if (Math.random() < 0.2 + Math.min(0.4, this.score / 200)) {
            this.spawnObstacle();
        }

        this.updatePhysics();
    }

    draw() {
        console.clear();
        console.log('═'.repeat(this.width));
        console.log(`  Счёт: ${this.score}   Рекорд: ${this.record}   Скорость: ${this.speed}`);
        console.log('═'.repeat(this.width));

        const grid = Array.from({ length: this.height }, () => Array(this.width).fill(' '));
        // Препятствия
        for (const x of this.obstacles) {
            if (x >= 0 && x < this.width) {
                grid[this.ground][x] = OBSTACLE;
            }
        }
        // Игрок
        const y = Math.floor(this.playerY);
        if (y >= 0 && y < this.height) {
            grid[y][this.playerX] = PLAYER;
        }

        for (const row of grid) {
            console.log('│' + row.join('') + '│');
        }
        console.log('═'.repeat(this.width));
        const status = this.paused ? 'ПАУЗА' : 'ИГРА';
        console.log(`  ${status}  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход`);
    }

    reset() {
        this.playerY = this.ground;
        this.playerVy = 0;
        this.onGround = true;
        this.obstacles = [];
        this.score = 0;
        this.speed = 1;
        this.gameOver = false;
        this.paused = false;
    }

    handleKey(ch, key) {
        if (!key) return;
        if (key.name === 'space' || key.name === 'up') {
            this.jump();
        } else if (key.name === 'p') {
            this.paused = !this.paused;
        } else if (key.name === 'r') {
            this.reset();
        } else if (key.name === 'q') {
            this.running = false;
            process.stdin.pause();
            process.exit(0);
        }
    }

    run() {
        // Настройка keypress
        keypress(process.stdin);
        process.stdin.setRawMode(true);
        process.stdin.resume();
        process.stdin.on('keypress', (ch, key) => this.handleKey(ch, key));

        const gameLoop = () => {
            if (!this.running) return;
            const now = Date.now();
            if (now - this.lastUpdate >= 1000 / 60) {
                this.update();
                this.draw();
                this.lastUpdate = now;
            }
            this.timer = setTimeout(gameLoop, 16);
        };
        gameLoop();
    }
}

const game = new RunnerGame();
game.run();
