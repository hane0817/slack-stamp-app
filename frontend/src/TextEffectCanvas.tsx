import React, { useRef, useEffect } from 'react';

type Props = {
    text: string;
    fontSize?: number;
    textColor?: string;
    backgroundColor?: string;
    selectedEffect?: 'none' | 'glitch' | 'jitter' | 'rotate' | 'shadow' | 'blur';
};

const TextEffectCanvas: React.FC<Props> = ({
    text,
    fontSize = 48,
    textColor = '#FFFFFF',
    backgroundColor = '#000000',
    selectedEffect = 'none',
}) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);

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

    return <canvas ref={canvasRef} />;
};

export default TextEffectCanvas;
