import axios, { AxiosError } from 'axios';
import { create } from 'zustand';
import {
	IAuthor,
	IComment,
	IRecipe,
	IUseRecipesStore,
} from '../../types/interfaces';
import { Filter } from '../../types/types';
import { deleteLastChar } from '../../utils/utils';
import { baseUrl } from '../auth/useAuthStore';

export const handleError = (error: unknown) => {
	if (axios.isAxiosError(error)) {
		const axiosError = error as AxiosError;
		const statusCode =
			axiosError?.response?.status ?? axiosError?.request?.status ?? 0;
		const errorMessage = getErrorMessage(statusCode);
		return errorMessage;
	} else if (typeof error === 'string') {
		const statusCodeMatch = error.match(/\b(4\d{2}|5\d{2})\b/);
		if (statusCodeMatch) {
			const statusCode = parseInt(statusCodeMatch[0]);
			const errorMessage = getErrorMessage(statusCode);
			return errorMessage;
		} else {
			return 'Произошла неизвестная ошибка.';
		}
	} else {
		return 'Произошла неизвестная ошибка.';
	}
};

const getErrorMessage = (statusCode: number) => {
	switch (statusCode) {
		case 400:
			return 'Некорректные данные запроса. Пожалуйста, проверьте введенные данные.';
		case 401:
			return 'Пользователь не авторизован. Пожалуйста, войдите в аккаунт.';
		case 404:
			return 'Запрошенный ресурс не найден. Пожалуйста, проверьте URL или введенные данные.';
		case 500:
			return 'Внутренняя ошибка сервера. Пожалуйста, попробуйте снова позже.';
		default:
			return `Произошла неизвестная ошибка с кодом ${statusCode}.`;
	}
};

