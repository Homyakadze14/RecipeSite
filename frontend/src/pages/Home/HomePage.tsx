import { ChangeEvent, useEffect, useState } from 'react';
import { RecipesList } from '../../components/RecipesList/RecipesList';
import { SearchInput } from '../../components/SearchInput/SearchInput';
import { Layout } from '../../layout/Layout';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';
import styles from './HomePage.module.scss';

export const HomePage = () => {
	const recipes = useRecipesStore(state => state.recipes);

	const isLoading = useRecipesStore(state => state.isLoading);

	const [searchQuery, setSearchQuery] = useState('');
	const [prevSearchQuery, setPrevSearchQuery] = useState('');

	const getFilteredRecipes = useRecipesStore(state => state.getFilteredRecipes);

	useEffect(() => {
		getFilteredRecipes({
			order_by: 1,
			query: searchQuery,
		});

		setPrevSearchQuery(searchQuery);
	}, [searchQuery, isLoading]);

	return (
		<Layout>
			<div className={styles.page}>
				<div className={styles.searchInputContainer}>
					<SearchInput
						value={searchQuery}
						onChange={(e: ChangeEvent<HTMLInputElement>) =>
							setSearchQuery(e.target.value)
						}
					/>
				</div>
				<div style={{ width: '100%', maxWidth: '100%' }}>
					<RecipesList
						data={recipes}
						isSearchQueryChanged={searchQuery !== prevSearchQuery}
					/>
				</div>
			</div>
		</Layout>
	);
};
