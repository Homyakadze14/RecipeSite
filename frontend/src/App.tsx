import { Route, Routes } from 'react-router-dom';
import { ChangePasswordPage } from './pages/ChangePasswordPage/ChangePasswordPage';
import { CreateRecipePage } from './pages/CreateRecipe/CreateRecipePage';
import { EditRecipePage } from './pages/EditRecipePage/EditRecipePage';
import { GenTokenPage } from './pages/GenTokenPage/GenTokenPage';
import { HomePage } from './pages/Home/HomePage';
import { NotFoundPage } from './pages/NotFound/NotFoundPage';
import { RecipeCardDetailPage } from './pages/RecipeCardDetail/RecipeCardDetailPage';
import { SignInPage } from './pages/SignIn/SignInPage';
import { SignUpPage } from './pages/SignUp/SignUpPage';
import { UserPage } from './pages/User/UserPage';
import { useRecipesStore } from './store/recipes/useRecipesStore';

export const App = () => {
	const recipeId = useRecipesStore(state => state.recipeId);

	return (
		<Routes>
			<Route path='/' element={<HomePage />} />
			<Route path='/signup' element={<SignUpPage />} />
			<Route path='/signin' element={<SignInPage />} />
			<Route path='/user/:login' element={<UserPage />} />
			<Route path='/create_recipe' element={<CreateRecipePage />} />
			<Route path='/recipe/:id' element={<RecipeCardDetailPage />} />
			<Route
				path='/edit_recipe/:id'
				element={<EditRecipePage recipeId={recipeId} />}
			/>
			<Route path='/edit_password' element={<ChangePasswordPage />} />
			<Route path='/gen_token' element={<GenTokenPage />} />
			<Route path='*' element={<NotFoundPage />} />
		</Routes>
	);
};
