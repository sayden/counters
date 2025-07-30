import { useState, useEffect, useCallback, useMemo } from 'react';

export default function InputFilter({
	filter, setFilter
}: {
	filter: string,
	setFilter: (filter: string) => void,
}) {

	return (
		<div className='flex flex-row'>
			<input
				className="input input-primary m-2 shrink"
				type="text"
				placeholder="Filter"
				value={filter}
				onChange={(e) => setFilter(e.target.value)} />
			<button
				className="btn btn m-2 p-2 hover:bg-gray-700 bg-gray-800 flex-none"
				onClick={() => setFilter("")}
			>
				Reset filter
			</button>
		</div>
	)
}

