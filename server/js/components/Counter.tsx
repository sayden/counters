import React from 'react';

export interface CounterData {
	counter: string,
	id: string,
	pretty_name: string
	filename: string
}

export default function Counter({
	counter
}: {
	counter: CounterData,
}) {
	return (
		<div className='card bg-base-100 w-auto shadow-sm'>
			<figure>
				<img id={counter.id} className="rounded-md m-2" src={counter.counter} alt={counter.id} />
			</figure>
			<div className="card-body">
				<p className="mt-1 text-xs/5 text-gray-500">{counter.pretty_name}</p>
			</div>
		</div>
	)
}

