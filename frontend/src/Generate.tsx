import React, { useState } from 'react';
import axios from "axios";
import { ColorResult, SketchPicker } from 'react-color';



function Generate() {

    const [text, setText] = useState<string>(`HELLO`);
    const [textColor, setTextColor] = useState('#FFFFFF');
    const [language, setLanguage] = useState<'japanese' | 'chinese'>('japanese');
    const changeTextColor = (color: ColorResult) => {
        setTextColor(color.hex)
    }

    // URLのstate
    const [imageUrl, setImageUrl] = useState<string | null>(null);

    const generateImage = async () => {
        try {
            const response = await axios.post("http://localhost:8080/api/generate", {
                text,
                textColor,
                language
            }, { responseType: "blob" }); // 画像をバイナリデータで受け取る

            const imageBlob = new Blob([response.data], { type: "image/png" });
            setImageUrl(URL.createObjectURL(imageBlob));
        } catch (error) {
            console.error("エラー:", error);
            alert("画像生成に失敗しました");
        }
    };


    return (
        <div>
            <div>
                <button onClick={() => setLanguage('japanese')} disabled={language === 'japanese'}>
                    日本語
                </button>
                <button onClick={() => setLanguage('chinese')} disabled={language === 'chinese'}>
                    中国語
                </button>
            </div>

            <form onSubmit={(e) => e.preventDefault()}>
                <input
                    type="text"
                    // text ステートが持っている入力中テキストの値を value として表示
                    value={text}
                    // onChange イベント（＝入力テキストの変化）を text ステートに反映する
                    onChange={(e) => setText(e.target.value)}
                />
                <input type="submit" />  {/* ← 省略 */}
            </form>

            {/* ↓ DOM のリアクティブな反応を見るためのサンプル */}
            <p>{text}</p>
            {/* ↑ あとで削除 */}

            <div>
                <SketchPicker
                    color={textColor}
                    onChange={changeTextColor}
                />
            </div>

            <div style={{ color: `${textColor}` }}>{text}</div>

            <button onClick={generateImage}>画像を生成</button>
            <br />
            {imageUrl && <img src={imageUrl} alt="Generated Image" />}
        </div>



    );
}

export default Generate;

