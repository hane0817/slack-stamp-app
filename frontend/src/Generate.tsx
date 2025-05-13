import React, { useState } from 'react';
import axios from "axios";
import { SketchPicker } from 'react-color';
import TextEffectCanvas from './TextEffectCanvas';
import RecentStamps from "./RecentStamps";
import './Generate.css';

type RGBA = {
    r: number;
    g: number;
    b: number;
    a: number;
};

const EFFECTS = ['none', 'glitch', 'jitter', 'rotate', 'shadow', 'blur'] as const;
type EffectType = typeof EFFECTS[number];

function rgbaToCSS(rgba: RGBA): string {
    return `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${rgba.a})`;
}

function Generate() {
    const [text, setText] = useState('HELLO');
    const [language, setLanguage] = useState<'japanese' | 'chinese'>('japanese');
    const [textColor, setTextColor] = useState<RGBA>({ r: 255, g: 255, b: 255, a: 1 });
    const [backgroundColor, setBackgroundColor] = useState<RGBA>({ r: 0, g: 0, b: 0, a: 1 });
    const [showPicker, setShowPicker] = useState<'text' | 'background' | null>(null);
    const [imageUrl, setImageUrl] = useState<string | null>(null);
    const [selectedEffect, setSelectedEffect] = useState<EffectType>('none');
    const [showEffectOptions, setShowEffectOptions] = useState(false);

    const generateImage = async () => {
        try {
            const response = await axios.post("http://localhost:8080/api/generate", {
                text,
                textColor,
                backgroundColor,
                language
            }, { responseType: "blob" });

            const imageBlob = new Blob([response.data], { type: "image/png" });
            setImageUrl(URL.createObjectURL(imageBlob));
        } catch (error) {
            console.error("エラー:", error);
            alert("画像生成に失敗しました");
        }
    };

    const { canvasRef, sendStamp } = TextEffectCanvas(
        text,
        rgbaToCSS(textColor),
        rgbaToCSS(backgroundColor),
        language,
        selectedEffect,
    );

    return (
        <div className="min-h-screen bg-gray-100 py-8 px-4">
            <div className="max-w-xl mx-auto bg-white p-6 rounded-lg shadow-md space-y-6">

                {/* 言語切り替え */}
                <div className="flex space-x-2">
                    <button
                        onClick={() => setLanguage('japanese')}
                        disabled={language === 'japanese'}
                        className={`px-4 py-2 rounded ${language === 'japanese' ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
                    >
                        日本語
                    </button>
                    <button
                        onClick={() => setLanguage('chinese')}
                        disabled={language === 'chinese'}
                        className={`px-4 py-2 rounded ${language === 'chinese' ? 'bg-blue-500 text-white' : 'bg-gray-200'}`}
                    >
                        中国語
                    </button>
                </div>

                {/* 入力 */}
                <form onSubmit={(e) => e.preventDefault()} className="flex space-x-2">
                    <input
                        type="text"
                        value={text}
                        onChange={(e) => setText(e.target.value)}
                        className="flex-grow border border-gray-300 px-4 py-2 rounded"
                    />
                    <input
                        type="submit"
                        value="決定"
                        className="bg-blue-500 text-white px-4 py-2 rounded"
                    />
                </form>

                {/* カラーピッカー切り替え */}
                <div className="flex space-x-2">
                    <button onClick={() => setShowPicker('text')} className="bg-indigo-500 text-white px-4 py-2 rounded">
                        文字色を選ぶ
                    </button>
                    <button onClick={() => setShowPicker('background')} className="bg-indigo-500 text-white px-4 py-2 rounded">
                        背景色を選ぶ
                    </button>
                </div>

                {/* カラーピッカー本体 */}
                <div className="space-y-4">
                    {showPicker === 'text' && (
                        <SketchPicker
                            color={textColor}
                            onChange={(color) => {
                                const rgba = color.rgb;
                                setTextColor({ r: rgba.r, g: rgba.g, b: rgba.b, a: rgba.a ?? 1 });
                            }}
                        />
                    )}
                    {showPicker === 'background' && (
                        <SketchPicker
                            color={backgroundColor}
                            onChange={(color) => {
                                const rgba = color.rgb;
                                setBackgroundColor({ r: rgba.r, g: rgba.g, b: rgba.b, a: rgba.a ?? 1 });
                            }}
                        />
                    )}
                </div>

                {/* プレビュー */}
                <div
                    className="text-center font-bold text-2xl py-4 rounded"
                    style={{
                        color: rgbaToCSS(textColor),
                        backgroundColor: rgbaToCSS(backgroundColor),
                    }}
                >
                    {text}
                </div>

                <div>
                    <button
                        onClick={generateImage}
                        className="w-full bg-green-500 hover:bg-green-600 text-white py-2 rounded"
                    >
                        画像を生成
                    </button>
                </div>

                {/* 画像表示 */}
                {imageUrl && (
                    <div className="text-center">
                        <img src={imageUrl} alt="Generated" className="mx-auto mt-4 rounded" />
                    </div>
                )}

                {/* エフェクト */}
                <div className="space-y-2">
                    <button
                        onClick={() => setShowEffectOptions(!showEffectOptions)}
                        className="bg-purple-500 text-white px-4 py-2 rounded"
                    >
                        エフェクトを追加
                    </button>

                    {showEffectOptions && (
                        <div className="flex flex-wrap gap-2">
                            {EFFECTS.map((effect) => (
                                <button
                                    key={effect}
                                    onClick={() => setSelectedEffect(effect)}
                                    className={`px-3 py-1 rounded border ${selectedEffect === effect
                                        ? 'bg-purple-600 text-white'
                                        : 'bg-gray-100'
                                        }`}
                                >
                                    {effect}
                                </button>
                            ))}
                        </div>
                    )}
                </div>

                {/* Canvas */}
                <canvas ref={canvasRef} className="w-full mt-4 border" />

                <div>
                    <button
                        onClick={sendStamp}
                        className="w-full bg-blue-600 text-white py-2 rounded mt-2"
                    >
                        送信
                    </button>
                </div>

                {/* 最近のスタンプ */}
                <div className="pt-6">
                    <h2 className="text-lg font-semibold mb-2">最近作成されたスタンプ</h2>
                    <RecentStamps />
                </div>
            </div>
        </div>
    );
}

export default Generate;
