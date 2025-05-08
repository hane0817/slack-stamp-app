import { useRef, useEffect } from 'react';
import axios from 'axios';

const EFFECTS = ['none', 'glitch', 'jitter', 'rotate', 'shadow', 'blur'] as const;
type EffectType = typeof EFFECTS[number];




export function TextEffectCanvas(
    text: string,
    textColor: string,
    backgroundColor: string,
    language: 'japanese' | 'chinese',
    selectedEffect: EffectType) {
    const canvasRef = useRef<HTMLCanvasElement | null>(null);

    useEffect(() => {
        const canvas = canvasRef.current;
        const ctx = canvas?.getContext('2d');
        if (!canvas || !ctx) return;

        const width = canvas.width;
        const height = canvas.height;

        ctx.clearRect(0, 0, width, height);
        ctx.fillStyle = backgroundColor;
        ctx.fillRect(0, 0, width, height);

        // リセット
        ctx.shadowColor = 'transparent';
        ctx.shadowBlur = 0;
        ctx.shadowOffsetX = 0;
        ctx.shadowOffsetY = 0;
        ctx.filter = 'none';

        ctx.fillStyle = textColor;
        ctx.font = '48px sans-serif';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';

        const centerX = width / 2;
        const centerY = height / 2;

        switch (selectedEffect) {
            case 'glitch':
                for (let i = 0; i < 5; i++) {
                    ctx.fillStyle = i % 2 === 0 ? 'red' : 'cyan';
                    ctx.fillText(text, centerX + Math.random() * 10 - 5, centerY + Math.random() * 10 - 5);
                }
                break;
            case 'jitter':
                ctx.fillText(text, centerX + Math.random() * 6 - 3, centerY + Math.random() * 6 - 3);
                break;
            case 'rotate':
                ctx.save();
                ctx.translate(centerX, centerY);
                ctx.rotate(Math.random() * 0.2 - 0.1);
                ctx.fillText(text, 0, 0);
                ctx.restore();
                break;
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
                ctx.fillText(text, centerX, centerY);
        }
    }, [text, textColor, backgroundColor, selectedEffect]);

    const sendStamp = async () => {
        try {
            await axios.post('http://localhost:8080/api/stamp/post', {
                text: text,
                textColor: textColor,
                backgroundColor: backgroundColor,
                language: language,
                selectedEffect: selectedEffect,
            });
            alert('スタンプを送信しました');
        } catch (err) {
            console.error('送信失敗:', err);
            alert('送信に失敗しました');
        }
    };

    return { canvasRef, sendStamp };
};

export default TextEffectCanvas;
