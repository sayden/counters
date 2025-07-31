export default function FileSelector({
	path,
	showFolderDialog,
}: {
	path: string,
	showFolderDialog: () => void,
}) {
	return (
		<div className="nav-file-selector">

			<p style={{ flexGrow: 1 }}>{path}</p>

			<button
				variant-="blue"
				style={{ marginLeft: "2ch", minWidth: "15ch" }}
				onClick={showFolderDialog}>
				Select File
			</button>
		</div>
	);
}
