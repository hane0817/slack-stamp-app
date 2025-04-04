import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Login from "./Login";
import Generate from "./Generate";
import Register from "./Register";

function App() {
  const [token, setToken] = useState<string | null>(localStorage.getItem("token"));

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    if (storedToken) {
      setToken(storedToken);
    }
  }, [token]);

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login setToken={setToken} />} />
        <Route path="/register" element={<Register />} />
        <Route path="/generate" element={token ? <Generate /> : <Navigate to="/login" />} />
        <Route path="*" element={<Navigate to={token ? "/generate" : "/login"} />} />
      </Routes>
    </Router>
  );
}

export default App;
