import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { RecipesList } from '../../components/RecipesList/RecipesList';
import { Tabs } from '../../components/Tabs/Tabs';
import { UserCard } from '../../components/UserCard/UserCard';
import { Layout } from '../../layout/Layout';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';
import { deleteLastChar, useUsersStore } from '../../store/users/useUsersStore';
import { NotFoundPage } from '../NotFound/NotFoundPage';
import styles from './UserPage.module.scss';

export const UserPage = () => {
	const user = useUsersStore(state => state.user);
	const login = useAuthStore(state => state.login);
	const getUser = useUsersStore(state => state.getUser);
	const isLoading = useRecipesStore(state => state.isLoading);
	const setIsLoading = useRecipesStore(state => state.setIsLoading);
	const isLiked = useRecipesStore(state => state.isLiked);
	const paramsLogin = useUsersStore(state => state.paramsLogin);
	const setParamsLogin = useUsersStore(state => state.setParamsLogin);

	useEffect(() => console.log('LOGIN !!!!!!!!!', login), [login]);

	const [selectedTab, setSelectedTab] = useState(
		localStorage.getItem('selectedTab') || 'added'
	);

	const navigate = useNavigate();

	const params = useParams();

	useEffect(() => {
		setParamsLogin(params.login || '');
	}, [useParams()]);

	useEffect(() => {
		getUser(paramsLogin);

		const timer = setTimeout(() => {
			setIsLoading(false);
		}, 300);

		console.log('USER IS: ', user);

		return () => clearTimeout(timer);
	}, [paramsLogin, user.login, user.about, user.icon_url, isLoading, isLiked]);

	useEffect(() => {
		const storedTab = localStorage.getItem('selectedTab');
		console.log('stored tab: ', storedTab);
		if (storedTab) {
			setSelectedTab(localStorage.getItem('selectedTab') || 'none');
		}
	}, []);

	useEffect(() => {
		localStorage.setItem('selectedTab', selectedTab);

		console.log('selected tab: ', selectedTab);
	}, [selectedTab]);

	if (paramsLogin != user?.login) {
		setTimeout(() => <NotFoundPage />, 500);
	}

	return (
		<Layout>
			<div className={styles.userPage}>
				<UserCard
					avatarSrc={user.icon_url}
					name={paramsLogin}
					description={user?.about || 'Ещё нет описания'}
				/>
				{login === user.login && (
					<div className={styles.createRecipeContainer}>
						<Tabs
							tabs={{ added: 'Добавленные', liked: 'Понравившиеся' }}
							selectedTab={selectedTab}
							setSelectedTab={setSelectedTab}
						/>
						<button
							className={styles.createRecipeButton}
							onClick={() => navigate('/create_recipe')}
						>
							Добавить новый рецепт
						</button>
					</div>
				)}

				<h2>Рецепты</h2>
				<RecipesList
					data={
						selectedTab === 'added'
							? user?.recipies
							: user?.liked_recipies?.map(recipe => ({
									...recipe,
									photos_urls: deleteLastChar(recipe.photos_urls),
							  })) || []
					}
					noDataText={
						selectedTab === 'added'
							? login === user.login
								? 'Вы ещё не добавили рецепты'
								: 'Пользователь ещё не добавил рецепты'
							: login === user.login
							? 'Вы ещё не лайкнули рецепты'
							: 'Пользователь ещё не лайкнул рецепты'
					}
					isTabChanged={selectedTab === 'added' ? false : true}
					listStartY={360}
				/>
			</div>
		</Layout>
	);
};