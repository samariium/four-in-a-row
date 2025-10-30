import { useEffect, useState } from "react";
import axios from "axios";

export default function Leaderboard() {
  const [players, setPlayers] = useState([]);

  useEffect(() => {
    axios
      .get("http://localhost:9090/leaderboard")
      .then((res) => setPlayers(res.data))
      .catch(() => setPlayers([]));
  }, []);

  return (
    <div style={{ marginTop: "2rem", textAlign: "center" }}>
      <h2>🏆 Leaderboard</h2>
      <table style={{ margin: "0 auto", color: "white" }}>
        <thead>
          <tr>
            <th>User</th>
            <th>Wins</th>
            <th>Draws</th>
            <th>Losses</th>
          </tr>
        </thead>
        <tbody>
          {players.map((p, i) => (
            <tr key={i}>
              <td>{p.username}</td>
              <td>{p.wins}</td>
              <td>{p.draws}</td>
              <td>{p.losses}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
