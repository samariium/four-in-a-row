export default function GameBoard({ board, color, turn, onMove }) {
  if (!board || !board.length) return null;

  return (
    <div style={{ textAlign: "center" }}>
      <h2>
        Your color: <span style={{ color: color === "R" ? "#ef4444" : "#facc15" }}>{color}</span>
      </h2>
      <h3>Turn: {turn}</h3>

      <div className="board">
        {board.map((row, rIndex) =>
          row.map((cell, cIndex) => (
            <div
              key={`${rIndex}-${cIndex}`}
              className={`cell ${cell === "R" ? "red" : cell === "Y" ? "yellow" : ""}`}
              onClick={() => onMove(cIndex)}
            />
          ))
        )}
      </div>
    </div>
  );
}
