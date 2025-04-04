import { useState } from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { LoginResponse } from "./types";

interface LoginProps {
    setToken: (token: string | null) => void;
}

const Login: React.FC<LoginProps> = ({ setToken }) => {
    const [name, setname] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError("");

        try {
            const response = await axios.post<LoginResponse>("http://localhost:8080/auth/login", { name, password },
                {
                    headers: {
                        "Content-Type": "application/json",
                    },
                }
            );
            const token = response.data.token;
            localStorage.setItem("token", token);
            setToken(token);
            console.log("Navigating to /Generate after login success...");
            navigate("/Generate");
        } catch (err) {
            setError("Login failed. Check your credentials.");
        }
    };

    return (
        <div>
            <h2>Login</h2>
            {error && <p style={{ color: "red" }}>{error}</p>}
            <form onSubmit={handleSubmit}>
                <input type="text" placeholder="name" value={name} onChange={(e) => setname(e.target.value)} />
                <input type="password" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                <button type="submit">Login</button>
            </form>
            <button onClick={() => navigate("/register")}>Register</button>
        </div>
    );
};

export default Login;
