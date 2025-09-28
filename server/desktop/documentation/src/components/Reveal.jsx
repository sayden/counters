import { useState } from "react";

function Reveal({ msg = "Show", children }) {
  const [visible, setVisible] = useState(false);

  return (
    <span>
      {!visible ? (
        <button
          type="button"
          className="my-[1ch] px-[1ch]"
          onClick={() => setVisible(!visible)}
        >
          {msg}
        </button>
      ) : (
        <span>{children}</span>
      )}
    </span>
  );
}

export default Reveal;
