// Runner.java - Бесконечный бегун на Java
import org.jline.terminal.Terminal;
import org.jline.terminal.TerminalBuilder;
import org.jline.utils.InfoCmp;
import com.google.gson.*;

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.*;

public class Runner {
    private static final int WIDTH = 40;
    private static final int HEIGHT = 10;
    private static final int GROUND = HEIGHT - 1;
    private static final char PLAYER = '@';
    private static final char OBSTACLE = '#';
    private static final String RECORD_FILE = "runner_record.json";

    private Terminal terminal;
    private int playerX;
    private double playerY;
    private double playerVy;
    private boolean onGround;
    private List<Integer> obstacles;
    private int score;
    private int speed;
    private boolean gameOver;
    private boolean paused;
    private boolean running;
    private int record;
    private long lastUpdate;
    private ScheduledExecutorService scheduler;

    public Runner() throws IOException {
        terminal = TerminalBuilder.builder().system(true).build();
        playerX = 5;
        playerY = GROUND;
        playerVy = 0;
        onGround = true;
        obstacles = new ArrayList<>();
        score = 0;
        speed = 1;
        gameOver = false;
        paused = false;
        running = true;
        record = loadRecord();
        lastUpdate = System.currentTimeMillis();
        scheduler = Executors.newSingleThreadScheduledExecutor();
    }

    private int loadRecord() {
        try {
            String content = new String(Files.readAllBytes(Paths.get(RECORD_FILE)));
            JsonObject obj = new Gson().fromJson(content, JsonObject.class);
            return obj.get("record").getAsInt();
        } catch (Exception e) {
            return 0;
        }
    }

    private void saveRecord() {
        try {
            JsonObject obj = new JsonObject();
            obj.addProperty("record", record);
            Files.write(Paths.get(RECORD_FILE), obj.toString().getBytes());
        } catch (Exception e) {}
    }

    private void jump() {
        if (onGround && !gameOver && !paused) {
            playerVy = -3.0;
            onGround = false;
        }
    }

    private void spawnObstacle() {
        obstacles.add(WIDTH - 1);
    }

    private void updatePhysics() {
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

    private void update() {
        if (gameOver || paused) return;

        // Движение препятствий
        Iterator<Integer> it = obstacles.iterator();
        while (it.hasNext()) {
            int x = it.next();
            x--;
            if (x < 0) it.remove();
            else {
                // обновляем в списке (через set)
                int idx = obstacles.indexOf(x); // неправильно, лучше использовать ListIterator
            }
        }
        // Используем ListIterator
        ListIterator<Integer> lit = obstacles.listIterator();
        while (lit.hasNext()) {
            int x = lit.next();
            x--;
            lit.set(x);
            if (x < 0) lit.remove();
        }

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

        if (Math.random() < 0.2 + Math.min(0.4, score / 200.0)) {
            spawnObstacle();
        }

        updatePhysics();
    }

    private void draw() {
        terminal.puts(InfoCmp.Capability.clear_screen);
        System.out.println("═".repeat(WIDTH));
        System.out.printf("  Счёт: %d   Рекорд: %d   Скорость: %d%n", score, record, speed);
        System.out.println("═".repeat(WIDTH));

        char[][] grid = new char[HEIGHT][WIDTH];
        for (int y = 0; y < HEIGHT; y++) Arrays.fill(grid[y], ' ');
        // Препятствия
        for (int x : obstacles) {
            if (x >= 0 && x < WIDTH) grid[GROUND][x] = OBSTACLE;
        }
        // Игрок
        int y = (int)playerY;
        if (y >= 0 && y < HEIGHT) grid[y][playerX] = PLAYER;

        for (int row = 0; row < HEIGHT; row++) {
            System.out.print('│');
            for (int col = 0; col < WIDTH; col++) {
                System.out.print(grid[row][col]);
            }
            System.out.println('│');
        }
        System.out.println("═".repeat(WIDTH));
        String status = paused ? "ПАУЗА" : "ИГРА";
        System.out.printf("  %s  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход%n", status);
    }

    private void reset() {
        playerY = GROUND;
        playerVy = 0;
        onGround = true;
        obstacles.clear();
        score = 0;
        speed = 1;
        gameOver = false;
        paused = false;
    }

    private void handleInput() {
        try {
            while (running) {
                int ch = terminal.reader().read();
                if (ch == -1) continue;
                char c = (char) ch;
                switch (c) {
                    case ' ': case '↑': jump(); break;
                    case 'p': case 'P': paused = !paused; break;
                    case 'r': case 'R': reset(); break;
                    case 'q': case 'Q': running = false; break;
                }
            }
        } catch (IOException e) {}
    }

    public void run() throws Exception {
        Thread inputThread = new Thread(this::handleInput);
        inputThread.setDaemon(true);
        inputThread.start();

        while (running) {
            long now = System.currentTimeMillis();
            if (now - lastUpdate >= 1000 / 60) {
                update();
                draw();
                lastUpdate = now;
            }
            Thread.sleep(16);
        }
        terminal.close();
        scheduler.shutdown();
    }

    public static void main(String[] args) throws Exception {
        Runner game = new Runner();
        game.run();
    }
}
