package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Radius   float32
	Color    rl.Color
}

type Bar struct {
	Rec   rl.Rectangle
	Speed rl.Vector2
	Color rl.Color
}

type Player struct {
	Bar       Bar
	Points    int16
	LeftKey   int
	RightKey  int
	RoundWon  bool
	LastTouch bool
}

type Game struct {
	Round         int16
	Player1       Player
	Player2       Player
	Paused        bool
	RoundFinished bool
	FramesCounter int
}

func (b *Ball) genRandomStart() {
	b.Position.X = float32(rl.GetScreenWidth() / 2)
	b.Position.Y = float32(rl.GetRandomValue(int32(0+b.Radius), int32(rl.GetScreenHeight()-int(b.Radius))))
	newXSpeedDirection := rl.GetRandomValue(-1, 1)
	newYSpeedDirection := rl.GetRandomValue(-1, 1)
	// TODO: Idk how to do it better atm
	if newXSpeedDirection > 0 {
		b.Speed.X *= 1
	} else if newXSpeedDirection < 0 {
		b.Speed.X *= -1
	}

	if newYSpeedDirection > 0 {
		b.Speed.Y *= 1
	} else if newYSpeedDirection < 0 {
		b.Speed.Y *= -1
	}

	// b.Speed.X *= float32(rl.GetRandomValue(-1, 1))
	// b.Speed.Y *= float32(rl.GetRandomValue(-1, 1))
}

