import { useNavigate } from "react-router-dom";

interface HeaderProps {
    setToken: (token: string | null) => void;
}

const Header: React.FC<HeaderProps> = ({ setToken }) => {
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem("token"); // JWTを削除
        setToken(null); // 状態をクリア
        navigate("/login"); // ログイン画面へリダイレクト
    };

    return (
        <header>
            <h1>Slack-stamp-app</h1>
            <button onClick={handleLogout}>Logout</button> {/* ログアウトボタン */}
        </header>
    );
};

export default Header;
