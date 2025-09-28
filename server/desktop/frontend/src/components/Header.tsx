
export default function Header({ className }: { className: string }) {
	return (
		<section className={className} >

			{/* Navigation */}
			<div className='flex justify-between gap-x-[2ch]'>
				<a href='/'>Home</a>
				<a href='/src/views/builder/index.html'>Builder</a>
				<a href='documentation/dist/index.html'>Documentation</a>
			</div>

			<h1 className="!my-0">Counters Visualizer</h1>

		</section >
	)
}

