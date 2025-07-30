import React from 'react';

import Counter, { CounterData } from './Counter';

export default function Grid({
	counters
}: {
	counters: CounterData[]
}) {

	return (
		<div className="grid grid-cols-12 gap-2 p-2">
			{
				counters.map((counter: CounterData) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</div>
	);
}

