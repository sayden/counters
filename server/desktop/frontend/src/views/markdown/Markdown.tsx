import { useEffect, useId, useState } from 'react';
import Header from '../../components/Header';

export default function Help() {
	const [markdown, setMarkdown] = useState('');

	useEffect(() => {
		fetch('/api/md.html', {
			headers: {
				'Accept': 'text/html'
			},
			method: 'GET'
		})
			.then(response => response.text())
			.then(text => setMarkdown(text))
	}, [])
	return (
		<main className="flex flex-col items-center">
			<div className='w-[80%] flex flex-col grow'>
				<Header className='flex justify-between items-baseline border-b-1 border-double !px-[1ch] !py-[1lh]' />
				<section>
					<div dangerouslySetInnerHTML={{ __html: markdown }} />;
				</section>
			</div>
		</main>
	)
}
