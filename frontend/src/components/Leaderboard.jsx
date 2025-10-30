import { useEffect, useState } from "react";
import axios from "axios";

export default function Leaderboard() {
  const [players, setPlayers] = useState([]);

  // ✅ Use environment variable for backend URL
  const backendUrl = import.meta.env.VITE_BACKEND_URL || "http://localhost:9090";

  useEffect(() => {
    axios
      .get(`${backendUrl}/leaderboard`)
      .then((res) => setPlayers(res.data))
      .catch((err) => {
        console.error("Failed to load leaderboard:", err);
        setPlayers([]);
      });
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
