import React from 'react';

import Counter, { CounterData } from './Counter';

export default function Grid({ counters }: {
	counters: CounterData[]
}) {

	return (
		<div className="grid-counters">
			{
				counters.map((counter: CounterData) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</div>
	);
}

