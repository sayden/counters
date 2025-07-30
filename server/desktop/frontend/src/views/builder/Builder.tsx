import { useState, useEffect, useRef, useCallback } from 'react';


import { GetImage } from "../../../wailsjs/go/main/App";

import Header from '../../components/Header';
import CodeEditor from './CodeEditor';
import Preview from './Preview';
import CopyToClipboardButton from '../../components/CopyToClipboard';

export default function Builder() {
  const [code, setCode] = useState(``);

  // const [code, setCode] = useState('{}');
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [imageSrc, setImageSrc] = useState<string>("");
  const [scale, setScale] = useState(200);

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
    <div className='h-screen w-screen overflow-hidden'>
      <Header />
      <div className='flex flex-row h-full'>
        <div className='m-2 w-1/2 h-full overflow-auto overscroll-contain bg-gray-800'>
          <CodeEditor code={code} setCode={setCode} />
        </div>
        <div className='flex flex-col m-2 w-1/2 h-[500px] items-center'>
          <Preview imageSrc={imageSrc} />
        </div>
      </div>
    </div>
  )
}
