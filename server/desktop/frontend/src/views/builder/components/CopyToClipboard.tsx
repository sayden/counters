import toast from 'react-hot-toast';

export default function ButtonCopyToClipboard({ code }: { code: string }) {
	const copyToClipboard = () => {
		navigator.clipboard.writeText(code)
			.then(() => {
				console.log('Copied to clipboard!');
				toast.success('Copied to clipboard!');
			})
			.catch(() => toast.error('Failed to copy'));
	};

	return (
		<div>
			<button
				className="min-w-[30ch] !border-1"
				onClick={copyToClipboard}>
				Copy to Clipboard
			</button>
		</div>
	);
}
