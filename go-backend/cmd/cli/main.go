package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type StartMsg struct {
	Type     string        `json:"type"`
	GameID   string        `json:"gameId"`
	Color    string        `json:"color"`
	Opponent string        `json:"opponent"`
	Board    [][]*string   `json:"board"`
	Turn     string        `json:"turn"`
}
type UpdateMsg struct {
	Type  string        `json:"type"`
	Move  struct {
		Row    int    `json:"row"`
		Col    int    `json:"col"`
		Player string `json:"player"`
	} `json:"move"`
	Board [][]*string `json:"board"`
	Turn  string      `json:"turn"`
}
type SimpleMsg struct {
	Type   string `json:"type"`
	Result string `json:"result,omitempty"`
	Msg    string `json:"message,omitempty"`
	Color  string `json:"color,omitempty"`
	Turn   string `json:"turn,omitempty"`
}

func printBoard(board [][]*string) {
	fmt.Println()
	for r := 0; r < len(board); r++ {
		fmt.Print("|")
		for c := 0; c < len(board[r]); c++ {
			cell := "."
			if board[r][c] != nil {
				cell = *board[r][c]
			}
			fmt.Printf(" %s ", cell)
		}
		fmt.Println("|")
	}
	fmt.Println("  0  1  2  3  4  5  6")
	fmt.Println()
}

func firstPlayableCol(board [][]*string) int {
	if len(board) == 0 {
		return 3
	}
	// prefer center outward
	order := []int{3, 2, 4, 1, 5, 0, 6}
	for _, c := range order {
		if board[0][c] == nil {
			return c
		}
	}
	// fallback
	for c := 0; c < len(board[0]); c++ {
		if board[0][c] == nil {
			return c
		}
	}
	return 3
}

func main() {
	server := flag.String("server", "ws://localhost:9090/ws", "WebSocket server URL")
	user := flag.String("user", "", "Username")
	auto := flag.Bool("auto", false, "Auto-play moves when it's your turn")
	flag.Parse()

	if strings.TrimSpace(*user) == "" {
		log.Fatal("provide -user <name>")
	}

	url := fmt.Sprintf("%s?username=%s", *server, *user)
	log.Printf("Connecting to %s ...", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	var myColor = ""
	var board [][]*string
	var nextTurn = "" // who moves next, "R" or "Y"

	// input reader for manual moves
	reader := bufio.NewReader(os.Stdin)

	// sender helper
	sendMove := func(col int) {
		payload := map[string]any{"type": "move", "col": col}
		_ = conn.WriteJSON(payload)
	}

	// prompt loop (manual)
	promptIfMyTurn := func() {
		if myColor != "" && nextTurn == myColor && !*auto {
			fmt.Print("Your move (enter column 0-6): ")
		}
	}

	go func() {
		// manual input goroutine
		for {
			if !*auto {
				line, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				line = strings.TrimSpace(line)
				if line == "" {
					promptIfMyTurn()
					continue
				}
				// only accept input when it's my turn
				if nextTurn != myColor {
					fmt.Println("Not your turn yet.")
					continue
				}
				var c int
				_, err = fmt.Sscanf(line, "%d", &c)
				if err != nil || c < 0 || c > 6 {
					fmt.Println("Enter a valid column (0-6).")
					promptIfMyTurn()
					continue
				}
				sendMove(c)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	// receiver loop
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		// sniff type
		var peek struct{ Type string `json:"type"` }
		_ = json.Unmarshal(data, &peek)

		switch peek.Type {
		case "queued":
			var m SimpleMsg
			_ = json.Unmarshal(data, &m)
			fmt.Println("⏳", m.Msg)

		case "start":
			var m StartMsg
			_ = json.Unmarshal(data, &m)
			myColor = m.Color
			board = m.Board
			nextTurn = m.Turn
			fmt.Printf("🎮 Game started! You are %s vs %s. Next turn: %s\n", myColor, m.Opponent, nextTurn)
			printBoard(board)
			if *auto && nextTurn == myColor {
				col := firstPlayableCol(board)
				fmt.Printf("🤖 Auto move -> %d\n", col)
				sendMove(col)
			}
			promptIfMyTurn()

		case "rejoined":
			var m StartMsg // same fields used
			_ = json.Unmarshal(data, &m)
			myColor = m.Color
			board = m.Board
			nextTurn = m.Turn
			fmt.Printf("🔁 Rejoined. You are %s vs %s. Next turn: %s\n", myColor, m.Opponent, nextTurn)
			printBoard(board)
			if *auto && nextTurn == myColor {
				col := firstPlayableCol(board)
				fmt.Printf("🤖 Auto move -> %d\n", col)
				sendMove(col)
			}
			promptIfMyTurn()

		case "update":
			var m UpdateMsg
			_ = json.Unmarshal(data, &m)
			board = m.Board
			nextTurn = m.Turn
			fmt.Printf("⬇️  %s played col %d (row %d). Next: %s\n", m.Move.Player, m.Move.Col, m.Move.Row, nextTurn)
			printBoard(board)
			if *auto && nextTurn == myColor {
				col := firstPlayableCol(board)
				time.Sleep(300 * time.Millisecond)
				fmt.Printf("🤖 Auto move -> %d\n", col)
				sendMove(col)
			}
			promptIfMyTurn()

		case "info":
			var m SimpleMsg
			_ = json.Unmarshal(data, &m)
			fmt.Println("ℹ️ ", m.Msg)

		case "gameOver":
			var m SimpleMsg
			_ = json.Unmarshal(data, &m)
			fmt.Println("🏁", m.Result)
			printBoard(board)
			return

		case "error":
			var m SimpleMsg
			_ = json.Unmarshal(data, &m)
			fmt.Println("❌", m.Msg)
			return

		default:
			fmt.Println("… unknown message:", string(data))
		}
	}
}

