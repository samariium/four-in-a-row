import { useState } from "react";
import UsernameForm from "./components/UsernameForm";
import GameBoard from "./components/GameBoard";
import Leaderboard from "./components/Leaderboard";
import { useGameSocket } from "./hooks/useGameSocket";
import "./styles/index.css";

export default function App() {
  const [username, setUsername] = useState("");
  const { status, gameState, sendMove, opponent } = useGameSocket(username);

  if (!username) return <UsernameForm onSubmit={setUsername} />;

  return (
    <div style={{ textAlign: "center" }}>
      <h1>4 in a Row</h1>
      <p>Status: {status}</p>
      {opponent && <p>Opponent: {opponent}</p>}
      {gameState.board && (
        <GameBoard
          board={gameState.board}
          color={gameState.color}
          turn={gameState.turn}
          onMove={sendMove}
        />
      )}
      <Leaderboard />
    </div>
  );
}
