export default function FileSelector(
	{
		path,
		showFolderDialog,
	}: {
		path: string,
		showFolderDialog: () => void,
	}) {
	return (
		<div className="flex flex-row">
			<div id="result" className="flex-4 px-4 py-2 m-2 text-sm text-gray-500 rounded-md truncate w-3/4">
				{path}
			</div>
			<button
				className="btn btn-primary flex-1 my-2 bg-gray-800 hover:bg-gray-700 p-2 w-1/4"
				onClick={showFolderDialog}
			>
				Select File
			</button>
		</div>
	);
}
