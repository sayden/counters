import { useState, useEffect, useRef } from 'react';

// Wails
import { GetImage } from "../../../wailsjs/go/main/App";

// Components
import Header from '../../components/Header';
import CodeEditor from './CodeEditor';
import Preview from './Preview';

export default function Builder() {
  const [code, setCode] = useState(``);

  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [imageSrc, setImageSrc] = useState<string>("");

  useEffect(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(async () => {
      // console.log("timeoutRef.current", timeoutRef.current);
      await GetImage(code)
        .then(blob =>
          setImageSrc(`data:image/png;base64,${blob}`));
    }, 500);
  }, [code])

  return (
    <div style={{ height: "100vh", width: "100vw", overflow: "hidden" }}>

      <Header />

      <div style={{ display: "flex", flex: "row", height: "100%" }}>

        <div
          style={{ width: "50%", height: "100%", overflow: "auto", overscrollBehavior: "contain" }}>
          <CodeEditor code={code} setCode={setCode} />
        </div>

        <div style={{ display: "flex", flex: "col", width: "50%", height: "50%", alignItems: "center" }}>
          <Preview imageSrc={imageSrc} />
        </div>

      </div>

    </div>
  )
}
