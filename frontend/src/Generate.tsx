import React, { useState } from 'react';
import axios from "axios";
import { SketchPicker } from 'react-color';

type RGBA = {
    r: number;
    g: number;
    b: number;
    a: number;
};

type Effects = {
    fog: boolean;
    glitch: boolean;
    sparkle: boolean;
};

const defaultEffects: Effects = {
    fog: false,
    glitch: false,
    sparkle: false,
};

function rgbaToCSS(rgba: RGBA): string {
    return `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})`;
}

function Generate() {

    const [text, setText] = useState<string>(`HELLO`);
    const [language, setLanguage] = useState<'japanese' | 'chinese'>('japanese');
    const [textColor, setTextColor] = useState<RGBA>({ r: 255, g: 255, b: 255, a: 1 });
    const [backgroundColor, setBackgroundColor] = useState<RGBA>({ r: 0, g: 0, b: 0, a: 1 });
    const [showPicker, setShowPicker] = useState<'text' | 'background' | null>(null);
    const [imageUrl, setImageUrl] = useState<string | null>(null);
    const [effects, setEffects] = useState<Effects>(defaultEffects);

    const handleEffectChange = (effect: keyof Effects) => {
        setEffects(prev => ({
            ...prev,
            [effect]: !prev[effect],
        }));
    };

    const generateImage = async () => {
        try {
            const response = await axios.post("http://localhost:8080/api/generate", {
                text,
                textColor,
                backgroundColor,
                language,
                effects: {
                    fog: effects.fog,
                    glitch: effects.glitch,
                    sparkle: effects.sparkle,
                },
            }, { responseType: "blob" });

            const imageBlob = new Blob([response.data], { type: "image/png" });
            setImageUrl(URL.createObjectURL(imageBlob));
        } catch (error) {
            console.error("エラー:", error);
            alert("画像生成に失敗しました");
        }
    };

    return (
        <div>
            {/* 言語切り替え */}
            <div>
                <button onClick={() => setLanguage('japanese')} disabled={language === 'japanese'}>日本語</button>
                <button onClick={() => setLanguage('chinese')} disabled={language === 'chinese'}>中国語</button>
            </div>

            {/* 入力 */}
            <form onSubmit={(e) => e.preventDefault()}>
                <input
                    type="text"
                    value={text}
                    onChange={(e) => setText(e.target.value)}
                />
                <input type="submit" />
            </form>

            {/* カラーピッカー切り替え */}
            <div>
                <button onClick={() => setShowPicker('text')}>文字色を選ぶ</button>
                <button onClick={() => setShowPicker('background')}>背景色を選ぶ</button>
            </div>

            {/* カラーピッカー本体 */}
            {showPicker === 'text' && (
                <SketchPicker
                    color={textColor}
                    onChange={(color) => {
                        const rgba = color.rgb;
                        setTextColor({
                            r: rgba.r,
                            g: rgba.g,
                            b: rgba.b,
                            a: rgba.a ?? 1  // aがundefinedなら1を使う
                        });
                    }}
                />
            )}
            {showPicker === 'background' && (
                <SketchPicker
                    color={backgroundColor}
                    onChange={(color) => {
                        const rgba = color.rgb;
                        setBackgroundColor({
                            r: rgba.r,
                            g: rgba.g,
                            b: rgba.b,
                            a: rgba.a ?? 1
                        });
                    }}

                />
            )}
            {/* エフェクトチェックボックス */}
            <fieldset>
                <legend>エフェクト</legend>
                <label>
                    <input
                        type="checkbox"
                        checked={effects.fog}
                        onChange={() => handleEffectChange('fog')}
                    />
                    fog
                </label>
                <label>
                    <input
                        type="checkbox"
                        checked={effects.glitch}
                        onChange={() => handleEffectChange('glitch')}
                    />
                    Glitch
                </label>
                <label>
                    <input
                        type="checkbox"
                        checked={effects.sparkle}
                        onChange={() => handleEffectChange('sparkle')}
                    />
                    Sparkle
                </label>
            </fieldset>

            {/* プレビュー */}
            <div style={{
                color: rgbaToCSS(textColor),
                backgroundColor: rgbaToCSS(backgroundColor),
                padding: '10px',
                display: 'inline-block'
            }}>
                {text}
            </div>

            <div><button onClick={generateImage}>画像を生成</button></div>

            <div>{imageUrl && <img src={imageUrl} alt="Generated" />}</div>
        </div>
    );
}

export default Generate;
