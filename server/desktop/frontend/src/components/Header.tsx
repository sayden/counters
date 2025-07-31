
import { useState, useEffect, useCallback, useMemo } from 'react';


export default function Header() {
	return (
		<div box-="square" className="nav-title-and-links">

			{/* Navigation */}
			<div style={{ display: "flex", justifyContent: "space-between", gap: "2ch" }}>
				<p><a href='/'>Home</a></p>
				<p><a href='/builder.html'>Builder</a></p>
			</div>

			<div style={{ display: "flex", justifyContent: "flex-end" }}>
				<h1>Counters Visualizer</h1>
			</div>

		</div >
	)
}

