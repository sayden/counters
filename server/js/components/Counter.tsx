import React from 'react';

export interface CounterData {
	counter: string,
	id: string,
	pretty_name: string
	filename: string
}

export default function Counter({ counter }: { counter: CounterData, }) {
	console.log(counter)
	return (
		<div className='card bg-base-100 shadow-sm'>
			<figure
				className="flex flex-col items-center justify-baseline !m-0">
				<img
					id={counter.id}
					className="rounded-md m-2"
					src={counter.counter}
					alt={counter.id} />
			</figure>
			<p className="text-center mt-1 text-xs/5 text-gray-500 truncate">{counter.pretty_name}</p>
		</div>
	)
}

