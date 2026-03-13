import { useState } from "react";
import { ClassifyFolder, Greet, Watch } from "../wailsjs/go/main/App";
import "./App.css";

function App() {
  const [resultText, setResultText] = useState("Please enter your name below 👇");
  const [name, setName] = useState("");
  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  const [isWatching, setIsWatching] = useState(false);

  function greet() {
    Greet(name).then(updateResultText);
  }

  const watch = () => {
    // 다운로드 파일 감시
    Watch();
  };

  return (
    <div id="App">
      {/* <div>
        <p>다운로드 폴더 자동 분류</p>
        <button onClick={watch}>감시</button>
      </div> */}

      <div>
        <div className="toolbar">
          <button>＋</button>
          <button>⚡</button>
          <button>⏸</button>
          <button>👁</button>
        </div>

        <div className="content">main content</div>
      </div>

      <div>
        <p>다운로드 폴더</p>
        <button onClick={() => ClassifyFolder()}>분류하기</button>
      </div>
    </div>
  );
}

export default App;
