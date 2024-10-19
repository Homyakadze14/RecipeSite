import { useEffect, useState } from 'react';
import { useLocation } from 'react-router-dom';

export const useCurrentPath = () => {
	const location = useLocation();
	const [currentPath, setCurrentPath] = useState(location.pathname);

	useEffect(() => {
		const handlePopState = () => {
			setCurrentPath(window.location.pathname);
		};

		window.addEventListener('popstate', handlePopState);

		return () => {
			window.removeEventListener('popstate', handlePopState);
		};
	}, []);

	return currentPath;
};