export const useRecipesStore = create<IUseRecipesStore>((set, get) => ({
	isLoading: true,
	setIsLoading: isLoading => set({ isLoading }),
	isLiked: true,
	setIsLiked: isLiked => set({ isLiked }),
	likes_count: 0,
	setLikesCount: likes_count => set({ likes_count }),
	recipe: {} as IRecipe,
	setRecipe: (recipe: IRecipe) => set({ recipe }),
	recipes: [],
	setRecipes: (recipes: IRecipe[]) => set({ recipes }),
	recipeId: Number(localStorage.getItem('recipeId')) || -1,

	comments: [] as IComment[],
	setComments: (comments: IComment[]) => set({ comments }),
	text: '',
	setText: text => set({ text }),

	setRecipeId: recipeId => set({ recipeId: recipeId }),

	createRecipeForm: {
		title: '',
		about: '',
		ingridients: '',
		author: {} as IAuthor,
		instructions: '',
		complexity: 1,
		need_time: '',
		selectedImage: null,
	},
	setCreateRecipeForm: form => set({ createRecipeForm: form }),

	editRecipeForm: {
		photos: null,
		about: '',
		complexity: 1,
		ingridients: '',
		instructions: '',
		need_time: '',
		title: '',
	},
	setEditRecipeForm: form => set({ editRecipeForm: form }),

	createRecipe: async (e, login, createRecipeForm, navigate) => {
		e.preventDefault();
		try {
			const formData = new FormData();
			formData.append('title', createRecipeForm.title);
			formData.append('about', createRecipeForm.about);
			formData.append('ingridients', createRecipeForm.ingridients);
			formData.append('instructions', createRecipeForm.instructions);
			formData.append('complexity', createRecipeForm.complexity.toString());
			formData.append('need_time', createRecipeForm.need_time);

			formData.append('photos', createRecipeForm.selectedImage);

			const axiosResponse = await axios.post(
				`${baseUrl}/user/${login}/recipe`,
				formData,
				{
					withCredentials: true,
				}
			);

			console.log('response data: ', axiosResponse.data);

			if (axiosResponse.status === 200) {
				set({ recipes: axiosResponse.data.recipes });
			}

			set({
				createRecipeForm: {
					title: '',
					about: '',
					ingridients: '',
					author: {} as IAuthor,
					instructions: '',
					complexity: 1,
					need_time: '',
					selectedImage: null,
				},
			});

			navigate(-1);
		} catch (err) {
			console.log('Error creating recipe: ', err);
			alert(`Ошибка создания рецепта: ${handleError(err)} `);
		}
	},

	getAllRecipes: async () => {
		try {
			const response = await axios.get(`${baseUrl}/recipe`, {
				withCredentials: true,
			});

			console.log('ALL: ', response.data.recipes);

			if (response.status === 200) {
				const updatedRecipes = response.data.recipes.map((recipe: IRecipe) => ({
					...recipe,
					photos_urls: deleteLastChar(recipe.photos_urls),
				}));

				set({ recipes: updatedRecipes });

				console.log('response: ', response.data.recipes);
			}
		} catch (err) {
			console.log('Error getting recipes: ', err);
		}
	},

	getFilteredRecipes: async (filter: Filter) => {
		try {
			const response = await axios.post(
				`${baseUrl}/recipe`,
				{ ...filter },
				{
					withCredentials: true,
				}
			);

			console.log('FILTERED: ', response.data.recipes);

			if (response.status === 200) {
				const updatedRecipes = response.data.recipes.map((recipe: IRecipe) => ({
					...recipe,
					photos_urls: deleteLastChar(recipe.photos_urls),
				}));

				set({ recipes: updatedRecipes });
				set({ isLoading: false });
			}
		} catch (err) {
			console.log('Error getting filtered recipes: ', err);
		}
	},

	getRecipeAndComments: async recipeId => {
		try {
			const response = await axios.get(`${baseUrl}/recipe/${recipeId}`, {
				withCredentials: true,
			});

			console.log('RECIPE: ', response.data.info.recipe);
			console.log('get rec response: ', response);

			if (response.status === 200) {
				const updatedRecipe = {
					...response.data.info.recipe,
					photos_urls: deleteLastChar(response.data.info.recipe.photos_urls),
				};

				const updatedComments = response.data.info.comments.map(
					(comment: IComment) => ({
						...comment,
						author: {
							...comment.author,
							icon_url: deleteLastChar(comment.author.icon_url),
						},
					})
				);

				console.log('res data comments: ', response);

				set({ isLiked: response.data.info.is_liked });
				set({ likes_count: response.data.info.likes_count });
				set({ recipe: updatedRecipe });
				set({ comments: updatedComments });
				set({ text: '' });
			}
		} catch (err) {
			console.log('Error getting recipe: ', err);
			alert(`Ошибка получения рецепта: ${handleError(err)}`);
		}
	},

	getAuthor: async login => {
		try {
			const response = await axios.get(`${baseUrl}/user/${login}`, {
				withCredentials: true,
			});

			console.log('author res: ', response);

			if (response.status === 200) {
				const author = response.data.user.login;
				return author;
			}
		} catch (err) {
			console.log('Error getting author: ', err);
		}
	},

	editRecipe: async (e, login, recipeId, editRecipeForm, navigate) => {
		e.preventDefault();
		try {
			const formData = new FormData();

			formData.append('title', editRecipeForm.title);
			formData.append('about', editRecipeForm.about);
			formData.append('ingridients', editRecipeForm.ingridients);
			formData.append('instructions', editRecipeForm.instructions);
			formData.append('complexity', editRecipeForm.complexity.toString());
			formData.append('need_time', editRecipeForm.need_time);

			formData.append('photos', editRecipeForm.photos);

			const axiosResponse = await axios.put(
				`${baseUrl}/user/${login}/recipe/${recipeId}`,
				formData,
				{
					withCredentials: true,
				}
			);

			console.log('AXIOS RES: ', axiosResponse);

			if (axiosResponse.status === 200) {
				set({ recipes: axiosResponse.data.recipes });
				set({
					editRecipeForm: {
						title: '',
						about: '',
						ingridients: '',
						instructions: '',
						complexity: 1,
						need_time: '',
						photos: null,
					},
				});

				navigate(-1);
			}
		} catch (err) {
			console.log('Error editing recipe: ', err);
			alert(`Ошибка изменения рецепта: ${handleError(err)}`);
		}
	},

	deleteRecipe: async (login, recipeId) => {
		try {
			const response = await axios.delete(
				`${baseUrl}/user/${login}/recipe/${recipeId}`,
				{ withCredentials: true }
			);

			if (response.status === 200) {
				console.log('LOAD: ', get().isLoading);

				set({ isLoading: true });
			}
		} catch (err) {
			console.log('Error deleting recipe: ', err);
			alert(`Ошибка удаления рецепта: ${handleError(err)}`);
		}
	},

	likeRecipe: async recipeId => {
		try {
			const response = await axios.post(
				`${baseUrl}/recipe/${recipeId}/like`,
				{},
				{
					withCredentials: true,
				}
			);

			if (response.status === 200) {
				set({ isLiked: true });
			}

			set({ isLiked: true });
		} catch (err) {
			console.log('Error liking recipe: ', err);
			alert(`Ошибка лайка рецепта: ${err}`);
		}
	},

	unlikeRecipe: async recipeId => {
		try {
			const response = await axios.post(
				`${baseUrl}/recipe/${recipeId}/unlike`,
				{},
				{
					withCredentials: true,
				}
			);

			if (response.status === 200) {
				set({ isLiked: false });
			}
			set({ isLiked: false });
		} catch (err) {
			console.log('Error unliking recipe: ', err);
			alert(`Ошибка удаления лайка рецепта: ${err}`);
		}
	},

	createComment: async (e, recipeId, text) => {
		e.preventDefault();
		try {
			const response = await axios.post(
				`${baseUrl}/recipe/${recipeId}/comment`,
				{ text },
				{
					withCredentials: true,
				}
			);

			console.log('comment res: ', response);

			if (response.status === 200) {
				set({ isLoading: true });
			}
		} catch (err) {
			console.log('Error creating comment: ', err);
			alert();
		}
	},

	editComment: async (commentId, text) => {
		try {
			const response = await axios.put(
				`${baseUrl}/recipe/${commentId}/comment`,
				{ id: commentId, text },
				{
					withCredentials: true,
				}
			);

			if (response.status === 200) {
				set({ isLoading: true });
			}
		} catch (err) {
			console.log('Error editing comment: ', err);
		}
	},

	deleteComment: async (commentId, recipeId) => {
		try {
			const response = await axios.delete(
				`${baseUrl}/recipe/${recipeId}/comment`,
				{
					data: { id: commentId },
					withCredentials: true,
				}
			);

			if (response.status === 200) {
				set({ isLoading: true });
			}
		} catch (err) {
			console.log('Error deleting comment: ', err);
		}
	},
}));
