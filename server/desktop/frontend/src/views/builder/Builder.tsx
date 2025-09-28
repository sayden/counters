import { useState, useEffect, useRef } from 'react';
import toast, { Toaster } from 'react-hot-toast';

// Wails
import { GetImage } from "../../../wailsjs/go/backend/App";

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
          setImageSrc(`data:image/png;base64,${blob}`))
        .catch(err => {
          toast.error(err);
          console.error(err);
        });
    }, 500);
  }, [code])

  return (
    <main className="flex flex-col items-center h-screen w-screen">
      <div className='w-[80%] flex flex-col grow'>

        <Header className='flex justify-between items-baseline border-b-1 border-double !px-[1ch] !py-[1lh]' />

        <section className="flex flex-row !mt-[1lh]">
          <div className="min-w-[80%]">
            <CodeEditor code={code} setCode={setCode} />
          </div>

          <div className="flex flex-col w-1/2 justify-center items-center">
            <Preview imageSrc={imageSrc} />
          </div>
        </section>

        <Toaster
          position="bottom-right"
          toastOptions={{
            duration: 3000,
          }}
        />
      </div>
    </main>
  )
}
