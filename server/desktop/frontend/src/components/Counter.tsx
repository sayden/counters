export interface CounterData {
	id: string,
	pretty_name: string
	filename: string
}

export default function Counter({ counter }: { counter: CounterData }) {
	return (
		<figure
			className="flex flex-col items-center justify-baseline !m-0">
			<img
				className="rounded-[1ch] !border-none !m-0"
				id={counter.id}
				src={counter.filename}
				alt={counter.id} />
			<span className="text-center text-wrap text-xs">{counter.pretty_name}</span>
		</figure>
	)
}


