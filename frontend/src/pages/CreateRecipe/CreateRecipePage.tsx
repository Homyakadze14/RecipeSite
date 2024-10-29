import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../../components/Button/Button';
import { Complexity } from '../../components/Complexity/Complexity';
import { ImageUploader } from '../../components/ImageUploader/ImageUploader';
import { TextField } from '../../components/TextField/TextField';
import { Layout } from '../../layout/Layout';
import { useAuthStore } from '../../store/auth/useAuthStore';
import { useRecipesStore } from '../../store/recipes/useRecipesStore';
import styles from './CreateRecipePage.module.scss';

export const CreateRecipePage = () => {
	const createRecipeForm = useRecipesStore(state => state.createRecipeForm);

	const setCreateRecipeForm = useRecipesStore(
		state => state.setCreateRecipeForm
	);

	const login = useAuthStore(state => state.login);

	const createRecipe = useRecipesStore(state => state.createRecipe);

	const navigate = useNavigate();

	const handleStarsChange = (stars: string | number) => {
		setCreateRecipeForm({
			...createRecipeForm,
			complexity: Number(stars),
		});
	};

	const handleImageUpload = (image: File | null) => {
		setCreateRecipeForm({
			...createRecipeForm,
			selectedImage: image,
		});
	};

	const handleTimeChange = (
		event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
	) => {
		setCreateRecipeForm({
			...createRecipeForm,
			need_time: event.target.value,
		});
		if (event.target.value.length === 5) {
			const inputValue = event.target.value;

			const regExp = /^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9])$/;

			if (!regExp.test(inputValue)) {
				setCreateRecipeForm({
					...createRecipeForm,
					need_time: '',
				});

				console.log('Invalid time format');
				alert('Неверный формат! Введите время в формате HH:MM');
				return;
			}
		} else if (event.target.value.length > 5) {
			setCreateRecipeForm({
				...createRecipeForm,
				need_time: '',
			});

			console.log('Invalid time format');
			alert('Неверный формат! Введите время в формате HH:MM');
			return;
		}
	};

	useEffect(() => {
		console.log('createRecipeForm: ', createRecipeForm);
	});

	return (
		<Layout>
			<div className={styles.page}>
				<h2>Создание рецепта</h2>
				<form>
					<TextField
						className={styles.titleField}
						direction='row'
						label='Название:'
						field='input'
						value={createRecipeForm.title || ''}
						onChange={e => {
							setCreateRecipeForm({
								...createRecipeForm,
								title: e.target.value,
							});
						}}
						placeholder='Введите название рецепта'
					/>
					<TextField
						direction='column'
						label='Описание:'
						field='textarea'
						value={createRecipeForm.about || ''}
						onChange={e => {
							setCreateRecipeForm({
								...createRecipeForm,
								about: e.target.value,
							});
						}}
						placeholder='Введите описание рецепта'
					/>
					<TextField
						direction='column'
						label='Ингредиенты:'
						field='textarea'
						value={createRecipeForm.ingridients || ''}
						onChange={e => {
							setCreateRecipeForm({
								...createRecipeForm,
								ingridients: e.target.value,
							});
						}}
						placeholder='Введите ингредиенты, требующиеся для рецепта'
					/>
					<TextField
						direction='column'
						label='Инструкция:'
						field='textarea'
						value={createRecipeForm.instructions || ''}
						onChange={e => {
							setCreateRecipeForm({
								...createRecipeForm,
								instructions: e.target.value,
							});
						}}
						placeholder='Напишите инструкцию по приготовлению данного рецепта'
					/>
					<label className={styles.label} style={{ fontSize: '28px' }}>
						Сложность:
					</label>
					<Complexity
						key={createRecipeForm.complexity}
						starsAmount={createRecipeForm.complexity || 1}
						setStars={handleStarsChange}
						width={40}
						height={40}
						isClickable={true}
					/>
					<TextField
						direction='column'
						label='Время приготовления:'
						field='input'
						value={createRecipeForm.need_time || ''}
						onChange={handleTimeChange}
						placeholder='HH:MM'
					/>
					<ImageUploader
						key={createRecipeForm.selectedImage}
						label='Изображение:'
						selectedImage={createRecipeForm.selectedImage || null}
						onImageUpload={handleImageUpload}
					/>
					<Button
						disabled={Object.values(createRecipeForm).some(value => !value)}
						title={
							Object.values(createRecipeForm).some(value => !value)
								? `Все поля должны быть заполнены${
										createRecipeForm.need_time.length &&
										createRecipeForm.need_time.length !== 5
											? ' и время в формате HH:MM'
											: ''
								  }`
								: ''
						}
						onClick={e => {
							createRecipe(e, login, createRecipeForm, navigate);
							console.log(createRecipeForm);
						}}
					>
						Опубликовать
					</Button>
				</form>
			</div>
		</Layout>
	);
};
