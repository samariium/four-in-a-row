package game

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yourname/fourinarow/internal/models"
	"github.com/yourname/fourinarow/internal/store"
	"github.com/yourname/fourinarow/internal/util"
)

type Manager struct {
	Store           *store.MongoStore
	MatchBotAfter   time.Duration
	RejoinGrace     time.Duration
	BotDelay        time.Duration

	upgrader websocket.Upgrader

	mu          sync.Mutex
	waiting     *waitingPlayer
	active      map[string]*state // gameId -> state
	userToGame  map[string]*userRef
}

type waitingPlayer struct {
	username string
	conn     *websocket.Conn
	timer    *time.Timer
}

type userRef struct {
	gameID string
	side   string // "R" or "Y"
}

type playerConn struct {
	username string
	conn     *websocket.Conn
	side     string
	bot      *Bot
}

type state struct {
	gameID   string
	p1       playerConn // R
	p2       playerConn // Y (or BOT)
	game     *GameLogic
	turn     string
	startAt  time.Time
	moves    []models.Move
	rejoinP1 *time.Timer
	rejoinP2 *time.Timer
}

func NewManager(store *store.MongoStore, matchBotMs, rejoinMs, botDelayMs int) *Manager {
	return &Manager{
		Store:         store,
		MatchBotAfter: time.Duration(matchBotMs) * time.Millisecond,
		RejoinGrace:   time.Duration(rejoinMs) * time.Millisecond,
		BotDelay:      time.Duration(botDelayMs) * time.Millisecond,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		active:     make(map[string]*state),
		userToGame: make(map[string]*userRef),
	}
}

func (m *Manager) HandleWS(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	gameID := r.URL.Query().Get("gameId")

	if username == "" {
		http.Error(w, "username required", http.StatusBadRequest)
		return
	}
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil { return }

	_ = m.Store.EnsurePlayer(r.Context(), username)

	if gameID != "" {
		m.tryRejoin(conn, username, gameID)
		return
	}
	m.enqueueOrMatch(conn, username)
}

func (m *Manager) enqueueOrMatch(conn *websocket.Conn, username string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ref, ok := m.userToGame[username]; ok {
		// already in a game; rejoin it
		m.mu.Unlock()
		m.tryRejoin(conn, username, ref.gameID)
		m.mu.Lock()
		return
	}

	if m.waiting != nil && m.waiting.username != username {
		wp := m.waiting
		if wp.timer != nil { wp.timer.Stop() }
		m.waiting = nil
		m.startGame(playerConn{username: wp.username, conn: wp.conn, side: "R"},
			playerConn{username: username, conn: conn, side: "Y"})
		return
	}

	// set waiting + bot fallback
	timer := time.AfterFunc(m.MatchBotAfter, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.waiting != nil && m.waiting.username == username {
			p1 := playerConn{username: username, conn: conn, side: "R"}
			p2 := playerConn{username: "BOT", conn: nil, side: "Y", bot: &Bot{Symbol: "Y"}}
			m.waiting = nil
			m.startGame(p1, p2)
		}
	})
	m.waiting = &waitingPlayer{username: username, conn: conn, timer: timer}
	sendJSON(conn, map[string]any{"type": "queued", "message": "Waiting for opponent..."})
}

func (m *Manager) startGame(p1, p2 playerConn) {
	st := &state{
		gameID:  util.NewID(10),
		p1:      p1,
		p2:      p2,
		game:    NewGame(),
		turn:    "R",
		startAt: time.Now(),
	}
	m.active[st.gameID] = st
	m.userToGame[p1.username] = &userRef{gameID: st.gameID, side: "R"}
	m.userToGame[p2.username] = &userRef{gameID: st.gameID, side: "Y"}

	startPayload := func(pc playerConn, opp string) map[string]any {
		return map[string]any{
			"type":     "start",
			"gameId":   st.gameID,
			"color":    pc.side,
			"opponent": opp,
			"board":    st.game.Board,
			"turn":     st.turn,
		}
	}
	sendJSON(p1.conn, startPayload(p1, p2.username))
	if p2.conn != nil {
		sendJSON(p2.conn, startPayload(p2, p1.username))
	}

	// readers
	go m.readLoop(st, p1)
	if p2.conn != nil {
		go m.readLoop(st, p2)
	}
}

