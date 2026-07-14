
---

## 💻 Код на 7 языках

### 1. Python – `runner.py`

```python
#!/usr/bin/env python3
# runner.py - Бесконечный бегун на Python

import os
import sys
import time
import random
import json
import threading
from dataclasses import dataclass
from typing import List, Optional

try:
    import keyboard
    from colorama import init, Fore, Style
    init(autoreset=True)
except ImportError:
    print("Установите зависимости: pip install keyboard colorama")
    sys.exit(1)

# Константы
WIDTH = 40
HEIGHT = 10
GROUND = HEIGHT - 1
PLAYER = '@'
OBSTACLE = '#'
EMPTY = ' '
RECORD_FILE = 'runner_record.json'

class RunnerGame:
    def __init__(self):
        self.width = WIDTH
        self.height = HEIGHT
        self.ground = GROUND
        self.player_x = 5
        self.player_y = GROUND
        self.player_vy = 0       # вертикальная скорость
        self.on_ground = True
        self.obstacles: List[int] = []  # x-координаты препятствий
        self.score = 0
        self.speed = 1
        self.game_over = False
        self.paused = False
        self.running = True
        self.record = self.load_record()
        self.lock = threading.Lock()

    def load_record(self) -> int:
        try:
            with open(RECORD_FILE, 'r') as f:
                data = json.load(f)
                return data.get('record', 0)
        except:
            return 0

    def save_record(self):
        with open(RECORD_FILE, 'w') as f:
            json.dump({'record': self.record}, f)

    def jump(self):
        if self.on_ground and not self.game_over and not self.paused:
            self.player_vy = -3.0
            self.on_ground = False

    def spawn_obstacle(self):
        # Препятствия появляются справа
        self.obstacles.append(self.width - 1)

    def update_physics(self):
        # Гравитация
        self.player_vy += 0.2
        self.player_y += self.player_vy

        # Проверка земли
        if self.player_y >= self.ground:
            self.player_y = self.ground
            self.player_vy = 0
            self.on_ground = True

        # Не выше потолка
        if self.player_y < 0:
            self.player_y = 0
            self.player_vy = 0

    def update(self):
        if self.game_over or self.paused:
            return

        # Движение препятствий
        for i in range(len(self.obstacles)):
            self.obstacles[i] -= 1
        # Удалить ушедшие за левую границу
        self.obstacles = [x for x in self.obstacles if x >= 0]

        # Столкновение (если препятствие на земле или выше, а игрок не в воздухе)
        for x in self.obstacles:
            # Проверяем, совпадает ли позиция препятствия с позицией игрока
            if x == self.player_x and self.player_y >= self.ground - 1:  # игрок на земле
                self.game_over = True
                if self.score > self.record:
                    self.record = self.score
                    self.save_record()
                return
            # Если препятствие на уровне выше, но игрок не допрыгнул
            # Для простоты проверяем совпадение по x и что игрок не на земле
            if x == self.player_x and self.player_y == self.ground - 1:
                self.game_over = True
                if self.score > self.record:
                    self.record = self.score
                    self.save_record()
                return

        # Счёт
        self.score += 1
        self.speed = 1 + self.score // 10

        # Спавн новых препятствий
        if random.random() < 0.2 + min(0.4, self.score / 200):
            self.spawn_obstacle()

        # Обновление физики игрока
        self.update_physics()

    def draw(self):
        os.system('cls' if os.name == 'nt' else 'clear')
        print('═' * self.width)
        print(f'  Счёт: {self.score}   Рекорд: {self.record}   Скорость: {self.speed}')
        print('═' * self.width)

        # Создаём сетку
        grid = [[' ' for _ in range(self.width)] for _ in range(self.height)]
        # Размещаем препятствия
        for x in self.obstacles:
            if 0 <= x < self.width:
                grid[self.ground][x] = Fore.RED + OBSTACLE + Style.RESET_ALL
        # Размещаем игрока
        if 0 <= self.player_y < self.height:
            grid[int(self.player_y)][self.player_x] = Fore.GREEN + PLAYER + Style.RESET_ALL

        # Вывод
        for row in grid:
            print('│' + ''.join(row) + '│')
        print('═' * self.width)
        status = "ПАУЗА" if self.paused else "ИГРА"
        print(f'  {status}  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход')

    def handle_input(self):
        while self.running:
            try:
                if keyboard.is_pressed('space') or keyboard.is_pressed('up'):
                    with self.lock:
                        self.jump()
                elif keyboard.is_pressed('p'):
                    with self.lock:
                        self.paused = not self.paused
                    time.sleep(0.2)
                elif keyboard.is_pressed('r'):
                    with self.lock:
                        self.reset()
                    time.sleep(0.2)
                elif keyboard.is_pressed('q'):
                    self.running = False
                    break
            except:
                pass
            time.sleep(0.05)

    def reset(self):
        self.player_y = self.ground
        self.player_vy = 0
        self.on_ground = True
        self.obstacles = []
        self.score = 0
        self.speed = 1
        self.game_over = False
        self.paused = False

    def run(self):
        # Запуск потока для ввода
        input_thread = threading.Thread(target=self.handle_input, daemon=True)
        input_thread.start()

        # Основной игровой цикл (60 FPS)
        last_update = time.time()
        while self.running:
            now = time.time()
            if now - last_update >= 1.0 / 60:
                with self.lock:
                    self.update()
                    self.draw()
                last_update = now
            time.sleep(0.01)

if __name__ == "__main__":
    game = RunnerGame()
    try:
        game.run()
    except KeyboardInterrupt:
        print("\nВыход...")
