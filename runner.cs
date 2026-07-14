// runner.cs - Бесконечный бегун на C#
using System;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Threading;

class RunnerGame
{
    private const int WIDTH = 40;
    private const int HEIGHT = 10;
    private const int GROUND = HEIGHT - 1;
    private const char PLAYER = '@';
    private const char OBSTACLE = '#';
    private const string RECORD_FILE = "runner_record.json";

    private int playerX = 5;
    private double playerY, playerVy;
    private bool onGround = true;
    private List<int> obstacles = new List<int>();
    private int score = 0;
    private int speed = 1;
    private bool gameOver = false;
    private bool paused = false;
    private bool running = true;
    private int record;
    private DateTime lastUpdate = DateTime.Now;
    private Random rand = new Random();

    public RunnerGame()
    {
        record = LoadRecord();
    }

    private int LoadRecord()
    {
        if (!File.Exists(RECORD_FILE)) return 0;
        string json = File.ReadAllText(RECORD_FILE);
        var data = JsonSerializer.Deserialize<Dictionary<string, int>>(json);
        return data != null && data.ContainsKey("record") ? data["record"] : 0;
    }

    private void SaveRecord()
    {
        var data = new Dictionary<string, int> { { "record", record } };
        File.WriteAllText(RECORD_FILE, JsonSerializer.Serialize(data));
    }

    private void Jump()
    {
        if (onGround && !gameOver && !paused)
        {
            playerVy = -3.0;
            onGround = false;
        }
    }

    private void SpawnObstacle()
    {
        obstacles.Add(WIDTH - 1);
    }

    private void UpdatePhysics()
    {
        playerVy += 0.2;
        playerY += playerVy;
        if (playerY >= GROUND)
        {
            playerY = GROUND;
            playerVy = 0;
            onGround = true;
        }
        if (playerY < 0)
        {
            playerY = 0;
            playerVy = 0;
        }
    }

    private void Update()
    {
        if (gameOver || paused) return;

        // Движение препятствий
        for (int i = obstacles.Count - 1; i >= 0; i--)
        {
            obstacles[i]--;
            if (obstacles[i] < 0) obstacles.RemoveAt(i);
        }

        // Столкновение
        int playerFloor = (int)playerY;
        foreach (int x in obstacles)
        {
            if (x == playerX && (playerFloor >= GROUND - 1 || playerFloor == GROUND - 1))
            {
                gameOver = true;
                if (score > record)
                {
                    record = score;
                    SaveRecord();
                }
                return;
            }
        }

        score++;
        speed = 1 + score / 10;

        if (rand.NextDouble() < 0.2 + Math.Min(0.4, score / 200.0))
            SpawnObstacle();

        UpdatePhysics();
    }

    private void Draw()
    {
        Console.Clear();
        Console.WriteLine(new string('═', WIDTH));
        Console.WriteLine($"  Счёт: {score}   Рекорд: {record}   Скорость: {speed}");
        Console.WriteLine(new string('═', WIDTH));

        char[,] grid = new char[HEIGHT, WIDTH];
        for (int y = 0; y < HEIGHT; y++)
            for (int x = 0; x < WIDTH; x++)
                grid[y, x] = ' ';
        foreach (int x in obstacles)
            if (x >= 0 && x < WIDTH) grid[GROUND, x] = OBSTACLE;
        int yPos = (int)playerY;
        if (yPos >= 0 && yPos < HEIGHT) grid[yPos, playerX] = PLAYER;

        for (int row = 0; row < HEIGHT; row++)
        {
            Console.Write('│');
            for (int col = 0; col < WIDTH; col++)
                Console.Write(grid[row, col]);
            Console.WriteLine('│');
        }
        Console.WriteLine(new string('═', WIDTH));
        string status = paused ? "ПАУЗА" : "ИГРА";
        Console.WriteLine($"  {status}  |  Пробел/↑ - прыжок  |  P - пауза  |  R - рестарт  |  Q - выход");
    }

    private void Reset()
    {
        playerY = GROUND;
        playerVy = 0;
        onGround = true;
        obstacles.Clear();
        score = 0;
        speed = 1;
        gameOver = false;
        paused = false;
    }

    public void Run()
    {
        while (running)
        {
            // Обработка ввода (неблокирующая)
            while (Console.KeyAvailable)
            {
                var key = Console.ReadKey(true);
                switch (key.Key)
                {
                    case ConsoleKey.Spacebar:
                    case ConsoleKey.UpArrow:
                        Jump();
                        break;
                    case ConsoleKey.P:
                        paused = !paused;
                        break;
                    case ConsoleKey.R:
                        Reset();
                        break;
                    case ConsoleKey.Q:
                        running = false;
                        break;
                }
            }

            // Обновление по таймеру
            if ((DateTime.Now - lastUpdate).TotalMilliseconds >= 1000.0 / 60)
            {
                Update();
                Draw();
                lastUpdate = DateTime.Now;
            }
            Thread.Sleep(10);
        }
    }

    static void Main()
    {
        var game = new RunnerGame();
        game.Run();
        Console.WriteLine("Игра завершена.");
    }
}
