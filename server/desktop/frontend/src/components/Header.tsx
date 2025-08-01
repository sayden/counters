
import { useState, useEffect, useCallback, useMemo } from 'react';


export default function Header() {
	return (
		<div box-="square" className="flex justify-between p-[3ch]" >

			{/* Navigation */}
			<div
				className='flex justify-between gap-[2ch]'>
				<p><a href='/'>Home</a></p>
				<p><a href='/builder.html'>Builder</a></p>
			</div>

			<div className='flex justify-end'>
				<h1>Counters Visualizer</h1>
			</div>

		</div >
	)
}

