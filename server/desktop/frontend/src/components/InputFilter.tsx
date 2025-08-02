export default function InputFilter({ filter, setFilter, className }: {
	filter: string,
	setFilter: (filter: string) => void,
	className?: string;
}) {
	return (
		<div className={className}>
			<input
				type="text"
				placeholder="Counter name"
				value={filter}
				className="grow"
				onChange={(e) => setFilter(e.target.value)} />

			<button
				className="min-w-[17ch] !border-1"
				onClick={() => setFilter("")}>
				Reset filter
			</button>

		</div >
	)
}

