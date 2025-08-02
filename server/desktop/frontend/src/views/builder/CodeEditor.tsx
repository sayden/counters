import React, { useState, useMemo, useCallback } from "react";

import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { dracula } from '@uiw/codemirror-theme-dracula';
import prettier from "prettier/standalone";
import parserBabel from "prettier/plugins/babel";
import estreePlugin from "prettier/plugins/estree";
import toast, { Toaster } from 'react-hot-toast';
import ButtonCopyToClipboard from "./components/CopyToClipboard";

interface Props {
  code: string,
  setCode: (n: string) => void
}

export default function CodeEditor({ code, setCode }: Props) {
  const handleFormat = useCallback(() => {
    console.log("formatting...", code);
    prettier.format(code,
      {
        parser: "json",
        plugins: [parserBabel, estreePlugin],
        semi: true,
        singleQuote: false,
      })
      .then(setCode)
      .catch(() => toast.error('Failed to copy'));
  }, [code]);

  return (
    <section>
      <div className="flex flex-row justify-between !mb-[1lh] border-b-1 border-dotted !pb-[1lh]">
        <ButtonCopyToClipboard code={code} />
        <button
          className="min-w-[15ch] !border-1"
          onClick={handleFormat}>
          Format
        </button>
      </div>
      <CodeMirror
        value={code}
        onChange={setCode}
        extensions={[json()]}
        theme={dracula}
        className='text-sm overflow-y-auto overscroll-contain'
      />
      <Toaster position="bottom-right" />
    </section>
  );
}
