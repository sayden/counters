import React from 'react';

import Counter, { CounterData } from './Counter';

export default function Grid({ counters }: {
	counters: CounterData[]
}) {

	return (
		<section className="grid grid-cols-[repeat(auto-fit,minmax(100px,1fr))] gap-[1ch]">
			{
				counters.map((counter: CounterData) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</section>
	);
}

