import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { RecipeCardDetail } from '../../components/RecipesList/RecipeCard/RecipeCardDetail/RecipeCardDetail';
import { Layout } from '../../layout/Layout';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';

export const RecipeCardDetailPage = () => {
	const getRecipeAndComments = useRecipesStore(
		state => state.getRecipeAndComments
	);
	const recipe = useRecipesStore(state => state.recipe);

	const isLoading = useRecipesStore(state => state.isLoading);
	const isLiked = useRecipesStore(state => state.isLiked);

	const location = useLocation();
	const recipeId = Number(location.pathname.split('/').pop());

	useEffect(() => {
		console.log('recipeId: ', recipeId);
	}, [recipeId]);

	useEffect(() => {
		getRecipeAndComments(recipeId);
	}, [isLoading, isLiked]);

	return (
		<Layout>
			<RecipeCardDetail recipe={recipe} />
		</Layout>
	);
};
