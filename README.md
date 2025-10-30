# four-in-a-row
# 🎮 Four-in-a-Row — Real-Time Multiplayer Web Game

A real-time **4-in-a-Row (Connect Four)** web game built with **Go (Golang)** for the backend and **React (Vite)** for the frontend.  
You can **play against your friends or challenge a smart bot**, with a persistent **Leaderboard** powered by MongoDB Atlas.

---

## 🚀 Live Demo

- 🎨 **Frontend (Vercel):** [https://four-in-a-bdgrhaofn-samars-projects-920aa19d.vercel.app](https://four-in-a-bdgrhaofn-samars-projects-920aa19d.vercel.app)
- ⚙️ **Backend (Render):** [https://fourinarow-backend.onrender.com](https://fourinarow-backend.onrender.com)
- 🏆 **Leaderboard API:** [https://fourinarow-backend.onrender.com/leaderboard](https://fourinarow-backend.onrender.com/leaderboard)

> 💤 **Note:**  
> The backend runs on **Render’s free tier**, so it may take **20–40 seconds to wake up** after inactivity.  
> To “wake” it, open the **Leaderboard API link** first before starting a match.

---

## 🧱 Tech Stack

### **Frontend**
- React + Vite ⚛️  
- Axios for API calls  
- WebSocket for real-time play  
- TailwindCSS for UI  

### **Backend (GoLang)**
- Go (net/http + gorilla/websocket) 🦫  
- MongoDB Atlas 🍃  
- Render (hosting) ☁️  
- REST + WebSocket architecture  

---

## 🕹️ Features

- 🧠 **Play vs Bot** or 🧑‍🤝‍🧑 **Play vs Friend**  
- ⚡ Real-time gameplay using WebSockets  
- 💾 Persistent leaderboard (MongoDB Atlas)  
- 🔁 Auto-rejoin on disconnect  
- 🌈 Interactive 7×6 game board with hover effects  
- ☁️ Fully deployed (Render + Vercel)  

---

## 📂 Project Structure
Four-in-a-row/
├── frontend/ # React + Vite app
│ ├── src/
│ ├── public/
│ ├── package.json
│ └── vite.config.js
│
├── go-backend/ # Go backend
│ ├── cmd/server/main.go
│ ├── internal/
│ │ ├── game/
│ │ ├── models/
│ │ ├── store/
│ │ ├── config/
│ │ └── util/
│ ├── go.mod
│ └── go.sum
│
└── README.md

---

## ⚙️ Local Setup Instructions

### 🖥️ Clone the Repository

git clone https://github.com/samariium/four-in-a-row.git
cd four-in-a-row

Backend Setup (GoLang)

Navigate to backend:

cd go-backend


Create a .env file:

MONGO_URI=mongodb+srv://<username>:<password>@cluster0.bhydljn.mongodb.net/?retryWrites=true&w=majority
PORT=9090


Run the backend:

go run ./cmd/server


Check it:

http://localhost:9090/health
