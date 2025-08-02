import React from 'react'
import { createRoot } from 'react-dom/client'

// Styles
import '../../style.css'
import '../../App.css'
import "@webtui/css";

// Main component
import Help from './Markdown'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
	<React.StrictMode>
		<Help />
	</React.StrictMode>
)


