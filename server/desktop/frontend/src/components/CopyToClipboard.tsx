import toast, { Toaster } from 'react-hot-toast';

export default function CopyToClipboardButton({ code }: { code: string }) {
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
				className="min-w-[30ch]"
				onClick={copyToClipboard}>
				Copy to Clipboard
			</button>
			{<Toaster position="bottom-right" />}
		</div>
	);
}
