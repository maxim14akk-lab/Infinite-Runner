// runner.cpp - Бесконечный бегун на C++17
#include <iostream>
#include <vector>
#include <random>
#include <thread>
#include <chrono>
#include <fstream>
#include <string>
#include <algorithm>
#include <cstdlib>

#ifdef _WIN32
    #include <conio.h>
    #include <windows.h>
    #define CLEAR() system("cls")
#else
    #include <ncurses.h>
    #include <termios.h>
    #include <unistd.h>
    #include <fcntl.h>
    #define CLEAR() system("clear")
#endif

const int WIDTH = 40;
const int HEIGHT = 10;
const int GROUND = HEIGHT - 1;
const char PLAYER = '@';
const char OBSTACLE = '#';
const std::string RECORD_FILE = "runner_record.json";

class RunnerGame {
public:
    RunnerGame() : playerX(5), playerY(GROUND), playerVy(0), onGround(true),
                   score(0), speed(1), gameOver(false), paused(false), running(true),
                   rng(std::random_device{}()) {
        record = loadRecord();
    }

    ~RunnerGame() {
#ifndef _WIN32
        endwin();
#endif
    }

    void run() {
#ifdef _WIN32
        // Для Windows используем _kbhit
        while (running) {
            handleInputWindows();
            update();
            draw();
            std::this_thread::sleep_for(std::chrono::milliseconds(16));
        }
#else
        initscr();
        raw();
        noecho();
        keypad(stdscr, TRUE);
        nodelay(stdscr, TRUE);
        curs_set(0);

        while (running) {
            handleInputNcurses();
            update();
            drawNcurses();
            napms(16);
        }
        endwin();
#endif
    }

private:
    int playerX;
    double playerY, playerVy;
    bool onGround;
    std::vector<int> obstacles;
    int score;
    int speed;
    bool gameOver, paused, running;
    int record;
    std::mt19937 rng;

    int loadRecord() {
        std::ifstream in(RECORD_FILE);
        if (!in) return 0;
        std::string content((std::istreambuf_iterator<char>(in)), std::istreambuf_iterator<char>());
        size_t pos = content.find("\"record\":");
        if (pos != std::string::npos) {
            pos += 9;
            size_t end = content.find(",", pos);
            if (end == std::string::npos) end = content.find("}", pos);
            return std::stoi(content.substr(pos, end - pos));
        }
        return 0;
    }

    void saveRecord() {
        std::ofstream out(RECORD_FILE);
        out << "{\"record\":" << record << "}";
    }

    void jump() {
        if (onGround && !gameOver && !paused) {
            playerVy = -3.0;
            onGround = false;
        }
    }

    void spawnObstacle() {
        obstacles.push_back(WIDTH - 1);
    }

    void updatePhysics() {
        playerVy += 0.2;
        playerY += playerVy;
        if (playerY >= GROUND) {
            playerY = GROUND;
            playerVy = 0;
            onGround = true;
        }
        if (playerY < 0) {
            playerY = 0;
            playerVy = 0;
        }
    }

    void update() {
        if (gameOver || paused) return;

        // Движение препятствий
        for (int& x : obstacles) x--;
        obstacles.erase(std::remove_if(obstacles.begin(), obstacles.end(),
            [](int x) { return x < 0; }), obstacles.end());

        // Столкновение
        int playerFloor = (int)playerY;
        for (int x : obstacles) {
            if (x == playerX && (playerFloor >= GROUND - 1 || playerFloor == GROUND - 1)) {
                gameOver = true;
                if (score > record) {
                    record = score;
                    saveRecord();
                }
                return;
            }
        }

        score++;
        speed = 1 + score / 10;

        double prob = 0.2 + std::min(0.4, score / 200.0);
        if ((double)rand() / RAND_MAX < prob) spawnObstacle();

        updatePhysics();
    }

    void draw() {
        CLEAR();
        std::cout << std::string(WIDTH, '═') << std::endl;
        std::cout << "  Счёт: " << score << "   Рекорд: " << record << "   Скорость: " << speed << std::endl;
        std::cout << std::string(WIDTH, '═') << std::endl;

        char grid[HEIGHT][WIDTH];
        for (int y = 0; y < HEIGHT; ++y)
            for (int x = 0; x < WIDTH; ++x)
                grid[y][x] = ' ';
        for (int x : obstacles) {
            if (x >= 0 && x < WIDTH) grid[GROUND][x] = OBSTACLE;
        }
        int y = (int)playerY;
        if (y >= 0 && y < HEIGHT) grid[y][playerX] = PLAYER;

        for (int row = 0; row < HEIGHT; ++row) {
            std::cout << '│';
            for (int col = 0; col < WIDTH; ++col)
                std::cout << grid[row][col];
            std::cout << '│' << std::endl;
        }
        std::cout << std::string(WIDTH, '═') << std::endl;
        std::cout << "  " << (paused ? "ПАУЗА" : "ИГРА") << "  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход" << std::endl;
    }

#ifdef _WIN32
    void handleInputWindows() {
        if (_kbhit()) {
            int ch = _getch();
            if (ch == 224) { // стрелки
                ch = _getch();
                if (ch == 72) jump(); // up
            } else {
                switch (tolower(ch)) {
                    case ' ': jump(); break;
                    case 'p': paused = !paused; break;
                    case 'r': reset(); break;
                    case 'q': running = false; break;
                }
            }
        }
    }
#else
    void handleInputNcurses() {
        int ch = getch();
        if (ch == ERR) return;
        switch (ch) {
            case ' ':
            case KEY_UP: jump(); break;
            case 'p': case 'P': paused = !paused; break;
            case 'r': case 'R': reset(); break;
            case 'q': case 'Q': running = false; break;
        }
    }
#endif

    void reset() {
        playerY = GROUND;
        playerVy = 0;
        onGround = true;
        obstacles.clear();
        score = 0;
        speed = 1;
        gameOver = false;
        paused = false;
    }
};

int main() {
    RunnerGame game;
    game.run();
    return 0;
}
