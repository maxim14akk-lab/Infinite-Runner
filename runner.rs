// runner.rs - Бесконечный бегун на Rust
use crossterm::{
    event::{self, Event, KeyCode},
    execute,
    terminal::{self, Clear, ClearType},
};
use std::io::{stdout, Write};
use std::time::{Duration, Instant};
use std::collections::VecDeque;
use serde::{Serialize, Deserialize};
use std::fs;
use rand::Rng;

const WIDTH: usize = 40;
const HEIGHT: usize = 10;
const GROUND: usize = HEIGHT - 1;
const PLAYER: char = '@';
const OBSTACLE: char = '#';
const RECORD_FILE: &str = "runner_record.json";

#[derive(Serialize, Deserialize)]
struct RecordData {
    record: u32,
}

struct RunnerGame {
    player_x: usize,
    player_y: f64,
    player_vy: f64,
    on_ground: bool,
    obstacles: VecDeque<usize>,
    score: u32,
    speed: u32,
    game_over: bool,
    paused: bool,
    running: bool,
    record: u32,
    rng: rand::rngs::ThreadRng,
}

impl RunnerGame {
    fn new() -> Self {
        Self {
            player_x: 5,
            player_y: GROUND as f64,
            player_vy: 0.0,
            on_ground: true,
            obstacles: VecDeque::new(),
            score: 0,
            speed: 1,
            game_over: false,
            paused: false,
            running: true,
            record: Self::load_record(),
            rng: rand::thread_rng(),
        }
    }

    fn load_record() -> u32 {
        if let Ok(data) = fs::read_to_string(RECORD_FILE) {
            if let Ok(rec) = serde_json::from_str::<RecordData>(&data) {
                return rec.record;
            }
        }
        0
    }

    fn save_record(record: u32) {
        let data = RecordData { record };
        let _ = fs::write(RECORD_FILE, serde_json::to_string(&data).unwrap());
    }

    fn jump(&mut self) {
        if self.on_ground && !self.game_over && !self.paused {
            self.player_vy = -3.0;
            self.on_ground = false;
        }
    }

    fn spawn_obstacle(&mut self) {
        self.obstacles.push_back(WIDTH - 1);
    }

    fn update_physics(&mut self) {
        self.player_vy += 0.2;
        self.player_y += self.player_vy;
        if self.player_y >= GROUND as f64 {
            self.player_y = GROUND as f64;
            self.player_vy = 0.0;
            self.on_ground = true;
        }
        if self.player_y < 0.0 {
            self.player_y = 0.0;
            self.player_vy = 0.0;
        }
    }

    fn update(&mut self) {
        if self.game_over || self.paused {
            return;
        }

        // Движение препятствий
        for x in self.obstacles.iter_mut() {
            if *x > 0 {
                *x -= 1;
            }
        }
        self.obstacles.retain(|&x| x > 0);

        // Столкновение
        let player_floor = self.player_y as usize;
        for &x in &self.obstacles {
            if x == self.player_x && (player_floor >= GROUND - 1 || player_floor == GROUND - 1) {
                self.game_over = true;
                if self.score > self.record {
                    self.record = self.score;
                    Self::save_record(self.record);
                }
                return;
            }
        }

        self.score += 1;
        self.speed = 1 + self.score / 10;

        let prob = 0.2 + f64::min(0.4, self.score as f64 / 200.0);
        if self.rng.gen_bool(prob) {
            self.spawn_obstacle();
        }

        self.update_physics();
    }

    fn draw(&self) {
        execute!(stdout(), Clear(ClearType::All)).unwrap();
        let mut out = stdout();
        writeln!(out, "{}", "═".repeat(WIDTH)).unwrap();
        writeln!(out, "  Счёт: {}   Рекорд: {}   Скорость: {}", self.score, self.record, self.speed).unwrap();
        writeln!(out, "{}", "═".repeat(WIDTH)).unwrap();

        let mut grid = vec![vec![' '; WIDTH]; HEIGHT];
        for &x in &self.obstacles {
            if x < WIDTH {
                grid[GROUND][x] = OBSTACLE;
            }
        }
        let y = self.player_y as usize;
        if y < HEIGHT {
            grid[y][self.player_x] = PLAYER;
        }

        for row in grid {
            write!(out, "│").unwrap();
            for ch in row {
                write!(out, "{}", ch).unwrap();
            }
            writeln!(out, "│").unwrap();
        }
        writeln!(out, "{}", "═".repeat(WIDTH)).unwrap();
        let status = if self.paused { "ПАУЗА" } else { "ИГРА" };
        writeln!(out, "  {}  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход", status).unwrap();
        out.flush().unwrap();
    }

    fn reset(&mut self) {
        self.player_y = GROUND as f64;
        self.player_vy = 0.0;
        self.on_ground = true;
        self.obstacles.clear();
        self.score = 0;
        self.speed = 1;
        self.game_over = false;
        self.paused = false;
    }

    fn run(&mut self) {
        terminal::enable_raw_mode().unwrap();
        let mut last_update = Instant::now();

        while self.running {
            // Обработка ввода
            if event::poll(Duration::from_millis(50)).unwrap() {
                if let Event::Key(key) = event::read().unwrap() {
                    match key.code {
                        KeyCode::Up | KeyCode::Char(' ') => self.jump(),
                        KeyCode::Char('p') | KeyCode::Char('P') => self.paused = !self.paused,
                        KeyCode::Char('r') | KeyCode::Char('R') => self.reset(),
                        KeyCode::Char('q') | KeyCode::Char('Q') => {
                            self.running = false;
                            break;
                        }
                        _ => {}
                    }
                }
            }

            let now = Instant::now();
            if now - last_update >= Duration::from_secs_f64(1.0 / 60.0) {
                self.update();
                self.draw();
                last_update = now;
            }
        }

        terminal::disable_raw_mode().unwrap();
    }
}

fn main() {
    let mut game = RunnerGame::new();
    game.run();
    println!("Игра завершена.");
}
