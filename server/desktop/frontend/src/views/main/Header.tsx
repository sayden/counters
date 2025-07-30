import { useState, useEffect, useCallback, useMemo } from 'react';

// Components
import InputFilter from '../../components/InputFilter';
import FileSelector from '../../components/FileSelector';
import HeaderComponent from '../../components/Header';


export default function Header({ filter, setFilter, path, showFolderDialog }: {
	filter: string,
	setFilter: (filter: string) => void,
	path: string,
	showFolderDialog: () => void
}) {
	return (
		<div className='flex flex-row bg-gray-900'>
			<HeaderComponent />

			<InputFilter
				filter={filter}
				setFilter={setFilter} />

			<FileSelector
				path={path}
				showFolderDialog={showFolderDialog} />
		</div>
	)
}

