import { useState } from "react";

export default function UsernameForm({ onSubmit }) {
  const [name, setName] = useState("");
  return (
    <div style={{ textAlign: "center" }}>
      <h1>🎯 4 in a Row</h1>
      <input
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Enter your username"
        style={{ padding: "10px", borderRadius: "8px", border: "none", marginRight: "8px" }}
      />
      <button onClick={() => onSubmit(name.trim())}>Start</button>
    </div>
  );
}
