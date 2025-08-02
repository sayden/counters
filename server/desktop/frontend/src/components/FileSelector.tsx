interface Props {
	path: string;
	showFolderDialog: () => void;
	className?: string;
}

export default function FileSelector({ path, showFolderDialog, className }: Props) {
	return (
		<div className={className}>

			<span className="truncate">{path}</span>

			<button
				className="min-w-[16ch] !border-1"
				onClick={showFolderDialog}
			>
				Select File
			</button>
		</div>
	);
}
