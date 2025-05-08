import React, { useEffect, useRef, useState } from "react";
import axios from "axios";

type StampData = {
    text: string;
    textColor: string;
    backgroundColor: string;
    selectedEffect: "none" | "glitch" | "jitter" | "rotate" | "shadow" | "blur";
};

export default function RecentStamps() {
    const [stamps, setStamps] = useState<StampData[]>([]);

    useEffect(() => {
        axios.get("http://localhost:8080/api/stamp/get").then((res) => {
            setStamps(res.data);
        });
    }, []);

    return (
        <div
            style={{
                overflowX: "auto",
                whiteSpace: "nowrap",
                padding: "1rem",
                border: "1px solid #ccc",
            }}
        >
            {stamps.map((stamp, index) => (
                <StampCanvas key={index} {...stamp} />
            ))}
        </div>
    );
}

function StampCanvas({
    text,
    textColor,
    backgroundColor,
    selectedEffect,
}: StampData) {
    const canvasRef = useRef<HTMLCanvasElement | null>(null);

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;

        const width = 128;
        const height = 128;
        canvas.width = width;
        canvas.height = height;

        const centerX = width / 2;
        const centerY = height / 2;

        // 背景描画
        ctx.fillStyle = backgroundColor;
        ctx.fillRect(0, 0, width, height);

        ctx.font = "bold 20px sans-serif";
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        ctx.fillStyle = textColor;

        // エフェクト処理
        switch (selectedEffect) {
            case "shadow":
                ctx.shadowColor = "rgba(0,0,0,0.5)";
                ctx.shadowOffsetX = 2;
                ctx.shadowOffsetY = 2;
                ctx.shadowBlur = 4;
                break;
            case "jitter":
                for (let i = 0; i < 5; i++) {
                    const offsetX = Math.random() * 4 - 2;
                    const offsetY = Math.random() * 4 - 2;
                    ctx.fillText(text, width / 2 + offsetX, height / 2 + offsetY);
                }
                return;
            case "glitch":
                for (let i = 0; i < 5; i++) {
                    ctx.fillStyle = i % 2 === 0 ? 'red' : 'cyan';
                    ctx.fillText(text, centerX + Math.random() * 10 - 5, centerY + Math.random() * 10 - 5);
                }
                break;
            case "rotate":
                ctx.save();
                ctx.translate(width / 2, height / 2);
                ctx.rotate(0.1);
                ctx.fillText(text, 0, 0);
                ctx.restore();
                return;
            case 'shadow':
                ctx.shadowColor = 'rgba(0, 0, 0, 0.5)';
                ctx.shadowOffsetX = 4;
                ctx.shadowOffsetY = 4;
                ctx.shadowBlur = 10;
                ctx.fillText(text, centerX, centerY);
                break;
            case 'blur':
                ctx.filter = 'blur(2px)';
                ctx.fillText(text, centerX, centerY);
                ctx.filter = 'none';
                break;
            default:
                break;
        }

        // 通常描画 or blur のベース
        ctx.fillText(text, width / 2, height / 2);

        if (selectedEffect === "blur") {
            canvas.style.filter = "blur(2px)";
        } else {
            canvas.style.filter = "none";
        }
    }, [text, textColor, backgroundColor, selectedEffect]);

    return (
        <canvas
            ref={canvasRef}
            style={{
                display: "inline-block",
                marginRight: "1rem",
                border: "1px solid #ddd",
                borderRadius: "8px",
                background: "#fff",
            }}
        />
    );
}
