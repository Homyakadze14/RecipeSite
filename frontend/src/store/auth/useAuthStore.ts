import axios from 'axios';
import Cookies from 'js-cookie';
import { NavigateFunction } from 'react-router-dom';
import { create } from 'zustand';
import { handleError } from '../recipes/useRecipesStore';

export interface IUseAuthStore {
	email: string;
	setEmail: (email: string) => void;
	login: string;
	setLogin: (login: string) => void;
	password: string;
	setPassword: (password: string) => void;
	isAuth: boolean;
	setIsAuth: (isAuth: boolean) => void;
	currentUserLogin: string;
	setCurrentUserLogin: (login: string) => void;

	signUp: (
		e: React.MouseEvent<HTMLButtonElement>,
		email: string,
		login: string,
		password: string,
		navigate: (to: string) => void,
		signIn: (
			e: React.MouseEvent<HTMLButtonElement>,
			email: string,
			password: string,
			navigate: (to: string) => void
		) => void
	) => void;
	signIn: (
		e: React.MouseEvent<HTMLButtonElement>,
		email: string,
		password: string,
		navigate: (to: string) => void
	) => void;
	logout: (navigate: NavigateFunction) => void;
}

const baseUrl = 'http://localhost:8080/api/v1/auth';

export const useAuthStore = create<IUseAuthStore>(set => ({
	email: '',
	setEmail: (email: string) => set({ email }),
	login: localStorage.getItem('login') || '',
	setLogin: (login: string) => set({ login }),
	password: '',
	setPassword: (password: string) => set({ password }),
	isAuth: JSON.parse(localStorage.getItem('isAuth') || 'false'),
	setIsAuth: (isAuth: boolean) => {
		set({ isAuth });
		localStorage.setItem('isAuth', JSON.stringify(isAuth));
	},
	currentUserLogin: '',
	setCurrentUserLogin: (login: string) => set({ currentUserLogin: login }),

	signIn: async (e, email, password, navigate) => {
		e.preventDefault();

		try {
			const response = await axios.post(`${baseUrl}/signin`, {
				email,
				password,
			});

			console.log('LOG:', response);

			if (response.status === 200) {
				set({ isAuth: true });
				localStorage.setItem('isAuth', JSON.stringify(true));
				set({ login: response.data.login });

				Cookies.set('session_id', response.data.session_id, {
					expires: 3,
				});

				localStorage.setItem('login', response.data.login);

				navigate(`/user/${response.data.login}`);

				console.log('login from lc st: ', localStorage.getItem('login'));

				console.log('response: ', response);
				console.log('login: ', useAuthStore.getState().login);
			}
		} catch (err) {
			console.log('Error signing in', err);
			alert(`Ошибка входа: ${handleError(err)}`);
		}
	},

	signUp: async (e, email, login, password, navigate, signIn) => {
		e.preventDefault();
		try {
			const response = await axios.post(`${baseUrl}/signup`, {
				email,
				login,
				password,
			});

			console.log('up res', response);

			if (response.status === 200) {
				set({ login: login });
				localStorage.setItem('login', login);

				signIn(e, email, password, navigate);
			}
		} catch (err) {
			console.log('Error signing up', err);
			alert(`Ошибка регистрации: ${handleError(err)}`);
		}
	},

	logout: async navigate => {
		try {
			const response = await axios.post(
				`${baseUrl}/logout`,
				{},
				{ withCredentials: true }
			);
			if (response.status === 200) {
				set({ isAuth: false, login: '', email: '', password: '' });

				localStorage.setItem('login', '');
				localStorage.setItem('isAuth', JSON.stringify(false));

				Cookies.remove('session_id');

				console.log('after reset login: ', useAuthStore.getState().login);

				navigate('/signin');
			}
		} catch (err) {
			console.log('Error logging out', err);
			alert(`Ошибка выхода: ${handleError(err)}`);
		}
	},
}));
