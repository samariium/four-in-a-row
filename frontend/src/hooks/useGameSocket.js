import { useEffect, useState, useRef } from "react";

export function useGameSocket(username) {
  const backendUrl = import.meta.env.VITE_BACKEND_URL || "http://localhost:9090";
  const WS_URL = backendUrl.replace("http", "ws") + "/ws";
  const [socket, setSocket] = useState(null);
  const [gameState, setGameState] = useState({ board: [], turn: null, color: null });
  const [status, setStatus] = useState("connecting");
  const [opponent, setOpponent] = useState(null);
  const gameIdRef = useRef(null);

  useEffect(() => {
    if (!username) return;
    const ws = new WebSocket(`${WS_URL}?username=${username}`);
    setSocket(ws);

    ws.onopen = () => setStatus("waiting");
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      switch (data.type) {
        case "queued":
          setStatus("waiting");
          break;
        case "start":
          setStatus("playing");
          setOpponent(data.opponent);
          gameIdRef.current = data.gameId;
          setGameState({ board: data.board, color: data.color, turn: data.turn });
          break;
        case "update":
          setGameState((s) => ({ ...s, board: data.board, turn: data.turn }));
          break;
        case "gameOver":
          setStatus("ended");
          alert(data.result);
          break;
        case "info":
          alert(data.message);
          break;
        case "rejoined":
          setStatus("rejoined");
          setOpponent(data.opponent);
          setGameState({ board: data.board, color: data.color, turn: data.turn });
          break;
        default:
          break;
      }
    };
    ws.onclose = () => setStatus("disconnected");
    return () => ws.close();
  }, [username]);

  const sendMove = (col) => {
    if (socket && status === "playing") {
      socket.send(JSON.stringify({ type: "move", col }));
    }
  };

  return { status, gameState, sendMove, opponent };
}
