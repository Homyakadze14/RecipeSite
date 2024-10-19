import { useMemo } from 'react';
import { IRecipe } from '../store/recipes/useRecipesStore';

export const useSearchedRecipes = (recipes: IRecipe[], searchQuery: string) => {
	const searchedRecipes = useMemo<IRecipe[]>(() => {
		if (!searchQuery) {
			return recipes;
		}
		return recipes.filter(
			recipe =>
				recipe.title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				recipe.about?.toLowerCase().includes(searchQuery.toLowerCase()) ||
				recipe.ingridients?.toLowerCase().includes(searchQuery.toLowerCase())
		);
	}, [searchQuery, recipes]);

	return searchedRecipes;
};
