
export default function Header({ className }: { className: string }) {
	return (
		<section className={className} >

			{/* Navigation */}
			<div className='flex justify-between gap-x-[2ch]'>
				<a href='/'>Home</a>
				<a href='/builder.html'>Builder</a>
				<a href='/help.html'>Help</a>
				<a href='/markdown.html'>Markdown</a>
			</div>

			<h1 className="!my-0">Counters Visualizer</h1>

		</section >
	)
}

