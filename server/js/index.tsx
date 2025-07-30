import { createRoot } from 'react-dom/client';
import Grid from './components/Grid';

const headerRoot = document.getElementById('react');
if (!headerRoot) {
	throw new Error('Could not find element with id react');
}

const headerReactRoot = createRoot(headerRoot);
headerReactRoot.render(<Grid />);
