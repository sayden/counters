// Components
import InputFilter from '../../components/InputFilter';
import FileSelector from '../../components/FileSelector';
import HeaderComponent from '../../components/Header';


export default function Header({
	filter,
	setFilter,
	path,
	showFolderDialog
}: {
	filter: string,
	setFilter: (filter: string) => void,
	path: string,
	showFolderDialog: () => void
}) {
	return (
		<nav className='!pb-[1ch]'>

			<HeaderComponent
				className="flex justify-between items-center border-b-1 border-double !px-[1ch] !py-[1lh]" />

			<section className="flex justify-stretch gap-[2ch] border-b-1 border-dotted !px-[1ch] !pb-[1lh] !mt-[1lh]">
				<FileSelector
					className="flex items-center justify-between gap-[1ch] max-w-1/2 text-wrap"
					path={path}
					showFolderDialog={showFolderDialog} />

				<InputFilter
					className="flex grow justify-between items-center gap-[1ch]"
					filter={filter}
					setFilter={setFilter} />
			</section>

		</nav>
	)
}

