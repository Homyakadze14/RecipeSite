import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { Complexity } from '../../components/Complexity/Complexity';
import { ImageUploader } from '../../components/ImageUploader/ImageUploader';
import { TextField } from '../../components/TextField/TextField';
import { Layout } from '../../layout/Layout';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';
import { IEditRecipePage } from '../../types/interfaces';
import styles from './EditRecipePage.module.scss';

export const EditRecipePage = ({ recipeId }: IEditRecipePage) => {
	const login = useAuthStore(state => state.login);

	const isLoading = useRecipesStore(state => state.isLoading);
	const setIsLoading = useRecipesStore(state => state.setIsLoading);

	const editRecipeForm = useRecipesStore(state => state.editRecipeForm);
	const setEditRecipeForm = useRecipesStore(state => state.setEditRecipeForm);

	const recipe = useRecipesStore(state => state.recipe);

	const getRecipeAndComments = useRecipesStore(
		state => state.getRecipeAndComments
	);
	const editRecipe = useRecipesStore(state => state.editRecipe);

	const navigate = useNavigate();

	useEffect(() => {
		console.log('editRecipeForm: ', editRecipeForm);
	}, [editRecipeForm]);

	useEffect(() => {
		getRecipeAndComments(recipeId);
	}, [recipeId]);

	useEffect(() => {
		setEditRecipeForm({
			...recipe,
			photos: recipe.photos_urls,
		});

		setIsLoading(false);
	}, [isLoading, recipe]);

	const handleStarsChange = (stars: string | number) => {
		setEditRecipeForm({
			...editRecipeForm,
			complexity: Number(stars),
		});
	};

	const handleImageUpload = (image: File | null) => {
		setEditRecipeForm({
			...editRecipeForm,
			photos: image,
		});
	};

	const TIME_FORMAT_REGEXP = /^([0-9]|0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$/;

	const handleTimeChange = (
		event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => {
		setEditRecipeForm({
			...editRecipeForm,
			need_time: event.target.value,
		});
		if (event.target.value.length === 5) {
			const inputValue = event.target.value;

			if (!TIME_FORMAT_REGEXP.test(inputValue)) {
				setEditRecipeForm({
					...editRecipeForm,
					need_time: '',
				});

				alert('Неверный формат! Введите время в формате HH:MM');
				return;
			}
		}
	};

	return (
		<Layout>
			<div className={styles.page}>
				<h2>Редактирование рецепта</h2>
				<form>
					<TextField
						direction='row'
						label='Название:'
						field='input'
						value={editRecipeForm.title || ''}
						onChange={e => {
							setEditRecipeForm({
								...editRecipeForm,
								title: e.target.value,
							});
						}}
						placeholder='Введите название рецепта'
					/>
					<TextField
						direction='column'
						label='Описание:'
						field='textarea'
						value={editRecipeForm.about || ''}
						onChange={e => {
							setEditRecipeForm({
								...editRecipeForm,
								about: e.target.value,
							});
						}}
						placeholder='Введите описание рецепта'
					/>
					<TextField
						direction='column'
						label='Ингредиенты:'
						field='textarea'
						value={editRecipeForm.ingridients || ''}
						onChange={e => {
							setEditRecipeForm({
								...editRecipeForm,
								ingridients: e.target.value,
							});
						}}
						placeholder='Введите ингредиенты, требующиеся для рецепта'
					/>
					<TextField
						direction='column'
						label='Инструкция:'
						field='textarea'
						value={editRecipeForm.instructions || ''}
						onChange={e => {
							setEditRecipeForm({
								...editRecipeForm,
								instructions: e.target.value,
							});
						}}
						placeholder='Напишите инструкцию по приготовлению данного рецепта'
					/>
					<label className={styles.label} style={{ fontSize: '28px' }}>
						Сложность:
					</label>
					<Complexity
						key={editRecipeForm.complexity}
						starsAmount={editRecipeForm.complexity || 1}
						setStars={handleStarsChange}
						width={40}
						height={40}
						isClickable={true}
					/>
					<TextField
						direction='column'
						label='Время приготовления:'
						field='input'
						value={editRecipeForm.need_time || ''}
						onChange={handleTimeChange}
						placeholder='HH:MM'
					/>
					<ImageUploader
						key={editRecipeForm.photos?.name || Math.random()}
						label='Изображение:'
						selectedImage={editRecipeForm.photos || null}
						onImageUpload={handleImageUpload}
					/>
					<Button
						disabled={
							Object.values(editRecipeForm).some(value => !value) ||
							editRecipeForm.need_time.length !== 5
						}
						title={
							Object.values(editRecipeForm).some(value => !value)
								? `Все поля должны быть заполнены${
										editRecipeForm.need_time?.length &&
										editRecipeForm.need_time?.length !== 5
											? ' и время в формате HH:MM'
											: ''
								  }`
								: ''
						}
						onClick={e =>
							editRecipe(e, login, recipeId, editRecipeForm, navigate)
						}
						className={styles.saveButton}
					>
						Применить
					</Button>
				</form>
			</div>
		</Layout>
	);
};
