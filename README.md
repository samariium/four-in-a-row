
# Four-in-a-Row Real-Time Multiplayer Web Game

A real-time **4-in-a-Row (Connect Four)** web game built with **Go (Golang)** for the backend and **React (Vite)** for the frontend.  
You can **play against your friends or challenge a smart bot**, with a persistent **Leaderboard** powered by MongoDB Atlas.

---

## Live Demo

- **Frontend (Vercel):** [https://four-in-a-bdgrhaofn-samars-projects-920aa19d.vercel.app](https://four-in-a-bdgrhaofn-samars-projects-920aa19d.vercel.app)
- **Backend (Render):** [https://fourinarow-backend.onrender.com](https://fourinarow-backend.onrender.com)
- **Leaderboard API:** [https://fourinarow-backend.onrender.com/leaderboard](https://fourinarow-backend.onrender.com/leaderboard)

> **Note:**  
> The backend runs on **Render’s free tier**, so it may take **20–40 seconds to wake up** after inactivity.  
> To “wake” it, open the **Leaderboard API link** first before starting a match.

---

## Tech Stack

### **Frontend**
- React + Vite     
- WebSocket for real-time play  
    

### **Backend (GoLang)**
- Go 
- MongoDB Atlas 
- Render (hosting)     

---

## Features

- **Play vs Bot** or **Play vs Friend**  
- Real-time gameplay using WebSockets  
- Persistent leaderboard (MongoDB Atlas)  
- Auto-rejoin on disconnect  
- Interactive 7×6 game board with hover effects  
- Fully deployed (Render + Vercel)  

---

## Project Structure
Four-in-a-row/<br>
├── frontend/ # React + Vite app<br>
│ ├── src/<br>
│ ├── public/<br>
│ ├── package.json<br>
│ └── vite.config.js<br>
│<br>
├── go-backend/ # Go backend<br>
│ ├── cmd/server/main.go<br>
│ ├── internal/<br>
│ │ ├── game/<br>
│ │ ├── models/<br>
│ │ ├── store/<br>
│ │ ├── config/<br>
│ │ └── util/<br>
│ ├── go.mod<br>
│ └── go.sum<br>
│<br>
└── README.md<br>

---



---

##  Local Setup Instructions

### Clone the Repository
```bash
git clone https://github.com/samariium/four-in-a-row.git
cd four-in-a-row
```
Backend Setup (GoLang)
Navigate to Backend
```bash
cd go-backend
```

Create a .env file
```bash
MONGO_URI=mongodb+srv://<username>:<password>@cluster0.bhydljn.mongodb.net/?retryWrites=true&w=majority
PORT=9090
```

Run the Backend
```bash
go run ./cmd/server
```

Test Health Endpoint
```bash
Visit → http://localhost:9090/health
```

Expected Output:
```bash
{"ok": true}
```

Frontend Setup (React + Vite)
Navigate to Frontend
```bash
cd frontend
```

Create a .env file
```bash
VITE_BACKEND_URL=http://localhost:9090
```

Install Dependencies
```bash
npm install
```

Run Locally
```bash
npm run dev
```

Visit in Browser
```bash
http://localhost:5173
```
