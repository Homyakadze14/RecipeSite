import { ChangeEvent } from 'react';
import { SearchIcon } from '../../assets/icons/SearchIcon';
import styles from './SearchInput.module.scss';

export interface ISearchInput {
	value: string;
	onChange: (e: ChangeEvent<HTMLInputElement>) => void;
}

export const SearchInput = ({ value, onChange }: ISearchInput) => {
	return (
		<div className={styles.searchInput}>
			<input
				placeholder='Введите текст описания или названия рецепта'
				value={value}
				onChange={onChange}
			/>
			<SearchIcon />
		</div>
	);
};
