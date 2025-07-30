
import { useState, useEffect, useCallback, useMemo } from 'react';


export default function Header() {
	return (
		<div className='flex flex-row justify-center'>

			{/* Sidebar */}
			<div className="content-center shadow-sm">
				<div className="flex-none m-2">
					<label htmlFor="my-drawer" className="btn btn-primary drawer-button">
						<svg xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							className="inline-block h-5 w-5 stroke-current">
							<path
								strokeLinecap="round"
								strokeLinejoin="round"
								strokeWidth="2"
								d="M4 6h16M4 12h16M4 18h16" />
						</svg>
					</label>
				</div>

				<div className="drawer">
					<input id="my-drawer" type="checkbox" className="drawer-toggle" />
					<div className="drawer-side">
						<label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay" />
						<ul className="menu bg-base-200 text-base-content min-h-full w-80 p-4">
							{/* Sidebar content here */}
							<li><a href='/'>Home</a></li>
							<li><a href='/builder.html'>Builder</a></li>
						</ul>
					</div>
				</div>
			</div>

			{/* Navbar */}
			<h1 className="flex flex-none p-2 mx-2 text-2xl content-center text-center">
				Counters Visualizer
			</h1>
		</div>
	)
}

