export default function FileSelector({
	path,
	showFolderDialog,
}: {
	path: string,
	showFolderDialog: () => void,
}) {
	return (
		<div className="max-w-1/2 flex items-center grow-1 justify-between gap-[1ch]">

			<p className="grow-1">{path}</p>

			<button
				variant-="blue"
				className="ml-[2ch] min-w-[15ch]"
				onClick={showFolderDialog}>
				Select File
			</button>
		</div>
	);
}
