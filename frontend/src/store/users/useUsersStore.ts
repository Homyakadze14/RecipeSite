import axios, { AxiosError } from 'axios';
import Cookies from 'js-cookie';
import { NavigateFunction } from 'react-router-dom';
import { create } from 'zustand';
import { useAuthStore } from './../auth/useAuthStore';
import { handleError, useRecipesStore } from './../recipes/useRecipesStore';

export interface IAuthor {
	login: string;
	icon_url: string;
}

export interface IRecipe {
	id: number;
	about: string;
	complexity: 1 | 2 | 3;
	created_at: string;
	creator_user_id?: number;
	ingridients: string;
	instructions: string;
	need_time: string;
	title: string;
	photos_urls: string;
	updated_at: string;
	author: IAuthor;
}

export interface IUser {
	id: number;
	login: string;
	about: string;
	icon_url: string;
	created_at: string;
	recipies: IRecipe[];
	liked_recipies: IRecipe[];
	is_subscribed: boolean;
}

export interface IUseUsersStore {
	user: IUser;
	setUser: (user: IUser) => void;
	users: IUser[];
	setUsers: (users: IUser[]) => void;

	paramsLogin: string;
	setParamsLogin: (login: string) => void;

	editUserForm: {
		icon: any;
		login: string;
		about: string;
	};
	setEditUserForm: (form: any) => void;

	password: string;
	setPassword: (password: string) => void;

	getUser: (login: string) => void;

	editUser: (
		e: React.MouseEvent<HTMLButtonElement>,
		login: string,
		editUserForm: { icon: File | null; login: string; about: string },
		navigate: NavigateFunction,
		setIsEditModalVisible: (isEditModalVisible: boolean) => void
	) => void;

	editPassword: (
		login: string,
		newPassword: string,
		navigate: NavigateFunction
	) => void;

	subscribe: (login: string) => void;
	unsubscribe: (login: string) => void;
}

export const deleteLastChar = (str: string) => {
	if (str.slice(-1) == ';') {
		return str.slice(0, -1);
	} else {
		return str;
	}
};

const baseUrl = 'http://localhost:8080/api/v1/user';

export const useUsersStore = create<IUseUsersStore>(set => ({
	user: {} as IUser,
	setUser: (user: IUser) => set({ user }),
	users: [],
	setUsers: (users: IUser[]) => set({ users }),

	paramsLogin: localStorage.getItem('paramsLogin') || '',
	setParamsLogin: (login: string) => set({ paramsLogin: login }),

	editUserForm: {
		icon: null,
		login: '',
		about: '',
	},
	setEditUserForm: form => set({ editUserForm: form }),

	password: '',
	setPassword: (password: string) => set({ password }),

	getUser: async login => {
		try {
			const response = await axios.get(`${baseUrl}/${login}`, {
				withCredentials: true,
			});

			console.log('user res: ', response);

			if (response.status === 200) {
				const updatedUser = {
					...response.data.user,
					icon_url: deleteLastChar(response.data.user.icon_url),
					recipies: response.data.user.recipies.map((recipe: IRecipe) => ({
						...recipe,
						photos_urls: deleteLastChar(recipe.photos_urls),
					})),
				};
				set({ user: updatedUser });
			}

			console.log('response data from user: ', response.data);
		} catch (err) {
			console.log('Error getting user: ', err);
		}
	},

	editUser: async (e, login, editUserForm, navigate, setIsEditModalVisible) => {
		e.preventDefault();
		try {
			const formData = new FormData();

			formData.append('login', editUserForm.login);
			formData.append('about', editUserForm.about);
			formData.append('icon', editUserForm.icon ?? '');

			const response = await axios.put(`${baseUrl}/${login}`, formData, {
				withCredentials: true,
			});

			console.log(response);

			if (response.status === 200) {
				localStorage.setItem('login', editUserForm.login);
				localStorage.setItem('paramsLogin', editUserForm.login);
				set({ paramsLogin: editUserForm.login });

				useAuthStore.setState({ login: editUserForm.login });
				useRecipesStore.setState({ isLoading: true });

				console.log('user is: ', useUsersStore.getState().user);

				const updatedPath = `/user/${editUserForm.login}`;
				window.history.pushState({}, '', updatedPath);
			}
		} catch (err: AxiosError | any) {
			if (err.response && err.response.status === 400) {
				navigate(`/user/${login}`);
				setIsEditModalVisible(true);

				const user = useUsersStore.getState().user;
				useAuthStore.setState({ login });
				useUsersStore.setState({
					editUserForm: {
						login: user.login,
						about: user.about,
						icon: user.icon_url,
					},
				});

				console.log('Error editing user: ', err);
				setTimeout(
					() =>
						alert(
							`Ошибка! Такой логин уже существует. Попробуйте ввести другой.`
						),
					400
				);
			}
		}
	},

	editPassword: async (login, password, navigate) => {
		try {
			const response = await axios.put(
				`${baseUrl}/${login}/password`,
				{ password },
				{ withCredentials: true }
			);

			if (response.status === 200) {
				alert('Пароль изменен успешно');

				useAuthStore.setState({
					isAuth: false,
					login: '',
					email: '',
					password: '',
				});

				localStorage.setItem('isAuth', JSON.stringify(false));
				localStorage.setItem('login', '');

				Cookies.remove('session_id');

				navigate('/signin');
			}
		} catch (err) {
			console.log('Error editing password: ', err);
			alert(`Ошибка изменения пароля: ${handleError(err)}`);
		}
	},

	subscribe: async login => {
		try {
			const response = await axios.post(
				`${baseUrl}/${login}/subscribe`,
				{},
				{ withCredentials: true }
			);

			console.log('subscribe res: ', response);

			if (response.status === 200) {
				console.log('subscribe res: ', response);
			}
		} catch (err) {
			console.log('Error subscribing: ', err);
		}
	},

	unsubscribe: async login => {
		try {
			const response = await axios.post(
				`${baseUrl}/${login}/unsubscribe`,
				{},
				{ withCredentials: true }
			);

			console.log('unsubscribe res: ', response);

			if (response.status === 200) {
				console.log('unsubscribe res: ', response);
			}
		} catch (err) {
			console.log('Error unsubscribing: ', err);
		}
	},
}));
