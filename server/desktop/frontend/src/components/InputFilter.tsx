export default function InputFilter({ filter, setFilter }: {
	filter: string,
	setFilter: (filter: string) => void,
}) {
	return (
		<div className="flex grow-1 w-auto justify-between items-center">
			<input
				type="text"
				placeholder="Counter name"
				className="ml-[1ch] grow-1"
				value={filter}
				onChange={(e) => setFilter(e.target.value)} />

			<button
				variant-="blue"
				className="min-w-[15ch] ml-[1ch]"
				onClick={() => setFilter("")}>
				Reset filter
			</button>

		</div >
	)
}

