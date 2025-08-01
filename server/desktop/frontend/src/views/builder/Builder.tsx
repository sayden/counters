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
    <div>

      <Header />

      <div className="flex flex-row h-full">
        <div className="min-w-1/2 h-full overflow-auto overscroll-contain">
          <CodeEditor code={code} setCode={setCode} />
        </div>

        <div className="flex flex-col w-1/2 h-1/2 items-center">
          <Preview imageSrc={imageSrc} />
        </div>
      </div>

    </div>
  )
}
