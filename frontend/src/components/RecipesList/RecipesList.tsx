import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import { ArrowIcon } from '../../assets/icons/ArrowIcon';
import { RecipeCardDetailPage } from '../../pages/RecipeCardDetail/RecipeCardDetailPage';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';
import {
	IRecipe,
	IRecipeCardDetail,
	IRecipesList,
} from '../../types/interfaces';
import { RecipeCard } from './RecipeCard/RecipeCard';
import { RecipeCardSkeletonLoader } from './RecipeCard/RecipeSkeletonLoader/RecipeSkeletonLoader';
import styles from './RecipesList.module.scss';

export const RecipesList = ({
	data = [],
	noDataText = 'Ещё нет рецептов',
	isSearchQueryChanged = false,
	isTabChanged,
	listStartY = 0,
}: IRecipesList) => {
	const login = useAuthStore(state => state.login);

	const deleteRecipe = useRecipesStore(state => state.deleteRecipe);

	const isLoading = useRecipesStore(state => state.isLoading);

	const setRecipe = useRecipesStore(state => state.setRecipe);

	const [selectedRecipe, setSelectedRecipe] =
		useState<IRecipeCardDetail | null>(null);

	const paramsLogin = useParams().login || '';

	const [currentPage, setCurrentPage] = useState(() => {
		const storedPage = localStorage.getItem('currentPage');
		return storedPage ? parseInt(storedPage, 10) : 1;
	});

	const recipesPerPage = 10;

	const indexOfLastRecipe = currentPage * recipesPerPage;
	const indexOfFirstRecipe = indexOfLastRecipe - recipesPerPage;
	const currentRecipes = data.slice(indexOfFirstRecipe, indexOfLastRecipe);

	useEffect(() => {
		localStorage.setItem('currentPage', String(currentPage));
		window.scrollTo(0, listStartY);
	}, [currentPage]);

	const nextPage = () => {
		if (currentPage < Math.ceil(data.length / recipesPerPage)) {
			setCurrentPage(currentPage + 1);
		}
	};

	const prevPage = () => {
		if (currentPage > 1) {
			setCurrentPage(currentPage - 1);
		}
	};

	const navigate = useNavigate();

	const handleCardClick = (recipe: IRecipeCardDetail) => {
		setSelectedRecipe(recipe);
		const recipeForSetRecipe: IRecipe = {
			...recipe,
			created_at: '',
			ingridients: '',
			updated_at: '',
		};
		setRecipe(recipeForSetRecipe);
		navigate(`/recipe/${recipe.id}`);
	};

	useEffect(() => {
		if (isSearchQueryChanged || isTabChanged || data.length === 0) {
			setCurrentPage(1);
		}
		console.log('paramsLogin in list: ', paramsLogin);
	}, [isSearchQueryChanged, isTabChanged, data.length, paramsLogin]);

	return (
		<div className={styles.container}>
			{!selectedRecipe && currentRecipes.length > 0 && (
				<TransitionGroup className={styles.list} component='ul'>
					{currentRecipes.map(item => (
						<CSSTransition
							key={item.id}
							timeout={500}
							classNames={{
								enter: styles.fadeEnter,
								enterActive: styles.fadeEnterActive,
								exit: styles.fadeExit,
								exitActive: styles.fadeExitActive,
							}}
							mountOnEnter
							unmountOnExit
						>
							<RecipeCard
								id={item.id}
								imageSrc={item.photos_urls}
								title={item.title}
								author={item.author}
								description={item.about}
								time={item.need_time}
								starsAmount={item.complexity}
								onClick={() => {
									handleCardClick({
										...item,
										ingredients: item.ingridients,
									});
								}}
								onDelete={() => deleteRecipe(login, item.id)}
							/>
						</CSSTransition>
					))}
				</TransitionGroup>
			)}

			{isLoading &&
				[...Array(recipesPerPage)].map((_, i) => (
					<RecipeCardSkeletonLoader key={i} />
				))}

			{data.length > 0 && (
				<div className={styles.paginationContainer}>
					<button
						onClick={prevPage}
						disabled={currentPage === 1}
						style={{ cursor: currentPage === 1 ? 'not-allowed' : 'pointer' }}
					>
						<ArrowIcon direction='left' fill='#fff' />
					</button>
					<span>
						{`${currentPage} из ${Math.ceil(data.length / recipesPerPage)}`}
					</span>
					<button
						onClick={nextPage}
						disabled={currentPage === Math.ceil(data.length / recipesPerPage)}
						style={{
							cursor:
								currentPage === Math.ceil(data.length / recipesPerPage)
									? 'not-allowed'
									: 'pointer',
						}}
					>
						<ArrowIcon direction='right' fill='#fff' />
					</button>
				</div>
			)}

			{!data.length && !isLoading && (
				<p className={styles.noDataText}>{noDataText}</p>
			)}

			{selectedRecipe && <RecipeCardDetailPage />}
		</div>
	);
};
