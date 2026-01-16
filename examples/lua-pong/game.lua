-- Lua Pong Example Game

local game = {}

-- Game state
local ball = { x = 0, y = 0, vx = 1, vy = 0.5 }
local paddle = { y = 0 }
local score = 0
local width, height = 40, 20

function game.init(ctx)
    width = ctx.width
    height = ctx.height
    ball.x = width / 2
    ball.y = height / 2
    paddle.y = height / 2
    score = 0
end

function game.start()
    -- Nothing to do
end

function game.stop()
    -- Nothing to do
end

function game.update(dt)
    -- Move ball
    ball.x = ball.x + ball.vx
    ball.y = ball.y + ball.vy

    -- Bounce off top/bottom
    if ball.y <= 1 or ball.y >= height - 2 then
        ball.vy = -ball.vy
    end

    -- Bounce off right wall
    if ball.x >= width - 2 then
        ball.vx = -ball.vx
    end

    -- Check paddle collision
    if ball.x <= 3 and ball.y >= paddle.y - 2 and ball.y <= paddle.y + 2 then
        ball.vx = -ball.vx
        score = score + 1
    end

    -- Ball out of bounds (left)
    if ball.x < 1 then
        -- Reset
        ball.x = width / 2
        ball.y = height / 2
        ball.vx = 1
        score = 0
    end
end

function game.render()
    slip.clear()

    -- Draw border
    slip.box(0, 0, slip.width(), slip.height())

    -- Draw score
    slip.color("cyan")
    slip.text(slip.width() / 2 - 5, 0, "Score: " .. score)

    -- Draw paddle
    slip.color("green")
    for i = -2, 2 do
        slip.set(2, math.floor(paddle.y) + i, "|")
    end

    -- Draw ball
    slip.color("yellow")
    slip.set(math.floor(ball.x), math.floor(ball.y), "o")
end

function game.input(key)
    if key == "up" and paddle.y > 3 then
        paddle.y = paddle.y - 1
    elseif key == "down" and paddle.y < height - 4 then
        paddle.y = paddle.y + 1
    end
end

return game
