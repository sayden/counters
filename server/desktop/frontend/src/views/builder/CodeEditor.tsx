import React, { useState, useMemo, useCallback } from "react";

import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { dracula } from '@uiw/codemirror-theme-dracula';
import CopyToClipboardButton from "../../components/CopyToClipboard";
import prettier from "prettier/standalone";
import parserBabel from "prettier/plugins/babel";
import estreePlugin from "prettier/plugins/estree";
import toast, { Toaster } from 'react-hot-toast';

interface Props {
  code: string,
  setCode: (n: string) => void
}

export default function CodeEditor({ code, setCode }: Props) {
  const handleFormat = useCallback(() => {
    prettier.format(code,
      {
        parser: "json",
        plugins: [parserBabel, estreePlugin],
        semi: true,
        singleQuote: false,
      })
      .then(formatted => {
        setCode(formatted);
        navigator.clipboard.writeText(code)
          .then(() => toast.success('Formatted!'));
      }).catch(() => toast.error('Failed to copy'));
  }, []);

  return (
    <div>
      <div className="flex flex-row justify-between">
        <CopyToClipboardButton code={code} />
        <button onClick={handleFormat}>
          Format
        </button>
      </div>
      <CodeMirror
        value={code}
        onChange={setCode}
        extensions={[json()]}
        theme={dracula}
        className='h-full'
      />
      <Toaster position="bottom-right" />
    </div>
  );
}
