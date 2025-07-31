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
		<nav>

			<HeaderComponent />

			<div box-="square" className="topbar">
				<FileSelector
					path={path}
					showFolderDialog={showFolderDialog} />

				<InputFilter
					filter={filter}
					setFilter={setFilter} />
			</div>

		</nav>
	)
}

