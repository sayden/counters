import { useRef, useCallback } from 'react';

export default function Preview({ imageSrc }: { imageSrc: string }) {
  return (
    <img src={imageSrc} alt="Preview" />
  );
}

function ScaleDropdownSelector({ setScale }: { setScale: (n: number) => void }) {
  const detailsRef = useRef<HTMLDetailsElement>(null);

  const handleClick = useCallback((n: number) => {
    setScale(n);
    if (detailsRef.current) {
      detailsRef.current.open = false; // Closes the details dropdown
    }
  }, []);

  return (
    <details className="dropdown" ref={detailsRef}>
      <summary className="btn m-1">Counter scale</summary>
      <ul className="menu dropdown-content bg-base-100 rounded-box z-1 w-52 shadow-sm">
        <li><a onClick={() => handleClick(200)}>200%</a></li>
        <li><a onClick={() => handleClick(100)}>100%</a></li>
        <li><a onClick={() => handleClick(50)}>50%</a></li>
        <li><a onClick={() => handleClick(300)}>300%</a></li>
        <li><a onClick={() => handleClick(400)}>400%</a></li>
      </ul>
    </details>
  );
}
