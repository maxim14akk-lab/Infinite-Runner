🏃 Infinite Runner – консольный бесконечный бегун

**Динамичная аркадная игра** на **7 языках программирования**.  
Управляйте бегуном, который мчится по бесконечной дороге, уворачиваясь от препятствий.  
Чем дальше – тем быстрее! Игра сохраняет рекорд, имеет анимацию и плавное управление.

---

## 🎮 Геймплей

- Бегун (`@`) бежит по дороге, постоянно ускоряясь.
- Препятствия (`#`) появляются справа и движутся влево.
- **Цель** – пробежать как можно дальше, не столкнувшись.
- **Управление**: `Пробел` или `↑` (стрелка вверх) – прыжок.
- **Дополнительно**: `P` – пауза, `R` – рестарт, `Q` – выход.
- Счёт увеличивается с каждым шагом; скорость растёт каждые 10 очков.
- Рекорд сохраняется в файл и загружается при старте.

---

## 🖥️ Пример игрового экрана
═══════════════════════════════════════
Счёт: 42 Рекорд: 50 Скорость: 3
═══════════════════════════════════════
│ │
│ │
│ # │
│ # │
│ # │
│ @ │
═══════════════════════════════════════
Пробел/↑ - прыжок | P - пауза | R - рестарт | Q - выход

text

---

## 🚀 Запуск

| Язык       | Файл          | Команда запуска                              | Зависимости                |
|------------|---------------|----------------------------------------------|----------------------------|
| Python     | `runner.py`   | `python runner.py`                           | `keyboard`, `colorama`     |
| JavaScript | `runner.js`   | `node runner.js`                             | `keypress`                 |
| Java       | `Runner.java` | `javac Runner.java && java Runner`           | `jline`                    |
| C++        | `runner.cpp`  | `g++ -std=c++17 runner.cpp -o runner && ./runner` | (Windows: `conio` / Unix: `ncurses`) |
| C#         | `runner.cs`   | `csc runner.cs && runner.exe`                | .NET SDK                   |
| Go         | `runner.go`   | `go run runner.go`                           | `keyboard`                 |
| Rust       | `runner.rs`   | `cargo run`                                  | `crossterm`, `rand`, `serde_json` |

---

## 📦 Установка зависимостей

### Python
```bash
pip install keyboard colorama
JavaScript (Node.js)
bash
npm install keypress
Java
Скачайте jline-3.21.0.jar и gson-2.8.9.jar (для сохранения рекорда) и добавьте в classpath.

C++
Windows: используется conio.h (встроена).

Unix/Linux: установите ncurses:

bash
sudo apt-get install libncurses-dev   # Debian/Ubuntu
sudo yum install ncurses-devel        # RHEL
Go
bash
go get github.com/eiannone/keyboard
Rust
В Cargo.toml добавьте:

toml
[dependencies]
crossterm = "0.27"
rand = "0.8"
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
🛠️ Продвинутые функции
Анимация прыжка – параболическая траектория (плавное изменение высоты).

Ускорение – скорость растёт с каждым пройденным шагом.

Рекорды – лучший результат сохраняется в JSON-файл.

Цветной вывод – препятствия красные, бегун зелёный.

Пауза/Рестарт – полное управление игровым процессом.

Кросс-платформенность – все версии работают на Windows, Linux и macOS.

📁 Структура репозитория
text
/
├── README.md
├── runner.py
├── runner.js
├── Runner.java
├── runner.cpp
├── runner.cs
├── runner.go
├── runner.rs
└── (для Rust) Cargo.toml + src/main.rs
📜 Лицензия
MIT – свободно используйте, улучшайте и распространяйте.