func (m *Manager) tryRejoin(conn *websocket.Conn, username, gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	st, ok := m.active[gameID]
	if !ok {
		sendJSON(conn, map[string]any{"type": "error", "message": "game not found or finished"})
		m.mu.Unlock()
		m.enqueueOrMatch(conn, username)
		m.mu.Lock()
		return
	}
	isP1 := st.p1.username == username
	isP2 := st.p2.username == username
	if !isP1 && !isP2 {
		sendJSON(conn, map[string]any{"type": "error", "message": "this game does not belong to you"})
		m.mu.Unlock()
		m.enqueueOrMatch(conn, username)
		m.mu.Lock()
		return
	}

	if isP1 {
		st.p1.conn = conn
		if st.rejoinP1 != nil { st.rejoinP1.Stop(); st.rejoinP1 = nil }
	} else {
		st.p2.conn = conn
		if st.rejoinP2 != nil { st.rejoinP2.Stop(); st.rejoinP2 = nil }
	}
	sendJSON(conn, map[string]any{
		"type": "rejoined", "gameId": st.gameID,
		"color": func() string { if isP1 { return "R" } else { return "Y" } }(),
		"opponent": func() string { if isP1 { return st.p2.username } else { return st.p1.username } }(),
		"board": st.game.Board, "turn": st.turn,
	})
	go m.readLoop(st, func() playerConn {
		if isP1 { return st.p1 }
		return st.p2
	}())
}

func (m *Manager) readLoop(st *state, pc playerConn) {
	conn := pc.conn
	if conn == nil { return }
	defer func() {
		// disconnection -> start rejoin timer
		m.onDisconnect(st, pc.side)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { return }
		var in struct {
			Type string `json:"type"`
			Col  int    `json:"col"`
		}
		if err := json.Unmarshal(msg, &in); err != nil { continue }
		if in.Type == "move" {
			m.applyMove(st, pc.side, in.Col)
		}
	}
}

func (m *Manager) applyMove(st *state, side string, col int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if st.turn != side { return }

	row, ok := st.game.DropDisc(col, side)
	if !ok { return }

	st.moves = append(st.moves, models.Move{Player: side, Col: col, Row: row, At: time.Now()})

	nextTurn := map[string]string{"R": "Y", "Y": "R"}[side]
	update := map[string]any{"type": "update", "move": map[string]any{"row": row, "col": col, "player": side}, "board": st.game.Board, "turn": nextTurn}
	sendJSON(st.p1.conn, update)
	sendJSON(st.p2.conn, update)

	win := st.game.CheckWinner(side)
	full := st.game.IsFull()
	if win || full {
		go m.finishGame(st, func() string {
			if win {
				if side == "R" { return st.p1.username }
				return st.p2.username
			}
			return "Draw"
		}())
		return
	}

	st.turn = nextTurn

	// bot
	if st.p2.bot != nil && st.turn == "Y" {
		time.AfterFunc(m.BotDelay, func() {
			col := st.p2.bot.ChooseMove(st.game)
			m.applyMove(st, "Y", col)
		})
	}
}

func (m *Manager) onDisconnect(st *state, side string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	timer := time.AfterFunc(m.RejoinGrace, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		winner := st.p1.username
		if side == "R" { winner = st.p2.username }
		go m.finishGame(st, "Forfeit:"+winner)
	})
	if side == "R" {
		st.rejoinP1 = timer
		if st.p2.conn != nil {
			sendJSON(st.p2.conn, map[string]any{"type": "info", "message": "Opponent disconnected, waiting 30s to rejoin..."})
		}
	} else {
		st.rejoinP2 = timer
		if st.p1.conn != nil {
			sendJSON(st.p1.conn, map[string]any{"type": "info", "message": "Opponent disconnected, waiting 30s to rejoin..."})
		}
	}
}

func (m *Manager) finishGame(st *state, winnerLabel string) {
	// persist + leaderboard
	duration := int(time.Since(st.startAt).Seconds())
	isDraw := winnerLabel == "Draw"
	isForfeit := len(winnerLabel) > 8 && winnerLabel[:8] == "Forfeit:"
	winner := winnerLabel
	if isForfeit {
		winner = winnerLabel[8:]
	}

	_ = m.Store.InsertGame(context.Background(), models.GameDoc{
		GameID:     st.gameID,
		Player1:    st.p1.username,
		Player2:    st.p2.username,
		Winner:     winner,
		Duration:   duration,
		FinalBoard: st.game.Board,
		Moves:      st.moves,
	})

	if isDraw {
		_ = m.Store.IncDraws(context.Background(), []string{st.p1.username, st.p2.username})
	} else {
		loser := st.p1.username
		if winner == st.p1.username { loser = st.p2.username }
		_ = m.Store.IncWinLoss(context.Background(), winner, loser)
	}

	sendJSON(st.p1.conn, map[string]any{"type": "gameOver", "result": func() string {
		if isDraw { return "Draw" }
		return winner + " wins"
	}()})
	sendJSON(st.p2.conn, map[string]any{"type": "gameOver", "result": func() string {
		if isDraw { return "Draw" }
		return winner + " wins"
	}()})

	delete(m.active, st.gameID)
	delete(m.userToGame, st.p1.username)
	delete(m.userToGame, st.p2.username)
}

func sendJSON(conn *websocket.Conn, v any) {
	if conn == nil { return }
	_ = conn.WriteJSON(v)
}
