export default function InputFilter({ filter, setFilter }: {
	filter: string,
	setFilter: (filter: string) => void,
}) {
	return (
		<div className="nav-input-filter">
			<input
				type="text"
				placeholder="Counter name"
				style={{ marginLeft: "1ch", flexGrow: 1 }}
				value={filter}
				onChange={(e) => setFilter(e.target.value)} />

			<button
				variant-="blue"
				style={{ minWidth: "15ch" }}
				onClick={() => setFilter("")}>
				Reset filter
			</button>

		</div >
	)
}

