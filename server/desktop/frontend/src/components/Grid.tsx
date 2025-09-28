import React from 'react';

import Counter, { CounterProps } from './Counter';

export default function Grid({ counters }: {
	counters: CounterProps[]
}) {

	return (
		<section className="grid grid-cols-[repeat(auto-fit,minmax(100px,1fr))] gap-[1ch]">
			{
				counters.map((counter: CounterProps) =>
					<Counter key={counter.id} counter={counter} />
				)
			}
		</section>
	);
}

