import React from 'react';

import Counter, { CounterData } from './Counter';

export default function Grid({ counters }: {
	counters: CounterData[]
}) {

	return (
		<div className="gap-[1ch] grid grid-cols-[repeat(auto-fit,minmax(150px,1fr))]">
			{
				counters.map((counter: CounterData) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</div>
	);
}