func main() {
	var wonText string
	rl.InitWindow(800, 450, "PonGo")
	rl.SetWindowMonitor(0)
	game := Game{
		Round: 1,
		Player1: Player{
			Bar: Bar{
				Rec: rl.Rectangle{
					X:      20,
					Y:      float32(rl.GetScreenHeight()) / 2,
					Width:  10,
					Height: 75,
				},
				Speed: rl.Vector2{
					X: 0,
					Y: 5.0,
				},
				Color: rl.RayWhite,
			},
			LeftKey:   rl.KeyA,
			RightKey:  rl.KeyD,
			Points:    0,
			RoundWon:  false,
			LastTouch: false,
		},
		Player2: Player{
			Bar: Bar{
				Rec: rl.Rectangle{
					X:      float32(rl.GetScreenWidth()) - 20,
					Y:      float32(rl.GetScreenHeight()) / 2,
					Width:  10,
					Height: 75,
				},
				Speed: rl.Vector2{
					X: 0,
					Y: 5.0,
				},
				Color: rl.RayWhite,
			},
			LeftKey:   rl.KeyLeft,
			RightKey:  rl.KeyRight,
			Points:    0,
			RoundWon:  false,
			LastTouch: false,
		},
		Paused:        false,
		RoundFinished: false,
		FramesCounter: 0,
	}
	ball := &Ball{
		Position: rl.Vector2{
			X: float32(rl.GetScreenWidth()) / 2.0,
			Y: float32(rl.GetScreenHeight()) / 2.0,
		},
		Speed: rl.Vector2{
			X: 6.0,
			Y: 4.0,
		},
		Radius: 10,
		Color:  rl.RayWhite,
	}
	ball.genRandomStart()

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyP) {
			game.Paused = !game.Paused
		}

		if rl.IsKeyPressed(rl.KeySpace) && game.RoundFinished {
			game.Round += 1
			game.RoundFinished = !game.RoundFinished
			ball.genRandomStart()
		}

		// Player 1 movement
		if !game.Paused && rl.IsKeyDown(int32(game.Player1.LeftKey)) {
			if game.Player1.Bar.Rec.Y >= 10 {
				game.Player1.Bar.Rec.Y -= float32(game.Player1.Bar.Speed.Y)
			}
		}

		if !game.Paused && rl.IsKeyDown(int32(game.Player1.RightKey)) {
			if game.Player1.Bar.Rec.Y <= float32(rl.GetScreenHeight())-game.Player1.Bar.Rec.Height-10 {
				game.Player1.Bar.Rec.Y += float32(game.Player1.Bar.Speed.Y)
			}
		}
		// End of player 1 movement

		// Player 2 movement
		if !game.Paused && rl.IsKeyDown(int32(game.Player2.LeftKey)) {
			if game.Player2.Bar.Rec.Y >= 10 {
				game.Player2.Bar.Rec.Y -= float32(game.Player2.Bar.Speed.Y)
			}
		}

		if !game.Paused && rl.IsKeyDown(int32(game.Player2.RightKey)) {
			if game.Player2.Bar.Rec.Y <= float32(rl.GetScreenHeight())-game.Player2.Bar.Rec.Height-10 {
				game.Player2.Bar.Rec.Y += float32(game.Player2.Bar.Speed.Y)
			}
		}
		// End of player 2 movement

		if !game.Paused {
			ball.Position.X += ball.Speed.X
			ball.Position.Y += ball.Speed.Y

			// Check Ball collision for Player 1 Bar
			// TODO: Player 1 and 2 should only touch the ball once (in part vecause teh ball get's kind og buggy sometimes and hits the bar multiple times) !game.Player1.LastTouch && ()
			// if !game.Player1.LastTouch && (ball.Position.X-float32(ball.Radius) < float32(game.Player1.Bar.Rec.X+game.Player1.Bar.Rec.Width+10) &&
			// 	ball.Position.X-float32(ball.Radius) > float32(game.Player1.Bar.Rec.X) && (ball.Position.Y-float32(ball.Radius) > game.Player1.Bar.Rec.Y-5 && ball.Position.Y-float32(ball.Radius) < float32(game.Player1.Bar.Rec.Y+game.Player1.Bar.Rec.Height+10))) {
			// 	game.Player1.LastTouch = true
			// 	game.Player2.LastTouch = false
			// 	ball.Speed.X *= -1.0
			// }

			if !game.Player1.LastTouch && rl.CheckCollisionCircleRec(ball.Position, float32(ball.Radius), game.Player1.Bar.Rec) {
				game.Player1.LastTouch = true
				game.Player2.LastTouch = false
				ball.Speed.X *= -1.0
			}

			// TODO: Add Player 2 collisions
			// Check Ball collision for Player 2 Bar
			// if !game.Player2.LastTouch && (ball.Position.X+float32(ball.Radius) < float32(game.Player2.Bar.Rec.X+game.Player2.Bar.Rec.Width+10) &&
			// 	ball.Position.X+float32(ball.Radius) > float32(game.Player2.Bar.Rec.X) && (ball.Position.Y-float32(ball.Radius) > game.Player2.Bar.Rec.Y-5 && ball.Position.Y-float32(ball.Radius) < float32(game.Player2.Bar.Rec.Y+game.Player2.Bar.Rec.Height+10))) {
			// 	game.Player1.LastTouch = false
			// 	game.Player2.LastTouch = true
			// 	ball.Speed.X *= -1.0
			// }
			if !game.Player2.LastTouch && rl.CheckCollisionCircleRec(ball.Position, float32(ball.Radius), game.Player2.Bar.Rec) {
				game.Player1.LastTouch = false
				game.Player2.LastTouch = true
				ball.Speed.X *= -1.0
			}

			if (ball.Position.Y >= (float32(rl.GetScreenHeight()) - ball.Radius)) || (ball.Position.Y <= float32(ball.Radius)) {
				ball.Speed.Y *= -1.0
			}
		} else {
			game.FramesCounter++
		}
		if !game.RoundFinished {
			if ball.Position.X < -float32(ball.Radius) {
				game.RoundFinished = !game.RoundFinished
				game.Player2.RoundWon = true
				game.Player2.Points += 1
				wonText = fmt.Sprintf("Player 2 won Round %03d", game.Round)
			}
			if ball.Position.X > (float32(rl.GetScreenWidth()) + ball.Radius) {
				game.RoundFinished = !game.RoundFinished
				game.Player1.RoundWon = true
				game.Player1.Points += 1
				wonText = fmt.Sprintf("Player 1 won Round %03d", game.Round)
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkGray)
		rl.DrawLine(int32(rl.GetScreenWidth()/2), 0, int32(rl.GetScreenWidth()/2), int32(rl.GetScreenHeight()), rl.RayWhite)

		// Draw Player 1 Bar
		rl.DrawRectangleRec(game.Player1.Bar.Rec, game.Player1.Bar.Color)

		// TODO: Add Player 2 Bar
		rl.DrawRectangleRec(game.Player2.Bar.Rec, game.Player2.Bar.Color)

		// Draw Ball
		rl.DrawCircleV(ball.Position, ball.Radius, ball.Color)

		// UI
		rl.DrawText("PRESS 'P' to PAUSE THE GAME", 10, int32(rl.GetScreenHeight())-25, 20, rl.LightGray)
		roundText := fmt.Sprintf("ROUND %03d", game.Round)
		rl.DrawTextEx(rl.GetFontDefault(), roundText, rl.Vector2{X: float32(rl.GetScreenWidth()/2 - int(rl.MeasureText(roundText, 20)/2)), Y: 10}, 20, 3, rl.Green)
		// Draw Points
		player1PointsText := fmt.Sprintf("Player 1: %d", game.Player1.Points)
		player2PointsText := fmt.Sprintf("Player 2: %d", game.Player2.Points)
		rl.DrawText(player1PointsText, int32(rl.GetScreenWidth()/2)-rl.MeasureText(roundText, 20)/2-rl.MeasureText(player1PointsText, 20)-20, 10, 20, rl.RayWhite)
		rl.DrawText(player2PointsText, int32(rl.GetScreenWidth()/2)+rl.MeasureText(roundText, 20)/2+20, 10, 20, rl.RayWhite)

		// Game Paused
		if game.Paused && ((game.FramesCounter/30)%2) != 0 {
			rl.DrawText("PAUSED", 350, 200, 30, rl.Gray)
		}

		// Round Won Text
		if game.RoundFinished {
			rl.DrawText(wonText, int32(rl.GetScreenWidth()/2)-rl.MeasureText(wonText, 30)/2, int32(rl.GetScreenHeight()/2-15), 30, rl.Green)
		}

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
