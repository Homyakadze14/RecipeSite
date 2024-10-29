import { useEffect, useState } from 'react';
import { SearchIcon } from '../../assets/icons/SearchIcon';
import { ISearchInput } from '../../types/interfaces';
import styles from './SearchInput.module.scss';

export const SearchInput = ({ value, onChange }: ISearchInput) => {
	const [placeholder, setPlaceholder] = useState(
		'Введите описание или название рецепта'
	);

	useEffect(() => {
		const handleResize = () => {
			if (window.innerWidth <= 520 && window.innerWidth > 430) {
				setPlaceholder('Описание или название рецепта');
			} else if (window.innerWidth <= 430) {
				setPlaceholder('Описание или название');
			} else {
				setPlaceholder('Введите описание или название рецепта');
			}
		};

		window.addEventListener('resize', handleResize);

		handleResize();

		return () => {
			window.removeEventListener('resize', handleResize);
		};
	}, []);

	return (
		<div className={styles.searchInput}>
			<input placeholder={placeholder} value={value} onChange={onChange} />
			<SearchIcon />
		</div>
	);
};
