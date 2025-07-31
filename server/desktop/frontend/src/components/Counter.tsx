export interface CounterData {
	id: string,
	pretty_name: string
	filename: string
}

export default function Counter({ counter }: { counter: CounterData }) {
	return (
		<div className='counter'>
			<figure
				style={{ display: "flex", justifyContent: "center" }}>
				<img
					id={counter.id}
					style={{ borderRadius: "1ch" }}
					src={counter.filename}
					alt={counter.id} />
			</figure>
			<div>
				<p className="counter-text">{counter.pretty_name}</p>
			</div>
		</div>
	)
}


