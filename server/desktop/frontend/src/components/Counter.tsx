export interface CounterData {
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
		<div className='card card-xs card-border bg-base-200 shadow-sm'>
			<figure className="p-2">
				<img id={counter.id} className="rounded-md" src={counter.filename} alt={counter.id} />
			</figure>
			<div className="card-body">
				<p className="text-xs/4 text-gray-200">{counter.pretty_name}</p>
			</div>
		</div>
	)
}


